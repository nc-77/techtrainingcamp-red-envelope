package middleware

import (
	"time"

	"red_envelope/api"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func Limiter() fiber.Handler {
	return limiter.New(
		limiter.Config{
			Max:        1,
			Expiration: time.Second * 1,
			KeyGenerator: func(c *fiber.Ctx) string {
				return c.FormValue("uid")
			},
			LimitReached: func(c *fiber.Ctx) error {
				return api.Response(c, api.FAILED, "")
			},
		})
}
