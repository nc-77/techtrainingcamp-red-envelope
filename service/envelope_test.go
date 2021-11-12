package service

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/rs/xid"
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
	testCount := 1
	amounts := make([]int64, testCount)
	sizes := make([]int64, testCount)

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < testCount; i++ {
		amounts[i] = 6e7
		sizes[i] = 6e4
		t.Run(fmt.Sprintf("amount:%v size:%v", amounts[i], sizes[i]), func(t *testing.T) {
			tmp := amounts[i] - 0.01*10000*sizes[i]
			tmp2 := tmp - 0.3*2000*sizes[i]
			tmp3 := float64(tmp2) / float64(sizes[i]) * 0.69
			InitRandomMoney([]EnvelopeDistribute{
				{
					0.01,
					10000,
					10000,
				},
				{
					0.3,
					tmp2,
					tmp2,
				},
				{
					0.69,
					int64(tmp3),
					int64(tmp3),
				},
			})
			beginAmount := amounts[i]
			for {
				envelopeId := xid.New().String()
				t.Log(envelopeId)
				money, ok := getRandomMoney(sizes[i], amounts[i], envelopeId)
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
