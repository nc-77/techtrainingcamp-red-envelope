package api

import (
	"encoding/json"
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

	// 检验uid是否达到次数上限
	if user.CurCount >= app.MaxCount {
		return Response(c, MAXCOUNT, "")
	}
	// 获取红包
	envelope := user.GetEnvelope(app.EnvelopeProducer)
	if envelope == nil {
		return Response(c, FAILED, "")
	}
	// 同步更新redis
	if err := service.WriteToRedis(user, envelope, app.RDB); err != nil {
		// 失败回滚
		app.EnvelopeProducer.Add(envelope)
		return Response(c, FAILED, "")
	}
	// 更新userCount
	app.UserCount.Store(uid, user.CurCount+1)

	// 异步更新mysql todo
	s, err := json.Marshal(envelope)
	if err != nil {
		// todo 打印错误
		// 这个是不允许的错误，相当于不能存到数据库
	} else {
		app.KafkaProducer.Send(s)
	}
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
