package event

type Exchange string

const (
	FanoutExchange Exchange = "random.single.instance.receives.event"
)
