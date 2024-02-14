package event

import (
	"encoding/json"
	"fmt"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/apedorenkoS/watermill-amqp-demo/specific-instance-receives-event/producer/internal/config"
)

type Transport interface {
	DirectPublish(exchange Exchange, routingKey string, payload interface{}) error
	Shutdown() error
}

type transportImpl struct {
	directPub *amqp.Publisher
}

func NewTransport(conf config.Config) (Transport, error) {
	pubSubConfig := amqp.NewDurablePubSubConfig(conf.RabbitURI, nil)
	pubSubConfig.Exchange = amqp.ExchangeConfig{
		GenerateName: parseTopic,
		Type:         "direct",
		Durable:      true,
	}
	pubSubConfig.Publish.GenerateRoutingKey = parseRoutingKey

	publisher, err := amqp.NewPublisher(pubSubConfig, &zerologWatermillAdapter{})
	if err != nil {
		return nil, fmt.Errorf("create amqp publisher: %v", err)
	}
	return &transportImpl{directPub: publisher}, nil
}

func (tr *transportImpl) DirectPublish(exchange Exchange, routingKey string, payload any) error {

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal %v of type %T to json: %v", body, body, err)
	}
	pubMess := &message.Message{Payload: body}
	return tr.directPub.Publish(concatenate(string(exchange), routingKey), pubMess)
}

func (tr *transportImpl) Shutdown() error {
	return tr.directPub.Close()
}
