package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
)

func Cache() fiber.Handler {
	return cache.New(cache.Config{
		Expiration:   time.Second * 10,
		CacheHeader:  "X-Cache",
		CacheControl: true,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.FormValue("uid")
		},
	})
}
