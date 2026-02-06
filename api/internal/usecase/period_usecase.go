package usecase

import (
	"api/internal/entity"
	"api/internal/entity/enum"
	"api/internal/model"
	"api/internal/repository"
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PeriodUseCase interface {
	GetOrCreatePeriodIdByDate(ctx context.Context, date string) (int, error)
}

type PeriodUseCaseImpl struct {
	DB               *gorm.DB
	Log              *logrus.Logger
	Validate         *validator.Validate
	PeriodRepository repository.PeriodRepository
}

func NewPeriodUseCase(
	db *gorm.DB,
	logger *logrus.Logger,
	validate *validator.Validate,
	employeeRepository repository.PeriodRepository,
) PeriodUseCase {
	return &PeriodUseCaseImpl{
		DB:               db,
		Log:              logger,
		Validate:         validate,
		PeriodRepository: employeeRepository,
	}
}

func (u *PeriodUseCaseImpl) GetOrCreatePeriodIdByDate(ctx context.Context, date string) (int, error) {
	// find period by date
	period, err := u.FindWeekByDate(ctx, date)
	if err != nil {
		u.Log.Warnf("Failed to find period by date: %+v", err)
		return 0, fiber.ErrInternalServerError
	}

	// create period if not found
	if period == nil {
		period, err := u.CreateByDate(ctx, date)
		if err != nil {
			u.Log.Warnf("Failed create period to database : %+v", err)
			return 0, err
		}

		return period.ID, nil
	}

	return period.ID, nil
}

func (u *PeriodUseCaseImpl) FindWeekByDate(ctx context.Context, date string) (*entity.Period, error) {
	// parse date
	parseDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		u.Log.Warnf("Failed arsing date : %+v", err)
		return nil, err
	}

	// find period
	period, err := u.PeriodRepository.FindByDate(u.DB.WithContext(ctx), enum.WEEKLY, parseDate)
	if err != nil {
		u.Log.Warnf("Failed find period to database : %+v", err)
		return nil, err
	}

	return period, nil
}

func (u *PeriodUseCaseImpl) CreateByDate(ctx context.Context, date string) (*entity.Period, error) {
	// parse date
	parseDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		u.Log.Warnf("Failed arsing date : %+v", err)
		return nil, err
	}

	weekInfo := u.CalculateWeekInfo(ctx, parseDate)

	// parse to entity
	period := &entity.Period{
		Type:       enum.WEEKLY,
		StartDate:  weekInfo.StartDate,
		EndDate:    weekInfo.EndDate,
		WeekNumber: weekInfo.WeekNumber,
		Month:      weekInfo.Month,
		Year:       weekInfo.Year,
		IsActive:   true,
		IsClosed:   false,
	}

	// create period
	newPeriod, err := u.PeriodRepository.Create(u.DB.WithContext(ctx), period)
	if err != nil {
		u.Log.Warnf("Failed find period to database : %+v", err)
		return nil, err
	}

	return newPeriod, nil
}

