package router

import (
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"red_packet/api"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func InitRouter() *fiber.App {
	router := fiber.New()

	v0 := router.Group("/v0")
	v0.Use(cors.New(), logger.New())
	v0.Use(limiter.New(
		limiter.Config{
			Max:        1000,
			Expiration: time.Second * 1,
			LimitReached: func(ctx *fiber.Ctx) error {
				return api.Response(ctx, api.FAILED, "")
			},
		}))
	v0.Post("/snatch", api.Snatch)
	v0.Post("/open", api.Open)
	v0.Post("/get_wallet_list", api.GetWalletList)

	return router
}
