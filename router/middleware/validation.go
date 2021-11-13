package middleware

import (
	"red_envelope/api"
	"red_envelope/service"

	"github.com/gofiber/fiber/v2"
)

func Validate() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var user service.User
		if err := c.BodyParser(&user); err != nil {
			return api.Response(c, api.ERRPARAM, "")
		}
		if !checkUid(user.Uid) {
			return api.Response(c, api.ERRPARAM, "")
		}
		c.Locals("uid", user.Uid)
		c.Locals("envelope_id", user.EnvelopeId)

		return c.Next()
	}
}

func checkUid(uid string) bool {
	return len(uid) != 0
}
