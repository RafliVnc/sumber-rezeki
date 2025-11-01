package http

import (
	"api/internal/delivery/http/middleware"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type HelloController struct {
}

func NewHelloController() *HelloController {
	return &HelloController{}
}

func (h *HelloController) SayHello(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	return ctx.Send([]byte(fmt.Sprintf("Hello %v", auth.ID)))
}
