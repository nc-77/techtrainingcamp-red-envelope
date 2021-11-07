package middleware

import (
	"github.com/gofiber/fiber/v2"
	"red_packet/api"
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
