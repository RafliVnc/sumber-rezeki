package middleware

import (
	"api/internal/model"
	"api/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func NewAuth(tokenUtil *utils.TokenUtil, log *logrus.Logger, redis *redis.Client) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		request := &model.VerifyUserRequest{Token: ctx.Get("Authorization", "NOT_FOUND")}
		log.Debugf("Authorization : %s", request.Token)

		userAuth, err := tokenUtil.ParseToken(ctx.UserContext(), request.Token)
		if err != nil {
			log.Warnf("Failed find user by token : %+v", err)
			return fiber.ErrUnauthorized
		}

		log.Debugf("User : %+v", userAuth.ID)
		ctx.Locals("auth", userAuth)
		return ctx.Next()
	}
}

func GetUser(ctx *fiber.Ctx) *model.Auth {
	return ctx.Locals("auth").(*model.Auth)
}
