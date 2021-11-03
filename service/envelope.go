package service

import (
	"context"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/panjf2000/ants/v2"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
)

var (
	ctx = context.Background()
)

type Producer struct {
	Amount int64
	Size   int64
	MaxLen int64
	mutex  sync.Mutex
}

func NewProducer(amount int64, size int64) *Producer {
	return &Producer{
		Amount: amount,
		Size:   size,
		MaxLen: size,
		mutex:  sync.Mutex{},
	}
}

// 启动runtimes个协程生产红包
func (p *Producer) Do(rdb *redis.Client, runtimes int) {
	defer ants.Release()
	var wg sync.WaitGroup
	for i := 0; i < runtimes; i++ {
		wg.Add(1)
		_ = ants.Submit(func() {
			defer wg.Done()
			p.do(rdb)
		})
	}
	wg.Wait()
}

// 根据Amount不断生产Size个红包放入rdb中
func (p *Producer) do(rdb *redis.Client) {
	logrus.Infof("begin producing %v envelope with %v account...", p.Size, p.Amount)
	pipe := rdb.Pipeline()
	for {
		p.mutex.Lock()
		value, ok := getRandomMoney(p.Size, p.Amount)
		if !ok {
			p.mutex.Unlock()
			break
		}
		p.Amount -= value
		p.Size--
		p.mutex.Unlock()

		cmd := pipe.XAdd(ctx, &redis.XAddArgs{
			Stream:     "envelope",
			NoMkStream: false,
			MaxLen:     p.MaxLen,
			Approx:     false,
			Values:     []interface{}{"envelope_id", xid.New().String(), "value", value, "opened", false, "snatch_time", time.Now().Unix()},
		})
		if cmd.Err() != nil {
			logrus.Error(cmd.Err())
		}

	}
	if _, err := pipe.Exec(ctx); err != nil {
		logrus.Error(err)
	}
	logrus.Info("finish producing envelopes")
}

// todo
func getRandomMoney(remainSize int64, remainMoney int64) (money int64, ok bool) {
	if remainSize <= 0 || remainMoney <= 0 {
		return
	}
	money = remainMoney / remainSize
	if money <= 0 {
		return
	}
	ok = true
	return
}
