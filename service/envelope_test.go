package service

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

const (
	amount int64 = 1e8
	size   int64 = 1e8
)

func TestProducer_do(t *testing.T) {
	var wg sync.WaitGroup
	producer := NewProducer(amount, size)

	wg.Add(1)
	go func() {
		defer wg.Done()
		producer.Do()
		close(producer.Chan)
	}()

	producer.MsgChan <- 1
	producer.MsgChan <- 0
	var sum int64
	for envelope := range producer.Chan {
		if envelope != nil {
			sum++
		}
	}
	assert.Equal(t, size, sum, "they should be equal")
	wg.Wait()
}
