package api

import (
	"strconv"

	"red_packet/model"
	"red_packet/service"

	"github.com/gofiber/fiber/v2"
)

func Snatch(c *fiber.Ctx) error {
	uid, err := strconv.ParseUint(c.FormValue("uid"), 10, 64)
	if err != nil {
		return response(c, ERRPARAM, "")
	}
	user := service.User(uid)
	// 检验uid是否在黑名单中
	if user.IsAllowed() {
		return response(c, DISABLED, "")
	}
	// 检验uid是否达到次数上限
	if user.IsMaxCount() {
		return response(c, MAXCOUNT, "")
	}
	// 获取红包
	envelope := user.GetEnvelope()
	if envelope == nil {
		return response(c, FAILED, "")
	}
	// 同步更新redis todo
	// rdb:=initialize.NewApp().RDB

	// 异步更新mysql todo
	// db:=initialize.NewApp().DB

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
