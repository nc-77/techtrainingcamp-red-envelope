package service

import (
	"encoding/json"
	"math"
	"math/rand"
	"sync"
	"time"

	"red_envelope/model"
	"red_envelope/utils"

	"github.com/go-redis/redis/v8"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
)

type Producer struct {
	Amount  int64
	Size    int64
	MaxLen  int64
	Chan    chan *model.Envelope
	MsgChan chan int // 启动消息通知
	Mutex   sync.Mutex
}

func NewProducer(amount int64, size int64) *Producer {
	return &Producer{
		Amount:  amount,
		Size:    size,
		Chan:    make(chan *model.Envelope, utils.Min(math.MaxUint16, utils.Max(math.MaxInt16, size))),
		MsgChan: make(chan int, 100),
		Mutex:   sync.Mutex{},
	}
}

func (p *Producer) Add(envelope *model.Envelope) {
	p.Chan <- envelope
}

func (p *Producer) Do() {
	for {
		//fmt.Println("wait")
		msg := <-p.MsgChan
		if msg == 0 {
			//fmt.Println("quit")
			return
		}
		size := p.Size
		logrus.Infof("begin producing %v envelopes with %v amount...", p.Size, p.Amount)
		for {
			p.Mutex.Lock()
			value, ok := getRandomMoney(p.Size, p.Amount)
			if !ok {
				p.Mutex.Unlock()
				break
			}
			p.Size--
			p.Amount -= value
			p.Mutex.Unlock()
			envelope := &model.Envelope{
				EnvelopeId: xid.New().String(),
				Value:      value,
				Opened:     false,
				UserId:     "",
			}
			p.Chan <- envelope

		}
		logrus.Infof("finish producing %v envelopes... ", size)

	}

}

// 二倍均值法随机分配红包
func getRandomMoney(remainSize int64, remainMoney int64) (money int64, ok bool) {
	if remainSize <= 0 || remainMoney <= 0 {
		return
	}
	n := utils.Max(remainMoney*2/remainSize-1, 1)
	rand.Seed(time.Now().UnixNano())
	money = utils.Min(rand.Int63n(n)+1, remainMoney)

	ok = true
	return
}

// 写回红包信息
func WriteToRedis(envelope *model.Envelope, rdb *redis.Client) error {
	pipe := rdb.Pipeline()
	data, err := json.Marshal(envelope)
	if err != nil {
		logrus.Error(err)
		return err
	}
	if err = pipe.HSet(ctx, envelope.UserId, envelope.EnvelopeId, data).Err(); err != nil {
		logrus.Error(err)
		return err
	}
	if err = pipe.IncrBy(ctx, "cur_amount", envelope.Value).Err(); err != nil {
		logrus.Error(err)
		return err
	}
	if err = pipe.IncrBy(ctx, "cur_size", 1).Err(); err != nil {
		logrus.Error(err)
		return err
	}
	_, err = pipe.Exec(ctx)
	return err
}

func UpdateRedis(envelope *model.Envelope, rdb *redis.Client) error {
	data, err := json.Marshal(envelope)
	if err != nil {
		logrus.Error(err)
		return err
	}
	if err = rdb.HSet(ctx, envelope.UserId, envelope.EnvelopeId, data).Err(); err != nil {
		logrus.Error(err)
		return err
	}
	return err
}
