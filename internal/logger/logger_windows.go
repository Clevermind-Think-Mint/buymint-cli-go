package logger

import (
	"time"

	"github.com/mattn/go-colorable"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Used to initialize human-freindly, colorize output
func logPrettyInit() error {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: colorable.NewColorableStdout(), TimeFormat: time.RFC3339})
	return nil
}
