package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var log zerolog.Logger

func InitLogger() {
	zerolog.TimeFieldFormat = time.RFC3339

	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "15:04:05",
		NoColor:    false,
	}

	log = zerolog.New(output).
		With().
		Timestamp().
		Caller().
		Logger()

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}

func LogDebug(format string, v ...any) {
	log.Debug().Msgf(format, v...)
}

func LogInfo(format string, v ...any) {
	log.Info().Msgf(format, v...)
}

func LogWarning(format string, v ...any) {
	log.Warn().Msgf(format, v...)
}

func LogError(format string, v ...any) {
	log.Error().Msgf(format, v...)
}

func PanicError(format string, v ...any) {
	log.Panic().Msgf(format, v...)
}
