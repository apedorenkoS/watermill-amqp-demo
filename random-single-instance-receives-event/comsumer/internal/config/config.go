package config

type Config struct {
	RabbitURI         string
	SubReaderPoolSize int
}

func LoadConfig() Config {
	return Config{
		RabbitURI:         "amqp://guest:guest@localhost:5672",
		SubReaderPoolSize: 4,
	}
}
