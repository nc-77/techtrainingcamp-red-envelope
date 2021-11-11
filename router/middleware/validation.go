package middleware

import (
	"red_envelope/api"

	"github.com/gofiber/fiber/v2"
)

func Validate() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !checkUid(c.FormValue("uid")) {
			return api.Response(c, api.ERRPARAM, "")
		}
		return c.Next()
	}
}

func checkUid(uid string) bool {
	return len(uid) != 0
}
