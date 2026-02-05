package repository

import (
	"api/internal/entity"
	"api/internal/entity/enum"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PeriodRepository interface {
	FindByDate(db *gorm.DB, periodType enum.PeriodType, date time.Time) (*entity.Period, error)
	FindLastClosedMothly(db *gorm.DB, month, year int) (*entity.Period, error)
	FindLastWeeklyInMonth(db *gorm.DB, month, year int) (*entity.Period, error)
	Create(db *gorm.DB, period *entity.Period) (*entity.Period, error)
	FindByStartDate(db *gorm.DB, startDate time.Time) (*entity.Period, error)
}

type periodRepositoryImpl struct {
	Log *logrus.Logger
}

func NewPeriodRepository(log *logrus.Logger) PeriodRepository {
	return &periodRepositoryImpl{
		Log: log,
	}
}

func (r *periodRepositoryImpl) FindByDate(db *gorm.DB, periodType enum.PeriodType, date time.Time) (*entity.Period, error) {
	var period entity.Period

	err := db.Where(
		"type = ? AND start_date <= ? AND end_date >= ? AND deleted_at IS NULL",
		periodType, date, date,
	).First(&period).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		r.Log.WithError(err).Error("error finding period by date")
		return nil, err
	}

	return &period, nil
}

func (r *periodRepositoryImpl) FindLastClosedMothly(db *gorm.DB, month, year int) (*entity.Period, error) {
	var period entity.Period

	err := db.Where(
		"type = ? AND month = ? AND year = ? AND is_closed = ? AND deleted_at IS NULL",
		enum.MONTHLY, month, year, true,
	).Order("week_number DESC").First(&period).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		r.Log.WithError(err).Error("error finding last closed weekly period")
		return nil, err
	}

	return &period, nil
}

func (r *periodRepositoryImpl) Create(db *gorm.DB, period *entity.Period) (*entity.Period, error) {
	if err := db.Create(period).Error; err != nil {
		r.Log.WithError(err).Error("error creating period")
		return nil, err
	}
	return period, nil
}

func (r *periodRepositoryImpl) FindLastWeeklyInMonth(db *gorm.DB, month, year int) (*entity.Period, error) {
	var period entity.Period
	err := db.Where("type = ? AND month = ? AND year = ?", "WEEKLY", month, year).
		Order("week_number DESC").
		First(&period).Error

	if err != nil {
		return nil, err
	}
	return &period, nil
}

func (r *periodRepositoryImpl) FindByStartDate(db *gorm.DB, startDate time.Time) (*entity.Period, error) {
	var period entity.Period
	err := db.Where("type = ? AND start_date = ?", "WEEKLY", startDate).
		First(&period).Error

	if err != nil {
		return nil, err
	}
	return &period, nil
}
