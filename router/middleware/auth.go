package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"red_envelope/api"
)

func Auth() fiber.Handler {
	return basicauth.New(basicauth.Config{
		Users: map[string]string{
			app.UserAuth: app.PasswdAuth,
		},
		Unauthorized: func(c *fiber.Ctx) error {
			return api.Response(c, api.UNAUTHORIZED, "")
		},
	})
}
