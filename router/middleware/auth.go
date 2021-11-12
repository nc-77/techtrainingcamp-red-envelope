package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
)

func Auth() fiber.Handler {
	return basicauth.New(basicauth.Config{
		Users: map[string]string{
			app.UserAuth: app.PasswdAuth,
		},
	})
}
