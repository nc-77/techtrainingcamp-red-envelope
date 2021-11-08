package service

import (
	"testing"
	"time"
)

func TestUser_GetEnvelope(t *testing.T) {
	onceApp = GetApp()
	onceApp.Run()
	user := NewUser("123")
	envelope := user.SnatchEnvelope(onceApp.EnvelopeProducer)

	if envelope == nil {
		t.Fatal()
	}
	t.Log(envelope)
}

func TestUser_isSnatched(t *testing.T) {
	user := NewUser("123")
	sum := 100
	snatched := 0
	for i := 0; i < sum; i++ {
		time.Sleep(time.Nanosecond)
		if user.isSnatched() {
			snatched++
		}
	}
	t.Log(float64(snatched) / float64(sum))

}
