package service

import (
	"testing"
	"time"
)

func TestUser_GetEnvelope(t *testing.T) {
	p := NewProducer(amount, size)
	user := NewUser("123")
	go p.Do()
	time.Sleep(time.Second)

	envelope := user.GetEnvelope(p)

	if envelope == nil {
		t.Fatal()
	}
	t.Log(envelope)
}
