package main

import (
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/apedorenkoS/watermill-amqp-demo/specific-instance-receives-event/consumer/internal/config"
	"github.com/apedorenkoS/watermill-amqp-demo/specific-instance-receives-event/consumer/internal/event"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	conf := config.LoadConfig()

	transport, err := event.NewTransport(conf, "specific.instance.receives.event", conf.RoutingKey)
	if err != nil {
		log.Panic().Msgf("Failed to create transport: %v", err)
	}

	err = transport.DirectSubscribe(event.DirectExchange, processEvent)
	if err != nil {
		log.Panic().Msgf("Failed to subscribe to ")
	}

	wg := &sync.WaitGroup{}
	setupGracefulShutdown(wg, func() {
		err := transport.Shutdown()
		if err != nil {
			log.Err(err).Msgf("Error shutting down transport: %v", err)
		}
	})
	wg.Wait()
}

func processEvent(msg *message.Message) {
	e := event.Event{}
	err := json.Unmarshal(msg.Payload, &e)
	if err != nil {
		log.Err(err).Msgf("Error unmarshalling event: %v", err)
		return
	}

	log.Info().Msgf("Received event with ID %d, message %s", e.ID, e.Message)
}

func setupGracefulShutdown(wg *sync.WaitGroup, shutdown func()) {
	wg.Add(1)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
		<-quit
		shutdown()
		wg.Done()
	}()
}
