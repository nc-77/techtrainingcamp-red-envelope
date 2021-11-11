package service

import (
	"encoding/json"

	"red_envelope/model"

	"github.com/Shopify/sarama"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type KafkaProducer struct {
	topic    string
	producer sarama.AsyncProducer
}

func getProducer(brokers []string) sarama.AsyncProducer {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认，确保信息传输完毕
	saramaConfig.Producer.Partitioner = sarama.NewRandomPartitioner // 随机选择partition
	saramaConfig.Producer.Return.Errors = true                      // 发送失败的消息会在error chan中返回
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

func (kafkaProducer *KafkaProducer) Send(msg []byte) {
	kafkaProducer.producer.Input() <- &sarama.ProducerMessage{
		Topic: kafkaProducer.topic,
		Value: sarama.StringEncoder(msg),
	}
}

func (kafkaProducer *KafkaProducer) HandleSendErr() {
	for {
		fail := <-kafkaProducer.producer.Errors()
		logrus.Error(fail.Err)
		data, err := fail.Msg.Value.Encode()
		if err != nil {
			logrus.Error(err)
			break
		}
		envelope := &model.Envelope{}
		if err = json.Unmarshal(data, envelope); err != nil {
			logrus.Error(err)
			break
		}
		// 发送失败的msg持久化到redis stream中
		if err = onceApp.RDB.XAdd(ctx, &redis.XAddArgs{
			Stream: "send_failed",
			Values: envelope.ToMap(),
		}).Err(); err != nil {
			logrus.Error(err)
		}
	}
}
