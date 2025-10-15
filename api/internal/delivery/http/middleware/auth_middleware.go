package middleware

import (
	"api/internal/model"
	"api/internal/usecase"
	"api/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func NewAuth(userUseCase *usecase.UserUseCaseImpl, tokenUtil *utils.TokenUtil) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		request := &model.VerifyUserRequest{Token: ctx.Get("Authorization", "NOT_FOUND")}
		userUseCase.Log.Debugf("Authorization : %s", request.Token)

		//auth, err := userUseCase.Verify(ctx.UserContext(), request)
		//if err != nil {
		//	userUseCase.Log.Warnf("Failed find user by token : %+v", err)
		//	return fiber.ErrUnauthorized
		//}

		auth, err := tokenUtil.ParseToken(ctx.UserContext(), request.Token)
		if err != nil {
			userUseCase.Log.Warnf("Failed find user by token : %+v", err)
			return fiber.ErrUnauthorized
		}

		userUseCase.Log.Debugf("User : %+v", auth.ID)
		ctx.Locals("auth", auth)
		return ctx.Next()
	}
}

func GetUser(ctx *fiber.Ctx) *model.Auth {
	return ctx.Locals("auth").(*model.Auth)
}
