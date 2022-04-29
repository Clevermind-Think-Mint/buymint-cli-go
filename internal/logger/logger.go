package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var prefix string

const (
	// DebugLevel defines debug log level.
	DebugLevel = zerolog.DebugLevel
	// InfoLevel defines info log level.
	InfoLevel = zerolog.InfoLevel
	// WarnLevel defines warn log level.
	WarnLevel = zerolog.WarnLevel
	// ErrorLevel defines error log level.
	ErrorLevel = zerolog.ErrorLevel
	// FatalLevel defines fatal log level.
	FatalLevel = zerolog.FatalLevel
	// PanicLevel defines panic log level.
	PanicLevel = zerolog.PanicLevel
	// NoLevel defines an absent log level.
	NoLevel = zerolog.Disabled
	// Disabled disables the logger.
	Disabled = zerolog.Disabled
)

// LogInit is ...
func LogInit(logLevel zerolog.Level, pretty bool) {
	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(logLevel)
	if pretty {
		//log.Logger = log.With().Caller().Logger()
		err := logPrettyInit()
		if err != nil {
			Warn("Failed to initialize log color: %s", err)
		}
	}
}

// IsProduction ...
func IsProduction() bool {
	return !(zerolog.GlobalLevel() == zerolog.DebugLevel)
}

// Debug ...
// We cannot use optional arguments nor adding default values
func Debug(message string, params ...interface{}) {
	log.Debug().Msgf(message, params...)
}

// Info ...
func Info(message string, params ...interface{}) {
	log.Info().Msgf(message, params...)
}

// Warn ...
func Warn(message string, params ...interface{}) {
	log.Warn().Msgf(message, params...)
}

// Error ...
func Error(err error, message string, params ...interface{}) {
	params = append(params, err)
	log.Error().Msgf(message+": %s", params...)
}

// Fatal ...
func Fatal(message string, params ...interface{}) {
	log.Fatal().Msgf(message, params...)
}
