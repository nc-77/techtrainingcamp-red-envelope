package service

import (
	"red_packet/model"
)

type User string

func (uid User) IsAllowed() bool {
	// todo
	return false
}

func (uid User) IsMaxCount() bool {
	// todo
	return false
}

// 根据概率判断是否抢到红包
func (uid User) isSnatched() bool {
	// todo
	return true
}

func (uid User) GetEnvelope() *model.Envelope {
	if !uid.isSnatched() {
		return nil
	}

	return &model.Envelope{}
}
