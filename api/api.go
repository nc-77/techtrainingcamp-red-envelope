package api

import (
	"context"
	"encoding/json"
	"sort"
	"sync"

	"red_envelope/model"
	"red_envelope/service"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

var (
	app = service.GetApp()
	ctx = context.Background()
)

func Snatch(c *fiber.Ctx) error {

	var mutex *sync.Mutex
	uid := c.Locals("uid").(string)
	user := service.NewUser(uid)
	defer func() {
		mutex.Unlock()
	}()
	val, _ := app.UserMutex.LoadOrStore(uid, new(sync.Mutex))
	mutex = val.(*sync.Mutex)
	mutex.Lock()
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

	// kafka异步更新mysql
	s, err := json.Marshal(envelope)
	if err != nil {
		logrus.Error(err)
		// 这个是不允许的错误，相当于不能存到数据库
	} else {
		app.KafkaProducer.Send(s)
	}

	return Response(c, SUCCESS, fiber.Map{
		"enveloped_id": envelope.EnvelopeId,
		"max_count":    app.MaxCount,
		"cur_count":    count + 1,
	})
}

func Open(c *fiber.Ctx) error {
	var envelope *model.Envelope
	var mutex *sync.Mutex
	var err error
	uid := c.Locals("uid").(string)
	user := service.NewUser(uid)
	defer func() {
		mutex.Unlock()
	}()
	val, _ := app.UserMutex.LoadOrStore(uid, new(sync.Mutex))
	mutex = val.(*sync.Mutex)
	mutex.Lock()
	if envelope, err = user.GetEnvelope(c.Locals("envelope_id").(string)); err != nil {
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

	// kafka异步更新mysql
	s, err := json.Marshal(envelope)
	if err != nil {
		logrus.Error(err)
		// 这个是不允许的错误，相当于不能存到数据库
	} else {
		app.KafkaProducer.Send(s)
	}

	return Response(c, SUCCESS, fiber.Map{
		"value": envelope.Value,
	})
}

func GetWalletList(c *fiber.Ctx) error {
	uid := c.Locals("uid").(string)
	user := service.NewUser(uid)
	wallet, err := user.GetWallet()
	if err != nil {
		return Response(c, FAILED, "")
	}
	var amount int64
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
			amount += envelopes[i].Value
		}
	}
	sort.Slice(envelopes, func(i, j int) bool {
		return envelopes[i].SnatchTime > envelopes[j].SnatchTime
	})
	return Response(c, SUCCESS, fiber.Map{
		"amount":        amount,
		"envelope_list": envelopes,
	})
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
	var config model.Config
	if err := c.BodyParser(&config); err != nil {
		return Response(c, ERRPARAM, "")
	}

	if config.SnatchedPr > 0 && config.SnatchedPr <= 100 {
		app.SnatchedPr = config.SnatchedPr
		updated = true
	}

	if config.MaxCount != 0 {
		app.MaxCount = config.MaxCount
		updated = true
	}

	if val := config.MaxAmount; val != 0 {
		app.AddAmount(val)
		updatedAmount = true
		if err := app.RDB.IncrBy(ctx, "max_amount", val).Err(); err != nil {
			app.RollbackAddAmount(val)
			updatedAmount = false
		}
		updated = updated || updatedAmount
	}

	if val := config.MaxSize; val != 0 {
		app.AddSize(val)
		updatedSize = true
		if err := app.RDB.IncrBy(ctx, "max_size", val).Err(); err != nil {
			app.RollbackAddSize(val)
			updatedSize = false
		}
		updated = updated || updatedSize
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
