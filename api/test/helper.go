package test

import (
	"api/internal/entity"
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func ClearAll() {
	ClearSalesRoutes()
	ClearSales()
	ClearEmployees()
	ClearRoutes()
	ClearUsers()
}

func ClearUsers() {
	err := db.Unscoped().Not("username = ?", "superadmin").Delete(&entity.User{}).Error
	if err != nil {
		log.Fatalf("Failed clear user data : %+v", err)
	}
}

func ClearEmployees() {
	err := db.Unscoped().Where("id IS NOT NULL").Delete(&entity.Employee{}).Error
	if err != nil {
		log.Fatalf("Failed clear employees data : %+v", err)
	}
}

func ClearSales() {
	err := db.Unscoped().Where("id IS NOT NULL").Delete(&entity.Sales{}).Error
	if err != nil {
		log.Fatalf("Failed clear sales data : %+v", err)
	}
}

func ClearSalesRoutes() {
	// Hapus junction table sales_routes
	err := db.Exec("DELETE FROM sales_routes WHERE id IS NOT NULL").Error
	if err != nil {
		log.Fatalf("Failed clear sales_routes data : %+v", err)
	}
}

func ClearRoutes() {
	err := db.Unscoped().Where("id IS NOT NULL").Delete(&entity.Route{}).Error
	if err != nil {
		log.Fatalf("Failed clear routes data : %+v", err)
	}
}

func GenerateTokenHelper() (string, error) {
	jwtSecret := viperConfig.GetString("secret_key")

	username := "superadmin"

	user := &entity.User{}
	if err := db.First(&user, "username = ?", username).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil
		}
		return "nil", err
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

	_, err = redisClient.SetEx(context.Background(), jwtToken, user.ID, time.Hour*25*30).Result()
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}
