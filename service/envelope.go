package service

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
	"red_packet/model"
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

func (p *Producer) Add(envelope *model.Envelope) {
	p.Chan <- envelope
}

func (p *Producer) Do() {
	logrus.Infof("begin producing %v envelopes with %v amount...", p.MaxLen, p.Amount)
	for {
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

// 写回红包信息以及user_count
func WriteToRedis(user *User, envelope *model.Envelope, rdb *redis.Client) error {
	data, err := json.Marshal(envelope)
	if err != nil {
		logrus.Error(err)
		return err
	}
	if err = rdb.HSet(ctx, envelope.UserId, envelope.EnvelopeId, data).Err(); err != nil {
		logrus.Error(err)
	}
	if err = rdb.HSet(ctx, "user_count", user.Uid, user.CurCount+1).Err(); err != nil {
		logrus.Error(err)
	}
	return err
}
