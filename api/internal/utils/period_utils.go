package utils

import (
	"time"
)

type PeriodUtils struct {
}

func NewPeriodUtils() *PeriodUtils {
	return &PeriodUtils{}
}

func (p *PeriodUtils) findFirstSunday(year, month int, loc *time.Location) time.Time {
	firstDay := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, loc)

	for firstDay.Weekday() != time.Sunday {
		firstDay = firstDay.AddDate(0, 0, 1)
		// Safety: if went to next month, return first day
		if firstDay.Month() != time.Month(month) {
			return time.Date(year, time.Month(month), 1, 0, 0, 0, 0, loc)
		}
	}
	return firstDay
}

func (p *PeriodUtils) calculateWeekNumberFromSunday(date, firstSunday time.Time) int {
	days := int(date.Sub(firstSunday).Hours() / 24)
	return (days / 7) + 1
}

func (p *PeriodUtils) getWeekBoundaries(firstSunday time.Time, weekNumber int) (time.Time, time.Time) {
	startDate := firstSunday.AddDate(0, 0, (weekNumber-1)*7)
	endDate := startDate.AddDate(0, 0, 6)
	return startDate, endDate
}

func (p *PeriodUtils) getPreviousMonth(year, month int) (int, int) {
	if month == 1 {
		return year - 1, 12
	}
	return year, month - 1
}
