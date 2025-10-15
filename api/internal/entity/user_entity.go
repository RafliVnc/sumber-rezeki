package entity

type User struct {
	ID       int    `gorm:"column:id;primaryKey;autoIncrement"`
	Name     string `gorm:"column:name"`
	Email    string `gorm:"column:email"`
	Password string `gorm:"column:password"`
}

func (u *User) TableName() string {
	return "users"
}
