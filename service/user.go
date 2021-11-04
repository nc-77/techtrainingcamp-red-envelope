package service

import (
	"github.com/sirupsen/logrus"
	"red_packet/model"
	"time"
)

type User struct {
	Uid string
}

func NewUser(uid string) *User {
	return &User{Uid: uid}
}

func CheckUid(uid string) bool {
	return uid != ""
}
func (user *User) IsAllowed() bool {
	// todo
	return false
}

func (user *User) GetCount() int {
	count, exist := onceApp.UserCount.LoadAndDelete(user.Uid)
	if !exist {
		return 0
	}
	cnt := count.(int)
	return cnt
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
