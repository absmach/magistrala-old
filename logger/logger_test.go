// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package logger_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"testing"

	mglog "github.com/absmach/magistrala/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	_      io.Writer = (*mockWriter)(nil)
	logger mglog.Logger
	err    error
	output logMsg
)

type mockWriter struct {
	value []byte
}

func (writer *mockWriter) Write(p []byte) (int, error) {
	writer.value = p
	return len(p), nil
}

func (writer *mockWriter) Read() (logMsg, error) {
	var output logMsg
	err := json.Unmarshal(writer.value, &output)
	return output, err
}

type logMsg struct {
	Level   string `json:"level"`
	Message string `json:"message"`
}

func TestDebug(t *testing.T) {
	ctx := context.Background()
	cases := []struct {
		desc   string
		input  string
		level  string
		output logMsg
	}{
		{
			desc:   "debug log ordinary string",
			input:  "input_string",
			level:  mglog.Debug.String(),
			output: logMsg{mglog.Debug.String(), "input_string"},
		},
		{
			desc:   "debug log empty string",
			input:  "",
			level:  mglog.Debug.String(),
			output: logMsg{mglog.Debug.String(), ""},
		},
		{
			desc:   "debug ordinary string lvl not allowed",
			input:  "input_string",
			level:  mglog.Info.String(),
			output: logMsg{"", ""},
		},
		{
			desc:   "debug empty string lvl not allowed",
			input:  "",
			level:  mglog.Info.String(),
			output: logMsg{"", ""},
		},
	}

	for _, tc := range cases {
		writer := mockWriter{}
		logger, err = mglog.New(&writer, tc.level)
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		logger.Debug(ctx, tc.input)
		output, err = writer.Read()
		assert.Equal(t, tc.output, output, fmt.Sprintf("%s: expected %s got %s", tc.desc, tc.output, output))
	}
}

func TestInfo(t *testing.T) {
	ctx := context.Background()
	cases := []struct {
		desc   string
		input  string
		level  string
		output logMsg
	}{
		{
			desc:   "info log ordinary string",
			input:  "input_string",
			level:  mglog.Info.String(),
			output: logMsg{mglog.Info.String(), "input_string"},
		},
		{
			desc:   "info log empty string",
			input:  "",
			level:  mglog.Info.String(),
			output: logMsg{mglog.Info.String(), ""},
		},
		{
			desc:   "info ordinary string lvl not allowed",
			input:  "input_string",
			level:  mglog.Warn.String(),
			output: logMsg{"", ""},
		},
		{
			desc:   "info empty string lvl not allowed",
			input:  "",
			level:  mglog.Warn.String(),
			output: logMsg{"", ""},
		},
	}

	for _, tc := range cases {
		writer := mockWriter{}
		logger, err = mglog.New(&writer, tc.level)
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		logger.Info(ctx, tc.input)
		output, err = writer.Read()
		assert.Equal(t, tc.output, output, fmt.Sprintf("%s: expected %s got %s", tc.desc, tc.output, output))
	}
}

func TestWarn(t *testing.T) {
	ctx := context.Background()
	cases := []struct {
		desc   string
		input  string
		level  string
		output logMsg
	}{
		{
			desc:   "warn log ordinary string",
			input:  "input_string",
			level:  mglog.Warn.String(),
			output: logMsg{mglog.Warn.String(), "input_string"},
		},
		{
			desc:   "warn log empty string",
			input:  "",
			level:  mglog.Warn.String(),
			output: logMsg{mglog.Warn.String(), ""},
		},
		{
			desc:   "warn ordinary string lvl not allowed",
			input:  "input_string",
			level:  mglog.Error.String(),
			output: logMsg{"", ""},
		},
		{
			desc:   "warn empty string lvl not allowed",
			input:  "",
			level:  mglog.Error.String(),
			output: logMsg{"", ""},
		},
	}

	for _, tc := range cases {
		writer := mockWriter{}
		logger, err = mglog.New(&writer, tc.level)
		require.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		logger.Warn(ctx, tc.input)
		output, err = writer.Read()
		assert.Equal(t, tc.output, output, fmt.Sprintf("%s: expected %s got %s", tc.desc, tc.output, output))
	}
}

func TestError(t *testing.T) {
	ctx := context.Background()
	cases := []struct {
		desc   string
		input  string
		output logMsg
	}{
		{
			desc:   "error log ordinary string",
			input:  "input_string",
			output: logMsg{mglog.Error.String(), "input_string"},
		},
		{
			desc:   "error log empty string",
			input:  "",
			output: logMsg{mglog.Error.String(), ""},
		},
	}

	writer := mockWriter{}
	logger, err := mglog.New(&writer, mglog.Error.String())
	require.Nil(t, err)
	for _, tc := range cases {
		logger.Error(ctx, tc.input)
		output, err := writer.Read()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		assert.Equal(t, tc.output, output, fmt.Sprintf("%s: expected %s got %s", tc.desc, tc.output, output))
	}
}
