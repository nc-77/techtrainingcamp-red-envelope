package router

import (
	"red_envelope/api"
	"red_envelope/router/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/pprof"
)

func InitRouter() *fiber.App {
	router := fiber.New()
	router.Use(pprof.New())

	v0 := router.Group("/v0")
	v0.Use(cors.New(), middleware.Logger(), middleware.Validate(), middleware.Limiter())
	v0.Post("/snatch", api.Snatch)
	v0.Post("/open", api.Open)
	v0.Post("/get_wallet_list", middleware.Cache(), api.GetWalletList)

	router.Post("/get_config", middleware.Logger(), api.GetConfig)
	router.Post("/config", middleware.Logger(), api.UpdateConfig)
	router.Post("/ping", middleware.Logger(), func(c *fiber.Ctx) error {
		return c.SendString("ping")
	})
	return router
}
