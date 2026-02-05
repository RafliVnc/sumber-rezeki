package repository

import (
	"api/internal/entity"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EmployeeAttendanceRepository interface {
	BatchUpsert(db *gorm.DB, employee []*entity.EmployeeAttendance) error
	BatchDeleteByDate(db *gorm.DB, date []time.Time) error
}

type employeeAttendanceRepositoryImpl struct {
	Log *logrus.Logger
}

func NewEmployeeAttendanceRepository(log *logrus.Logger) EmployeeAttendanceRepository {
	return &employeeAttendanceRepositoryImpl{
		Log: log,
	}
}

func (r *employeeAttendanceRepositoryImpl) BatchUpsert(db *gorm.DB, attendances []*entity.EmployeeAttendance) error {
	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "date"},
			{Name: "employee_id"},
		},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"status":     gorm.Expr("EXCLUDED.status"),
			"period_id":  gorm.Expr("EXCLUDED.period_id"),
			"updated_at": gorm.Expr("EXCLUDED.updated_at"),
			"deleted_at": nil, // Restore soft deleted
		}),
	}).Create(&attendances).Error
}

func (r *employeeAttendanceRepositoryImpl) BatchDeleteByDate(db *gorm.DB, date []time.Time) error {
	return db.Where("date IN ?", date).Delete(&entity.EmployeeAttendance{}).Error
}
