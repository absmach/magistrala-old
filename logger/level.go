// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package logger

import (
	"errors"
	"log/slog"
)

// ErrInvalidLogLevel indicates an unrecognized log level.
var ErrInvalidLogLevel = errors.New("unrecognized log level")

// UnmarshalText is a helper function to convert text to slog's Level type.
func UnmarshalText(text string) (slog.Level, error) {
	var level slog.Level

	err := level.UnmarshalText([]byte(text))
	if err != nil {
		return 0, ErrInvalidLogLevel
	}
	return level, nil
}
