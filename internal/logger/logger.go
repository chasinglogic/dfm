package logger

import "github.com/rs/zerolog"

var logger zerolog.Logger

func SetLogger(l zerolog.Logger) {
	logger = l
}

func Debug() *zerolog.Event {
	return logger.Debug()
}

func Info() *zerolog.Event {
	return logger.Info()
}
