package service

import (
	"github.com/stretchr/testify/assert"
	"os"
	"sync"
	"testing"
)

const (
	amount int64 = 1e3
	size   int64 = 1e3
)

func TestMain(m *testing.M) {
	m.Run()

	os.Exit(0)
}

func TestProducer_do(t *testing.T) {
	var wg sync.WaitGroup
	producer := NewProducer(amount, size)

	wg.Add(1)
	go func() {
		defer wg.Done()
		producer.Do()
	}()

	var sum int64
	for envelope := range producer.Chan {
		if envelope != nil {
			sum++
		}
	}
	assert.Equal(t, size, sum, "they should be equal")
	wg.Wait()
}
