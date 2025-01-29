// internal/log/logger.go
package log

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogger() {
	log.Logger = NewLogger("default")
	log.Info().Msg("Logger initialized successfully")
}

func NewLogger(component string) zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339

	level := zerolog.InfoLevel
	if lvl, err := zerolog.ParseLevel(os.Getenv("LOG_LEVEL")); err == nil {
		level = lvl
	}
	zerolog.SetGlobalLevel(level)

	logger := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("component", component).
		Logger()

	return logger
}

func Info() *zerolog.Event {
	return log.Info()
}

func Error() *zerolog.Event {
	return log.Error()
}

func Debug() *zerolog.Event {
	return log.Debug()
}
