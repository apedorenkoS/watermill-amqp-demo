package event

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/rs/zerolog/log"
)

type zerologWatermillAdapter struct {
	fields watermill.LogFields
}

func (zwa *zerologWatermillAdapter) Error(msg string, err error, fields watermill.LogFields) {
	log.Err(err).Fields(zwa.fields.Add(fields)).Msg(msg)
}
func (zwa *zerologWatermillAdapter) Info(msg string, fields watermill.LogFields) {
	log.Info().Fields(zwa.fields.Add(fields)).Msg(msg)
}
func (zwa *zerologWatermillAdapter) Debug(msg string, fields watermill.LogFields) {
	log.Debug().Fields(zwa.fields.Add(fields)).Msg(msg)
}
func (zwa *zerologWatermillAdapter) Trace(msg string, fields watermill.LogFields) {
	log.Trace().Fields(zwa.fields.Add(fields)).Msg(msg)
}
func (zwa *zerologWatermillAdapter) With(fields watermill.LogFields) watermill.LoggerAdapter {
	return &zerologWatermillAdapter{fields: zwa.fields.Add(fields)}
}
