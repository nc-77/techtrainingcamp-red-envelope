package service

import (
	"github.com/sirupsen/logrus"
	"red_packet/model"
	"strconv"
	"time"
)

type User struct {
	Uid      string
	CurCount int
}

func NewUser(uid string) *User {
	return &User{
		Uid:      uid,
		CurCount: getCount(uid),
	}
}

func CheckUid(uid string) bool {
	return uid != ""
}

func getCount(uid string) (cnt int) {
	var err error
	// 先从内存中取，如果不存在则从redis中取
	count, exist := onceApp.UserCount.Load(uid)
	if !exist {
		var result string
		if result, err = onceApp.RDB.HGet(ctx, "user_count", uid).Result(); err != nil {
			// 从redis中也没有为新用户
			//logrus.Info("new user")
			return
		}
		//logrus.Info("from redis",result)
		if cnt, err = strconv.Atoi(result); err != nil {
			logrus.Error(err)
		}
		// 更新内存
		onceApp.UserCount.Store(uid, cnt)
		return
	}
	cnt = count.(int)
	//logrus.Info("from map",cnt)
	return
}

// 根据概率判断是否抢到红包
func (user *User) isSnatched() bool {
	// todo
	return true
}

func (user *User) GetEnvelope(p *Producer) *model.Envelope {
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
