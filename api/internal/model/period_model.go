package model

import "time"

type WeekInfo struct {
	Month      int
	Year       int
	WeekNumber int
	StartDate  time.Time
	EndDate    time.Time
}
