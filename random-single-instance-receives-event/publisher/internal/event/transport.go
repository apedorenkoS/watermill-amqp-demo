package event

import (
	"encoding/json"
	"fmt"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/apedorenkoS/watermill-amqp-demo/random-single-instance-receives-event/publisher/internal/config"
)

type Transport interface {
	FanoutPublish(exchange Exchange, payload interface{}) error
	Shutdown() error
}

type transportImpl struct {
	fanoutPub *amqp.Publisher
}

func NewTransport(conf config.Config) (Transport, error) {
	pubSubConfig := amqp.NewDurablePubSubConfig(conf.RabbitURI, nil)

	publisher, err := amqp.NewPublisher(pubSubConfig, &zerologWatermillAdapter{})
	if err != nil {
		return nil, fmt.Errorf("create amqp publisher: %v", err)
	}
	return &transportImpl{fanoutPub: publisher}, nil
}

func (tr *transportImpl) FanoutPublish(exchange Exchange, payload any) error {

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal %v of type %T to json: %v", body, body, err)
	}
	pubMess := &message.Message{Payload: body}
	return tr.fanoutPub.Publish(string(exchange), pubMess)
}

func (tr *transportImpl) Shutdown() error {
	return tr.fanoutPub.Close()
}