func (u *PeriodUseCaseImpl) CalculateWeekInfo(ctx context.Context, date time.Time) *model.WeekInfo {
	year, _ := date.ISOWeek()

	// find closest sunday
	weekday := int(date.Weekday())
	daysFromSunday := weekday
	if weekday == 0 {
		daysFromSunday = 0
	}

	startDate := date.AddDate(0, 0, -daysFromSunday)
	endDate := startDate.AddDate(0, 0, 6)

	// Cari Minggu pertama di bulan dimana startDate berada
	firstDayOfMonth := time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, startDate.Location())
	firstSunday := firstDayOfMonth
	firstWeekday := int(firstDayOfMonth.Weekday())

	if firstWeekday != 0 {
		daysToSunday := 7 - firstWeekday
		firstSunday = firstDayOfMonth.AddDate(0, 0, daysToSunday)
	}

	// Hitung nomor minggu dan tentukan bulan yang tepat
	var week int
	var monthValue int

	if startDate.Before(firstSunday) {
		// Minggu ini dimulai di bulan sebelumnya
		lastDayPrevMonth := firstDayOfMonth.AddDate(0, 0, -1)
		prevMonthWeekInfo := u.CalculateWeekInfo(ctx, lastDayPrevMonth)
		week = prevMonthWeekInfo.WeekNumber + 1
		monthValue = prevMonthWeekInfo.Month
	} else {
		// Minggu ini di bulan yang sama dengan startDate

		// CEK: Apakah firstSunday sudah dipakai bulan sebelumnya?
		firstSundayWeekInfo, _ := u.getWeeklyPeriodByStartDate(ctx, firstSunday)

		if firstSundayWeekInfo != nil && firstSundayWeekInfo.Month != int(startDate.Month()) {
			// FirstSunday sudah dipakai bulan lain
			// Ini adalah minggu pertama atau seterusnya di bulan baru

			// Cari minggu kedua setelah firstSunday (yang merupakan minggu pertama bulan baru)
			secondSunday := firstSunday.AddDate(0, 0, 7)

			if startDate.Equal(secondSunday) || startDate.After(secondSunday) {
				// Hitung dari secondSunday sebagai week 1
				daysDiff := startDate.Sub(secondSunday).Hours() / 24
				week = int(daysDiff/7) + 1
			} else {
				// Ini tidak mungkin terjadi karena startDate >= firstSunday
				week = 1
			}

			monthValue = int(startDate.Month())
		} else {
			// FirstSunday belum dipakai atau dipakai bulan yang sama
			// Cek apakah bulan sebelumnya ada yang closed
			prevMonth := firstDayOfMonth.AddDate(0, -1, 0)
			lastClosedPeriod, _ := u.getLastClosedMonthlyPeriod(ctx, prevMonth.Year(), int(prevMonth.Month()))

			if lastClosedPeriod != nil {
				// Ada monthly period yang closed, reset ke bulan baru
				daysDiff := startDate.Sub(firstSunday).Hours() / 24
				week = int(daysDiff/7) + 1
				monthValue = int(startDate.Month())
			} else {
				// Tidak ada yang closed, cek apakah ada period weekly di bulan sebelumnya
				lastWeeklyPeriod, _ := u.getLastWeeklyPeriodInMonth(ctx, prevMonth.Year(), int(prevMonth.Month()))

				if lastWeeklyPeriod != nil {
					// Ada period weekly sebelumnya, lanjutkan dari bulan tersebut
					daysDiff := startDate.Sub(firstSunday).Hours() / 24
					week = lastWeeklyPeriod.WeekNumber + int(daysDiff/7) + 1
					monthValue = lastWeeklyPeriod.Month
				} else {
					// Tidak ada period weekly sebelumnya (bulan pertama sistem)
					daysDiff := startDate.Sub(firstSunday).Hours() / 24
					week = int(daysDiff/7) + 1
					monthValue = int(startDate.Month())
				}
			}
		}
	}

	return &model.WeekInfo{
		Year:       year,
		WeekNumber: week,
		Month:      monthValue,
		StartDate:  startDate,
		EndDate:    endDate,
	}
}

func (u *PeriodUseCaseImpl) getLastClosedMonthlyPeriod(ctx context.Context, year, month int) (*entity.Period, error) {
	// find last period month
	period, err := u.PeriodRepository.FindLastClosedMothly(u.DB.WithContext(ctx), month, year)
	if err != nil {
		u.Log.Warnf("Failed find period to database : %+v", err)
		return nil, err
	}

	return period, nil
}

func (u *PeriodUseCaseImpl) getLastWeeklyPeriodInMonth(ctx context.Context, year, month int) (*entity.Period, error) {
	period, err := u.PeriodRepository.FindLastWeeklyInMonth(u.DB.WithContext(ctx), month, year)
	if err != nil {
		u.Log.Warnf("Failed find weekly period: %+v", err)
		return nil, err
	}
	return period, nil
}

func (u *PeriodUseCaseImpl) getWeeklyPeriodByStartDate(ctx context.Context, startDate time.Time) (*entity.Period, error) {
	period, err := u.PeriodRepository.FindByStartDate(u.DB.WithContext(ctx), startDate)
	if err != nil {
		return nil, err
	}
	return period, nil
}
