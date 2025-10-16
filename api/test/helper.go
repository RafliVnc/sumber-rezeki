package test

import (
	"api/internal/entity"
	"strconv"
)

func ClearAll() {
	ClearUsers()
}

func ClearUsers() {
	err := db.Where("id is not null").Delete(&entity.User{}).Error
	if err != nil {
		log.Fatalf("Failed clear user data : %+v", err)
	}
}

func CreateUsers(total int) []entity.User {
	users := make([]entity.User, total)

	for i := 0; i < total; i++ {
		user := &entity.User{
			Name:     "rafli",
			Username: "rafli" + strconv.Itoa(i) + "@gmail.com",
			Password: "rahasia",
		}
		err := db.Create(user).Error
		if err != nil {
			log.Fatalf("Failed create user data : %+v", err)
		}
		users[i] = *user
	}
	return users
}
