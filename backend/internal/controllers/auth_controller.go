package controllers

import (
	"github.com/example/git-ops/backend/internal/services"
	"github.com/gofiber/fiber/v2"
)

type AuthController struct{ auth *services.AuthService }

func NewAuthController(auth *services.AuthService) *AuthController { return &AuthController{auth: auth} }

func (c *AuthController) Register(r fiber.Router) {
	r.Post("/auth/login", c.Login)
}

func (c *AuthController) Login(ctx *fiber.Ctx) error {
	t, _ := c.auth.Token("demo-user")
	return ctx.JSON(fiber.Map{"token": t})
}
