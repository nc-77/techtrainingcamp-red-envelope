package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"red_envelope/model"
	"red_envelope/service"
	"strconv"
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
	// 是否还有红包
	if len(app.EnvelopeProducer.Chan) == 0 {
		logrus.Info("no envelope")
		return Response(c, FAILED, "")
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
	if err = service.UpdateRedis(envelope, app.RDB); err != nil {
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

func GetConfig(c *fiber.Ctx) error {
	curAmount, err := app.GetCurAmount()
	if err != nil {
		return Response(c, FAILED, "")
	}
	curSize, err := app.GetCurSize()
	if err != nil {
		return Response(c, FAILED, "")
	}
	return Response(c, SUCCESS, fiber.Map{
		"snatched_pr": app.SnatchedPr,
		"max_count":   app.MaxCount,
		"max_amount":  app.MaxAmount,
		"max_size":    app.MaxSize,
		"cur_amount":  curAmount,
		"cur_size":    curSize,
	})
}

func UpdateConfig(c *fiber.Ctx) error {
	var updated, updatedAmount, updatedSize bool
	snatchedPr := c.FormValue("snatched_pr")
	count := c.FormValue("max_count")
	amount := c.FormValue("amount")
	size := c.FormValue("size")

	if val, ok := service.CheckSnatchedPr(snatchedPr); ok {
		app.SnatchedPr = val
		updated = true
	}

	if val, err := strconv.Atoi(count); err == nil {
		app.MaxCount = val
		updated = true
	}

	if val, err := strconv.ParseInt(amount, 10, 64); err == nil {
		app.MaxAmount += val
		app.RemainingAmount += val
		app.EnvelopeProducer.Mutex.Lock()
		app.EnvelopeProducer.Amount += val
		app.EnvelopeProducer.Mutex.Unlock()
		updatedAmount = true
		updated = true
	}

	if val, err := strconv.ParseInt(size, 10, 64); err == nil {
		app.MaxSize += val
		app.RemainingSize += val
		app.EnvelopeProducer.Mutex.Lock()
		app.EnvelopeProducer.Size += val
		app.EnvelopeProducer.Mutex.Unlock()
		updatedSize = true
		updated = true
	}
	if updatedAmount || updatedSize {
		app.EnvelopeProducer.MsgChan <- 1
	}
	if !updated {
		return Response(c, ERRPARAM, "")
	}

	return Response(c, SUCCESS, fiber.Map{
		"snatched_pr": app.SnatchedPr,
		"max_count":   app.MaxCount,
		"max_amount":  app.MaxAmount,
		"max_size":    app.MaxSize,
	})
}
