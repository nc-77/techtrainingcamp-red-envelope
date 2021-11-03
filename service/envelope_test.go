package service

import (
	"os"
	"red_packet/initialize"
	"testing"
)

const (
	amount int64 = 1e6
	size   int64 = 1e6
)

var (
	app = initialize.NewApp()
)

func TestMain(m *testing.M) {
	app.OpenRedis()

	m.Run()

	os.Exit(0)
}

func TestProducer_Do(t *testing.T) {
	runtimes := 1000
	producer := NewProducer(amount, size)
	producer.Do(app.RDB, runtimes)
}
