package config

import "flag"

type Config struct {
	RabbitURI         string
	RoutingKey        string
	SubReaderPoolSize int
}

func LoadConfig() Config {
	routingKey := flag.String("routing.key", "#", "Routing key to subscribe to queue with")
	flag.Parse()

	return Config{
		RabbitURI:         "amqp://guest:guest@localhost:5672",
		RoutingKey:        *routingKey,
		SubReaderPoolSize: 4,
	}
}
