package service

import (
	"encoding/json"
	"fmt"
	"red_envelope/config"
	"red_envelope/model"
	"red_envelope/utils"
	"testing"
	"time"
)

func TestKafka_Client(t *testing.T) {

	kafkaBrokers := utils.GetEnv("KAFKA_ADDRS", config.DefaultKafkaBrokers)
	brokers := utils.GetArgs(kafkaBrokers)
	topic := utils.GetEnv("KAFKA_TOPIC", config.DefaultKafkaTopic)
	kafkaProducer := GetKafkaProducer(topic, brokers)
	defer kafkaProducer.producer.Close()

	sum := 100
	for i := 0; i < sum; i++ {
		now := time.Now().Unix()
		tmp := model.Envelope{
			EnvelopeId: fmt.Sprintf("test message EnvelopeId %v from kafkatest %v", i, now),
			Value:      0,
			Opened:     false,
			SnatchTime: now,
			UserId:     fmt.Sprintf("test message UserId %v from kafkatest %v", i, now),
		}
		encode, _ := json.Marshal(tmp)
		kafkaProducer.Send(encode)
	}
	select {
	case err := <-kafkaProducer.producer.Errors():
		t.Log(err.Msg.Value)
		t.Fatal(err.Error())
	case <-time.After(time.Second * 1):
		break
	}

}
