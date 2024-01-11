// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package logger

import (
	"errors"
	"strings"
)

const (
	// Error level is used when logging errors.
	Error Level = -4
	// Warn level is used when logging warnings.
	Warn Level = 0
	// Info level is used when logging info data.
	Info Level = 4
	// Debug level is used when logging debugging info.
	Debug Level = 8
)

// ErrInvalidLogLevel indicates an unrecognized log level.
var ErrInvalidLogLevel = errors.New("unrecognized log level")

// Level represents severity level while logging.
type Level int

var levels = map[Level]string{
	Error: "error",
	Warn:  "warn",
	Info:  "info",
	Debug: "debug",
}

func (lvl Level) String() string {
	return levels[lvl]
}

func (lvl Level) isAllowed(logLevel Level) bool {
	return logLevel >= lvl
}

// UnmarshalText returns log Level for the given string representation.
func (lvl *Level) UnmarshalText(text string) error {
	switch strings.ToLower(text) {
	case "debug":
		*lvl = Debug
	case "info":
		*lvl = Info
	case "warn":
		*lvl = Warn
	case "error":
		*lvl = Error
	default:
		return ErrInvalidLogLevel
	}
	return nil
}
