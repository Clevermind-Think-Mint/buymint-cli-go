//go:build !windows && !linux && !darwin
// +build !windows,!linux,!darwin

package logger

import (
	"errors"
	"runtime"
)

// Used to initialize human-freindly, colorize output
func logPrettyInit() error {
	return errors.New("Unsupported log color on: " + runtime.GOOS)
}
