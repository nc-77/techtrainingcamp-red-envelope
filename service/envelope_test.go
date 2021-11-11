package service

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	amount int64 = 100
	size   int64 = 100
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

func Test_getRandomMoney(t *testing.T) {
	testCount := 10
	amounts := make([]int64, testCount)
	sizes := make([]int64, testCount)
	rand.Seed(time.Now().UnixNano())
	
	for i := 0; i < testCount; i++ {
		amounts[i] = rand.Int63n(math.MaxUint16)
		sizes[i] = rand.Int63n(amounts[i])
		t.Run(fmt.Sprintf("amount:%v size:%v", amounts[i], sizes[i]), func(t *testing.T) {
			beginAmount := amounts[i]
			for {
				money, ok := getRandomMoney(sizes[i], amounts[i], "c6576gjbu3ifgt3emvrg")
				if !ok {
					break
				}
				amounts[i] -= money
				sizes[i]--
			}
			assert.Equal(t, int64(0), sizes[i])
			assert.GreaterOrEqual(t, amounts[i], int64(0))
			assert.LessOrEqual(t, amounts[i], beginAmount/100)
		})
	}

}
