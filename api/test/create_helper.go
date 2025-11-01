package test

import (
	"api/internal/entity"
	"api/internal/entity/enum"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CreateUsers(total int) []entity.User {
	users := make([]entity.User, total)

	for i := 0; i < total; i++ {
		user := &entity.User{
			ID:       uuid.New(),
			Name:     "UserTest",
			Username: "user" + strconv.Itoa(i),
			Password: "password",
			Phone:    "123456" + strconv.Itoa(i),
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

func CreateSales(total int) []entity.Sales {
	sales := make([]entity.Sales, total)

	for i := 0; i < total; i++ {
		sales[i] = entity.Sales{
			Name:  "Created Sales",
			Phone: "123456" + strconv.Itoa(i),
		}

		dbErr := db.Create(&sales[i]).Error
		if dbErr != nil {
			log.Fatalf("Failed create sales data : %+v", dbErr)
		}
	}
	return sales
}

func CreateRoutes(total int) []entity.Route {
	routes := make([]entity.Route, total)

	for i := 0; i < total; i++ {
		routes[i] = entity.Route{
			Name:        "Created Route " + strconv.Itoa(i),
			Description: "Created Description Route" + strconv.Itoa(i),
		}

		dbErr := db.Create(&routes[i]).Error
		if dbErr != nil {
			log.Fatalf("Failed create route data : %+v", dbErr)
		}
	}
	return routes
}

func CreateSalesWithRoutes(count int, routeIDs []int) []entity.Sales {
	salesList := make([]entity.Sales, count)

	for i := 0; i < count; i++ {
		sales := entity.Sales{
			Name:  fmt.Sprintf("Created Sales %d", i),
			Phone: fmt.Sprintf("123456%d", i),
		}

		dbErr := db.Create(&sales).Error
		if dbErr != nil {
			log.Fatalf("Failed create sales data : %+v", dbErr)
		}

		// add routes if exist
		if len(routeIDs) > 0 {
			var routes []entity.Route
			err := db.Where("id IN ?", routeIDs).Find(&routes).Error
			if err != nil {
				log.Fatalf("Failed find routes data : %+v", err)
			}

			err = db.Model(&sales).Association("Routes").Replace(routes)
			if err != nil {
				log.Fatalf("Failed create sales routes association : %+v", err)
			}
		}

		salesList[i] = sales
	}

	return salesList
}
