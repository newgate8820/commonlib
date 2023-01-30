package dbcommon

import (
	"time"

	"github.com/Shopify/sarama"
)

/*************************
* author: Dev0026
* createTime: 19-1-20
* updateTime: 19-1-20
* description:
*************************/

var senders = make([]*Sender, 0)

// ############### kafka Sender(异步) ###############
type Sender struct {
	prod sarama.AsyncProducer
}

func WithSenderDefaultConfig() *sarama.Config {
	cfg := sarama.NewConfig()
	//cfg.Producer.RequiredAcks = sarama.WaitForLocal
	cfg.Producer.MaxMessageBytes = int(sarama.MaxRequestSize)
	cfg.Producer.Partitioner = sarama.NewRandomPartitioner // 随机partition
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true
	cfg.Producer.Timeout = 5 * time.Second
	return cfg
}

// NewSender 创建一个异步kafka
func NewSender(addrs ...string) (*Sender, error) {
	prod, err := sarama.NewAsyncProducer(addrs, WithSenderDefaultConfig())
	if err != nil {
		return nil, err
	}
	sender := &Sender{
		prod: prod,
	}
	senders = append(senders, sender)
	return sender, nil
}

// sender 错误时
func (s *Sender) OnError(errFunc func(*sarama.ProducerError)) {
	if s == nil {
		return
	}

	for err := range s.prod.Errors() {
		if err != nil {
			errFunc(err)
		}
	}
}

// sender 成功时
func (s *Sender) OnSuccess(succFunc func(*sarama.ProducerMessage)) {
	if s == nil {
		return
	}

	for succ := range s.prod.Successes() {
		if succFunc != nil {
			succFunc(succ)
		}
	}
}

// sender 发送
func (s *Sender) Send(topic, key string, data []byte) {
	if s == nil {
		return
	}
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(data),
	}
	s.prod.Input() <- msg
}

func CloseKafka() {
	for _, one := range senders {
		one.prod.Close()
	}
}
