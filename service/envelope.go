package service

import (
	"github.com/go-redis/redis/v8"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
	"red_packet/model"
	"strings"
	"sync"
)

type Producer struct {
	Amount int64
	Size   int64
	MaxLen int64
	Chan   chan *model.Envelope
	mutex  sync.Mutex
}

func NewProducer(amount int64, size int64) *Producer {
	return &Producer{
		Amount: amount,
		Size:   size,
		MaxLen: size,
		Chan:   make(chan *model.Envelope, size),
		mutex:  sync.Mutex{},
	}
}

func (p *Producer) Do() {
	logrus.Infof("begin producing %v envelopes with %v amount...", p.MaxLen, p.Amount)
	for i := int64(0); i < p.MaxLen; i++ {
		p.mutex.Lock()
		value, ok := getRandomMoney(p.Size, p.Amount)
		if !ok {
			p.mutex.Unlock()
			break
		}
		p.Size--
		p.Amount -= value
		p.mutex.Unlock()
		envelope := &model.Envelope{
			EnvelopeId: xid.New().String(),
			Value:      value,
			Opened:     false,
			UserId:     "",
		}
		p.Chan <- envelope
	}
	close(p.Chan)
	logrus.Infof("finish producing %v envelopes...", p.MaxLen)

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

func WriteToRedis(user *User, envelope *model.Envelope, rdb *redis.Client) (err error) {
	var key strings.Builder
	key.WriteString(user.Uid)
	key.WriteString("-")
	key.WriteString(envelope.EnvelopeId)
	if err = rdb.HMSet(ctx, key.String(), "snatch_time", envelope.SnatchTime, "value", envelope.Value, "opened", envelope.Opened).Err(); err != nil {
		logrus.Error(err)
	}
	return
}
