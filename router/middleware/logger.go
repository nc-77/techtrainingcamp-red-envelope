package middleware

import (
	"github.com/gofiber/fiber/v2"
	logger2 "github.com/gofiber/fiber/v2/middleware/logger"
)

func Logger() fiber.Handler {
	return logger2.New(logger2.Config{
		TimeZone: "Asia/ShangHai",
	})
}
