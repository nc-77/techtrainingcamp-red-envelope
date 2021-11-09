package service

import (
	"fmt"
	"red_packet/config"
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
		kafkaProducer.Send([]byte(fmt.Sprintf("test message %v from kafka-client-go-test", i)))
	}
	time.Sleep(5 * time.Second)
	select {
	case <-kafkaProducer.producer.Successes():
	case err := <-kafkaProducer.producer.Errors():
		panic(err.Error)
	}
}
