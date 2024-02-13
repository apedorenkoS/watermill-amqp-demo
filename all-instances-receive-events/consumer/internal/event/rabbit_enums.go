package event

type Exchange string

const (
	FanoutExchange Exchange = "all.instances.receive.events"
)
