package service

import (
	"math/rand"
	"time"

	"red_envelope/model"
	"red_envelope/utils"

	"github.com/sirupsen/logrus"
)

type User struct {
	Uid string `json:"uid"`
}

func NewUser(uid string) *User {
	return &User{
		Uid: uid,
	}
}

func (user *User) GetCount() (cnt int, err error) {
	uid := user.Uid
	// 先从cache中取，如果不存在则从redis中取
	count, exist := onceApp.UserCount.Get(uid)
	if !exist {
		var result int64
		if result, err = onceApp.RDB.HLen(ctx, uid).Result(); err != nil {
			logrus.Error(err)
			return
		}
		//logrus.Info("from redis")
		cnt = int(result)
		onceApp.UserCount.SetDefault(user.Uid, cnt)
		return
	}
	cnt = count.(int)
	//logrus.Info("from cache",cnt)
	return
}

func (user *User) GetWallet() (wallet []*model.Envelope, err error) {
	uid := user.Uid
	// 先从cache中取，如果不存在则从redis中取
	envelopes, exist := onceApp.UserWallet.Get(uid)
	if !exist {
		var result map[string]string
		if result, err = onceApp.RDB.HGetAll(ctx, uid).Result(); err != nil {
			logrus.Error(err)
			return
		}
		if wallet, err = utils.DecodeWallet(result); err != nil {
			logrus.Error(err)
			return
		}
		//logrus.Info("from redis")
		onceApp.UserWallet.SetDefault(user.Uid, wallet)
		return
	}
	//logrus.Info("from cache")
	wallet = envelopes.([]*model.Envelope)

	return
}

// 根据概率判断是否抢到红包
func (user *User) isSnatched() bool {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(100)
	return n <= onceApp.SnatchedPr
}

// 从Producer中按照取一个红包
func (user *User) SnatchEnvelope(p *Producer) *model.Envelope {
	if !user.isSnatched() {
		return nil
	}
	select {
	case envelope := <-p.Chan:
		if envelope != nil {
			envelope.SnatchTime = time.Now().Unix()
			envelope.UserId = user.Uid
		}
		return envelope
	case <-time.After(time.Second * 1):
		logrus.Error("get envelope time out...")
		return nil
	}
}

// 根据uid以及envelopeId获取红包信息
func (user *User) GetEnvelope(envelopeId string) (envelope *model.Envelope, err error) {
	var envelopes []*model.Envelope
	if envelopes, err = user.GetWallet(); err != nil {
		return nil, err
	}
	for _, envelope := range envelopes {
		if envelope != nil && envelope.EnvelopeId == envelopeId && !envelope.Opened {
			return envelope, nil
		}
	}
	return nil, err
}
