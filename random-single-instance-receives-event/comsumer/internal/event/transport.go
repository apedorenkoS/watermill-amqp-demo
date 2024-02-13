package event

import (
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/apedorenkoS/watermill-amqp-demo/random-single-instance-receives-event/consumer/internal/config"
)

type Transport interface {
	FanoutSubscribe(
		exchange Exchange,
		onMessage func(message *message.Message)) error
	Shutdown() error
}

type transportImpl struct {
	fanoutSub         *amqp.Subscriber
	subReaderPoolSize int
}

func NewTransport(conf config.Config, queueName string) (Transport, error) {
	subConfig := amqp.NewDurablePubSubConfig(
		conf.RabbitURI,
		amqp.GenerateQueueNameConstant(queueName))

	subscriber, err := amqp.NewSubscriber(subConfig, &zerologWatermillAdapter{})
	if err != nil {
		return nil, fmt.Errorf("create amqp subscriber: %v", err)
	}
	return &transportImpl{
		fanoutSub:         subscriber,
		subReaderPoolSize: conf.SubReaderPoolSize,
	}, nil
}

func (tr *transportImpl) FanoutSubscribe(
	exchange Exchange,
	onMessage func(message *message.Message)) error {

	messages, err := tr.fanoutSub.Subscribe(context.Background(), string(exchange))
	if err != nil {
		return fmt.Errorf("subscribe to %s: %v", exchange, err)
	}

	for i := 0; i < tr.subReaderPoolSize; i++ {
		go processMessages(messages, onMessage)
	}
	return nil
}

func processMessages(
	messages <-chan *message.Message,
	onMessage func(message *message.Message)) {

	for {
		select {
		case msg, ok := <-messages:
			if !ok {
				return
			}
			onMessage(msg)
			msg.Ack()
		}
	}
}

func (tr *transportImpl) Shutdown() error {
	return tr.fanoutSub.Close()
}
