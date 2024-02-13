package main

import (
	"github.com/apedorenkoS/watermill-amqp-demo/random-single-instance-receives-event/publisher/internal/config"
	"github.com/apedorenkoS/watermill-amqp-demo/random-single-instance-receives-event/publisher/internal/event"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	conf := config.LoadConfig()

	transport, err := event.NewTransport(conf)
	if err != nil {
		log.Panic().Msgf("Failed to create transport: %v", err)
	}

	loop := true
	wg := &sync.WaitGroup{}
	setupGracefulShutdown(wg, func() {
		loop = false
		err := transport.Shutdown()
		if err != nil {
			log.Err(err).Msgf("Error shutting down transport: %v", err)
		}
	})

	i := 0
	for loop {
		e := event.Event{ID: i}
		log.Info().Msgf("Sending event with ID %d", i)

		err := transport.FanoutPublish(event.FanoutExchange, e)
		if err != nil {
			log.Err(err).Msgf("Error sending event %v: %v", e, err)
		}

		time.Sleep(2 * time.Second)
		i++
	}

	wg.Wait()
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
