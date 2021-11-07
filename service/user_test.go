package service

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	onceApp = GetApp()
	onceApp.Run()
	m.Run()
	os.Exit(0)
}

func TestUser_GetEnvelope(t *testing.T) {

	user := NewUser("123")
	envelope := user.SnatchEnvelope(onceApp.EnvelopeProducer)

	if envelope == nil {
		t.Fatal()
	}
	t.Log(envelope)
}
