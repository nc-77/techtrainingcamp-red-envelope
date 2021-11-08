package service

import (
	"fmt"

	"github.com/Shopify/sarama"
)

type KafkaProducer struct {
	topic    string
	producer sarama.AsyncProducer
}

func getProducer(brokers []string) sarama.AsyncProducer {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认，确保信息传输完毕
	saramaConfig.Producer.Partitioner = sarama.NewRandomPartitioner // 随机选择partition
	saramaConfig.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回, 用于删除消息
	client, err := sarama.NewAsyncProducer(brokers, saramaConfig)
	if err != nil {
		panic("client kafka err: " + err.Error())
	}
	return client
}

func GetKafkaProducer(topic string, addrs []string) KafkaProducer {
	return KafkaProducer{
		topic,
		getProducer(addrs),
	}
}

func (kafkaProducer *KafkaProducer) Send(msg string) {
	kafkaProducer.producer.Input() <- &sarama.ProducerMessage{
		Topic: kafkaProducer.topic,
		Value: sarama.StringEncoder(msg),
	}
}

func (kafkaProducer *KafkaProducer) HandleSendErr() {
	for {
		select {
		case <-kafkaProducer.producer.Successes():
			fmt.Println("Send Ok")
		case fail := <-kafkaProducer.producer.Errors():
			fmt.Println("err: ", fail.Err)
			// todo 失败后存储信息
		}
	}
}
