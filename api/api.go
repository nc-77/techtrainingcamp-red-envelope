package api

import (
	"red_packet/model"
	"red_packet/service"

	"github.com/gofiber/fiber/v2"
)

var (
	app = service.GetApp()
)

func Snatch(c *fiber.Ctx) error {
	uid := c.FormValue("uid")
	// 检验uid
	if !service.CheckUid(uid) {
		return Response(c, ERRPARAM, "")
	}
	user := service.NewUser(uid)
	// 检验uid是否在黑名单中
	if user.IsAllowed() {
		return Response(c, DISABLED, "")
	}
	// 检验uid是否达到次数上限
	count := user.GetCount()
	if count >= app.MaxCount {
		return Response(c, MAXCOUNT, "")
	}
	// 获取红包
	envelope := user.GetEnvelope(app.EnvelopeProducer)
	if envelope == nil {
		return Response(c, FAILED, "")
	}
	// 更新userCount
	app.UserCount.Store(uid, count+1)

	// 同步更新redis todo
	if err := service.WriteToRedis(user, envelope, app.RDB); err != nil {
		// 失败回滚
	}

	// 异步更新mysql todo
	// db:=initialize.GetApp().DB
	return Response(c, SUCCESS, envelope)
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
			EnvelopeId: "1",
			Value:      10,
			Opened:     true,
			SnatchTime: 123456,
		},
		{
			EnvelopeId: "2",
			Opened:     false,
			SnatchTime: 123456,
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
