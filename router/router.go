package router

import (
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"red_envelope/api"
	"red_envelope/router/middleware"
	"time"

	"github.com/gofiber/fiber/v2"
)

func InitRouter() *fiber.App {
	router := fiber.New()

	v0 := router.Group("/v0")
	v0.Use(cors.New(), middleware.Logger(), middleware.Validate())
	v0.Use(limiter.New(
		limiter.Config{
			Max:        5,
			Expiration: time.Second * 1,
			KeyGenerator: func(c *fiber.Ctx) string {
				return c.FormValue("uid")
			},
			LimitReached: func(c *fiber.Ctx) error {
				return api.Response(c, api.FAILED, "")
			},
		}))
	v0.Post("/snatch", api.Snatch)
	v0.Post("/open", api.Open)
	v0.Post("/get_wallet_list", cache.New(cache.Config{
		Expiration:   time.Second * 30,
		CacheHeader:  "X-Cache",
		CacheControl: true,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.FormValue("uid")
		},
	}),
		api.GetWalletList)

	router.Post("/get_config", middleware.Logger(), api.GetConfig)
	router.Post("/config", middleware.Logger(), api.UpdateConfig)
	return router
}
