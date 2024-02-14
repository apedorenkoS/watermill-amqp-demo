package config

type Config struct {
	RabbitURI string
}

func LoadConfig() Config {
	return Config{
		RabbitURI: "amqp://guest:guest@localhost:5672",
	}
}
