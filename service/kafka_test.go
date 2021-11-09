package service

import (
	"encoding/json"
	"fmt"
	"red_packet/config"
	"red_packet/model"
	"red_packet/utils"
	"testing"
	"time"
)

func TestKafka_Client(t *testing.T) {
	kafkaBrokers := utils.GetEnv("KAFKA_ADDRS", config.DefaultKafkaBrokers)
	brokers := utils.GetArgs(kafkaBrokers)
	topic := utils.GetEnv("KAFKA_TOPIC", config.DefaultKafkaTopic)
	kafkaProducer := GetKafkaProducer(topic, brokers)
	defer kafkaProducer.producer.Close()
	for i := 0; i < 100; i++ {
		time := time.Now().Unix()
		tmp := model.Envelope{
			EnvelopeId: fmt.Sprintf("test message EnvelopeId %v from kafkatest %v", i, time),
			Value:      0,
			Opened:     false,
			SnatchTime: time,
			UserId:     fmt.Sprintf("test message UserId %v from kafkatest %v", i, time),
		}
		encode, _ := json.Marshal(tmp)
		kafkaProducer.Send(encode)
	}
	time.Sleep(5 * time.Second)
	for i := 0; i < 100; i++ {
		select {
		case <-kafkaProducer.producer.Successes():
		case err := <-kafkaProducer.producer.Errors():
			panic(err.Error)
		}
	}
}
