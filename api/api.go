package api

import (
	"github.com/gofiber/fiber/v2"
	"red_packet/model"
)

func Snatch(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"code": 0,
		"msg":  "success",
		"data": fiber.Map{
			"envelope_id": 123,
			"max_count":   5,
			"cur_count":   3,
		},
	})
}

func Open(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"code": 0,
		"msg":  "success",
		"data": fiber.Map{
			"value": 5,
		},
	})
}

func GetWalletList(c *fiber.Ctx) error {
	envelopes := []model.Envelope{
		{
			EnvelopeId: 1,
			Value:      10,
			Opened:     true,
			SnatchTime: "123456",
		},
		{
			EnvelopeId: 2,
			Opened:     false,
			SnatchTime: "123456",
		},
	}

	return c.JSON(fiber.Map{
		"code": 0,
		"msg":  "success",
		"data": fiber.Map{
			"amount":        12,
			"envelope_list": envelopes,
		},
	})
}
