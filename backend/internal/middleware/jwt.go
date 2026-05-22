package middleware

import "github.com/gofiber/fiber/v2"

func JWT(_ string) fiber.Handler { return func(c *fiber.Ctx) error { return c.Next() } }
func RequestLogger() fiber.Handler { return func(c *fiber.Ctx) error { return c.Next() } }
