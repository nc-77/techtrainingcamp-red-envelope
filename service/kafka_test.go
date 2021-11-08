package service

import (
	"red_envelope/config"
	"red_envelope/utils"
	"testing"
)

func TestKafka_Client(t *testing.T) {
	kafkaBrokers := utils.GetEnv("KAFKA_ADDRS", config.DefaultKafkaBrokers)
	brokers := utils.GetArgs(kafkaBrokers)
	topic := utils.GetEnv("KAFKA_TOPIC", config.DefaultKafkaTopic)
	kafkaProducer := GetKafkaProducer(topic, brokers)

	kafkaProducer.Send("test message from kafka-client-go-test")
	select {
	case <-kafkaProducer.producer.Successes():
	case err := <-kafkaProducer.producer.Errors():
		panic(err.Error)
	}
}
