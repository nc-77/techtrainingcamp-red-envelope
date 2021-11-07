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
	user := service.NewUser(uid)

	// 检验uid是否达到次数上限
	count, err := user.GetCount()
	if err != nil {
		return Response(c, FAILED, "")
	}
	if count >= app.MaxCount {
		return Response(c, MAXCOUNT, "")
	}
	// 获取红包
	envelope := user.SnatchEnvelope(app.EnvelopeProducer)
	if envelope == nil {
		return Response(c, FAILED, "")
	}
	// 同步更新redis
	if err = service.WriteToRedis(envelope, app.RDB); err != nil {
		// 失败回滚
		app.EnvelopeProducer.Add(envelope)
		return Response(c, FAILED, "")
	}
	// 更新缓存
	app.UserCount.SetDefault(user.Uid, count+1)
	app.UserWallet.Delete(user.Uid)

	// 异步更新mysql todo
	// db:=initialize.GetApp().DB
	return Response(c, SUCCESS, fiber.Map{
		"enveloped_id": envelope.EnvelopeId,
		"max_count":    app.MaxCount,
		"cur_count":    count + 1,
	})
}

func Open(c *fiber.Ctx) error {
	var envelope *model.Envelope
	var err error
	uid := c.FormValue("uid")
	user := service.NewUser(uid)

	if envelope, err = user.GetEnvelope(c.FormValue("envelope_id")); err != nil {
		return Response(c, FAILED, "")
	}

	if envelope == nil {
		return Response(c, ERRPARAM, "")
	}
	// 同步更新redis
	envelope.Opened = true
	if err = service.WriteToRedis(envelope, app.RDB); err != nil {
		return Response(c, FAILED, "")
	}
	// 更新缓存
	app.UserWallet.Delete(user.Uid)
	// 异步更新Mysql todo

	return Response(c, SUCCESS, fiber.Map{
		"value": envelope.Value,
	})
}

func GetWalletList(c *fiber.Ctx) error {
	uid := c.FormValue("uid")
	user := service.NewUser(uid)
	wallet, err := user.GetWallet()
	if err != nil {
		return Response(c, FAILED, "")
	}
	// 隐藏value字段
	envelopes := make([]model.RespEnvelope, len(wallet))
	for i := range wallet {
		envelopes[i] = model.RespEnvelope{
			EnvelopeId: wallet[i].EnvelopeId,
			Opened:     wallet[i].Opened,
			SnatchTime: wallet[i].SnatchTime,
		}
		if envelopes[i].Opened {
			envelopes[i].Value = wallet[i].Value
		}
	}
	return Response(c, SUCCESS, envelopes)
}
