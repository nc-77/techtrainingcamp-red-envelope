package service

import (
	"fmt"
	"time"

	"github.com/Shopify/sarama"
)

type KafkaProducer struct {
	SendFail map[string]time.Ticker
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
		make(map[string]time.Ticker),
		topic,
		getProducer(addrs),
	}
}

func (kafkaProducer *KafkaProducer) Send(msg []byte) {
	// todo 如下
	kafkaProducer.SendFail[string(msg)] = time.Ticker{}
	kafkaProducer.producer.Input() <- &sarama.ProducerMessage{
		Topic: kafkaProducer.topic,
		Value: sarama.StringEncoder(msg),
	}
}

func (kafkaProducer *KafkaProducer) HandleSendErr() {
	for {
		select {
		case ok := <-kafkaProducer.producer.Successes():
			// 这里应该有一个更稳妥的存储方案，例如Successes后给用户返回或者其他更好的解决方案
			// 这里的逻辑应该是默认发送失败
			encode, err := ok.Value.Encode()
			if err != nil {
				// todo 打印
			} else {
				delete(kafkaProducer.SendFail, string(encode))
				fmt.Println("Send Ok")
			}
		case fail := <-kafkaProducer.producer.Errors():
			fmt.Println("err: ", fail.Err)
			// todo 失败后的逻辑
			encode, err := fail.Msg.Value.Encode()
			if err != nil {
				// todo 打印
			} else {
				fmt.Println(encode)
				// kafkaProducer.SendFail[string(encode)].C 开始计时或者其他，这里只是模拟一下set
				// 进入守护线程
			}
		}
	}
}
