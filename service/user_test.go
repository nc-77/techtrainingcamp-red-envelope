package service

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	m.Run()
	os.Exit(0)
}

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
