// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package logger

import (
	"fmt"
	"io"
	"log/slog"
	"time"
)

// New returns a new slog logger.
func New(w io.Writer, levelText string) (slog.Logger, error) {
	var slogLevel slog.Level
	err := slogLevel.UnmarshalText([]byte(levelText))
	if err != nil {
		return slog.Logger{}, fmt.Errorf(`{"level":"error","message":"%s: %s","ts":"%s"}`, err, levelText, time.RFC3339Nano)
	}

	logHandler := slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: slogLevel,
	})

	sLogger := slog.New(logHandler)

	return *sLogger, nil
}
