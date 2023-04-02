package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Trace() *zerolog.Event {
	return log.Trace().Caller()
}

func Debug() *zerolog.Event {
	return log.Debug().Caller()
}

func Info() *zerolog.Event {
	return log.Info().Caller()
}

func Warn() *zerolog.Event {
	return log.Warn().Caller()
}

func Error() *zerolog.Event {
	return log.Error().Caller()
}

func Fatal() *zerolog.Event {
	return log.Fatal().Caller()
}

func Panic() *zerolog.Event {
	return log.Panic().Caller()
}

func Print() *zerolog.Event {
	return log.Log().Caller()
}
