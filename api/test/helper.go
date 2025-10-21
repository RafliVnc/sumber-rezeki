package test

import (
	"api/internal/entity"
	"api/internal/entity/enum"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func ClearAll() {
	ClearUsers()
}

func ClearUsers() {
	err := db.Unscoped().Where("id is not null").Delete(&entity.User{}).Error
	if err != nil {
		log.Fatalf("Failed clear user data : %+v", err)
	}
}

func CreateUsers(total int) []entity.User {
	users := make([]entity.User, total)

	for i := 0; i < total; i++ {
		user := &entity.User{
			ID:       uuid.New(),
			Name:     "rafli",
			Username: "rafli" + strconv.Itoa(i),
			Password: "rahasia",
			Phone:    "1231241" + strconv.Itoa(i),
			Role:     enum.SUPER_ADMIN,
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Failed generate password : %+v", err)
		}

		user.Password = string(hashedPassword)

		dbErr := db.Create(user).Error
		if dbErr != nil {
			log.Fatalf("Failed create user data : %+v", dbErr)
		}
		users[i] = *user
	}
	return users
}

func GenerateTokenHelper() (string, error) {
	jwtSecret := "testSecret"

	// create User
	user := &entity.User{
		ID:       uuid.New(),
		Name:     "admin",
		Username: "admin",
		Password: "rahasia",
		Phone:    "999999999",
		Role:     enum.SUPER_ADMIN,
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed generate password : %+v", err)
	}

	user.Password = string(hashedPassword)

	dbErr := db.Create(user).Error
	if dbErr != nil {
		log.Fatalf("Failed create user data : %+v", dbErr)
	}

	// create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":     user.ID,
		"expire": time.Now().Add(time.Hour * 24 * 30).UnixMilli(),
	})

	jwtToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}
