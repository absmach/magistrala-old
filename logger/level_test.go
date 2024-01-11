// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package logger

import (
	"fmt"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalText(t *testing.T) {
	cases := []struct {
		desc   string
		input  string
		output slog.Level
		err    error
	}{
		{
			desc:   "select log level Not_A_Level",
			input:  "Not_A_Level",
			output: 0,
			err:    ErrInvalidLogLevel,
		},
		{
			desc:   "select log level Bad_Input",
			input:  "Bad_Input",
			output: 0,
			err:    ErrInvalidLogLevel,
		},
		{
			desc:   "select log level debug",
			input:  "debug",
			output: slog.LevelDebug,
			err:    nil,
		},
		{
			desc:   "select log level DEBUG",
			input:  "DEBUG",
			output: slog.LevelDebug,
			err:    nil,
		},
		{
			desc:   "select log level info",
			input:  "info",
			output: slog.LevelInfo,
			err:    nil,
		},
		{
			desc:   "select log level INFO",
			input:  "INFO",
			output: slog.LevelInfo,
			err:    nil,
		},
		{
			desc:   "select log level warn",
			input:  "warn",
			output: slog.LevelWarn,
			err:    nil,
		},
		{
			desc:   "select log level WARN",
			input:  "WARN",
			output: slog.LevelWarn,
			err:    nil,
		},
		{
			desc:   "select log level Error",
			input:  "Error",
			output: slog.LevelError,
			err:    nil,
		},
		{
			desc:   "select log level ERROR",
			input:  "ERROR",
			output: slog.LevelError,
			err:    nil,
		},
	}

	for _, tc := range cases {
		level, err := UnmarshalText(tc.input)
		assert.Equal(t, tc.output, level, fmt.Sprintf("%s: expected %s got %v", tc.desc, tc.output, level))
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %v got %v", tc.desc, tc.err, err))
	}
}
