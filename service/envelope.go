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

var randomMoney []int64

const randomMoneyBit = 13
const randomMoneyLen = 1 << randomMoneyBit

// todo 按照config来生成EnvelopeDistribute数组
// 举个例子配置里将红包的分成5个等级
// 比分占比15 40 30 10 5
// 下限和上限为(1~30) (30~200) (200~1000) (1000~5000) (5000~10000)
// 为了让金额尽量用完，配置的期望要等于总金额
type EnvelopeDistribute struct {
	Probability float64
	UpperLimit  int64
	LowerLimit  int64
}

// 更改配置则重新调用,randomMoney更新根据EnvelopeDistribute，在开启服务前一定要调用
func InitRandomMoney(envelopeDistribute []EnvelopeDistribute) {
	if len(envelopeDistribute) < 1 {
		logrus.Errorln("len(envelopeDistribute) < 1")
	}
	res := make([]int64, randomMoneyLen)
	index := 0
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for _, ele := range envelopeDistribute {
		for i := float64(randomMoneyLen) * ele.Probability; i > 0 && index < randomMoneyLen; i-- {
			res[index] = ele.LowerLimit + r.Int63n(ele.UpperLimit-ele.LowerLimit+1)
		}
	}
	for ; index < randomMoneyLen; index++ {
		res[index] = envelopeDistribute[0].LowerLimit
	}
	rand.Shuffle(len(res), func(i, j int) { res[i], res[j] = res[j], res[i] })
	randomMoney = res
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
			envelopeId := xid.New().String()
			value, ok := getRandomMoney(p.Size, p.Amount, envelopeId)
			if !ok {
				p.Mutex.Unlock()
				break
			}
			p.Size--
			p.Amount -= value
			p.Mutex.Unlock()
			envelope := &model.Envelope{
				EnvelopeId: envelopeId,
				Value:      value,
				Opened:     false,
				UserId:     "",
			}
			p.Chan <- envelope

		}
		logrus.Infof("finish producing %v envelopes... ", size)

	}

}

func envelopeIdToIndex(envelopeId string) int {
	index := 0
	for i := 19; i > 20-randomMoneyBit; i-- {
		index ^= (int(envelopeId[i]) & 1) << (19 - i)
	}
	return index
}

// 二倍均值法随机分配红包 -> 打表预分配，靠envelope_id来拿
func getRandomMoney(remainSize int64, remainMoney int64, envelopeId string) (money int64, ok bool) {
	if remainSize <= 0 || remainMoney <= 0 {
		return
	}
	index := envelopeIdToIndex(envelopeId)
	money = utils.Min(randomMoney[index], remainMoney)
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
