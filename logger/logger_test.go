// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package logger_test

import (
	"encoding/json"
	"fmt"
	"testing"

	mglog "github.com/absmach/magistrala/logger"
	"github.com/stretchr/testify/assert"
)

type mockWriter struct {
	value []byte
}

const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
)

func (writer *mockWriter) Write(p []byte) (int, error) {
	writer.value = append(writer.value, p...)
	return len(p), nil
}

func (writer *mockWriter) Read() (logMsg, error) {
	if len(writer.value) == 0 {
		return logMsg{}, nil
	}

	var output logMsg
	err := json.Unmarshal(writer.value, &output)
	return output, err
}

type logMsg struct {
	Time    string `json:"time"`
	Level   string `json:"level"`
	Message string `json:"msg"`
}

func TestDebug(t *testing.T) {
	cases := []struct {
		desc   string
		input  string
		level  string
		output logMsg
	}{
		{
			desc:   "debug log ordinary string",
			input:  "input_string",
			level:  LevelDebug,
			output: logMsg{Level: "DEBUG", Message: "input_string"},
		},
		{
			desc:   "debug log empty string",
			input:  "",
			level:  LevelDebug,
			output: logMsg{Level: "DEBUG", Message: ""},
		},
		{
			desc:   "debug ordinary string lvl not allowed",
			input:  "input_string",
			level:  LevelInfo,
			output: logMsg{Level: "", Message: ""},
		},
		{
			desc:   "debug empty string lvl not allowed",
			input:  "",
			level:  LevelInfo,
			output: logMsg{Level: "", Message: ""},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			writer := &mockWriter{}
			logger, err := mglog.New(writer, tc.level)
			assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))

			logger.Debug(tc.input)
			output, err := writer.Read()
			assert.Nil(t, err)

			if tc.level != LevelDebug {
				assert.Empty(t, output.Message, fmt.Sprintf("%s: expected no output got %s", tc.desc, output.Message))
			} else {
				assert.Equal(t, tc.output.Level, output.Level, fmt.Sprintf("%s: expected Level %v got %v", tc.desc, tc.output.Level, output.Level))
				assert.Equal(t, tc.output.Message, output.Message, fmt.Sprintf("%s: expected Message %v got %v", tc.desc, tc.output.Message, output.Message))
			}
		})
	}
}

func TestInfo(t *testing.T) {
	cases := []struct {
		desc   string
		input  string
		level  string
		output logMsg
	}{
		{
			desc:   "info log ordinary string",
			input:  "input_string",
			level:  LevelInfo,
			output: logMsg{Level: "INFO", Message: "input_string"},
		},
		{
			desc:   "info log empty string",
			input:  "",
			level:  LevelInfo,
			output: logMsg{Level: "INFO", Message: ""},
		},
		{
			desc:   "info ordinary string lvl not allowed",
			input:  "input_string",
			level:  LevelWarn,
			output: logMsg{Level: "", Message: ""},
		},
		{
			desc:   "info empty string lvl not allowed",
			input:  "",
			level:  LevelWarn,
			output: logMsg{Level: "", Message: ""},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			writer := &mockWriter{}
			logger, err := mglog.New(writer, tc.level)
			assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))

			logger.Info(tc.input)
			output, err := writer.Read()
			assert.Nil(t, err)

			if tc.level != LevelInfo && tc.level != LevelDebug {
				assert.Empty(t, output.Message, fmt.Sprintf("%s: expected no output got %s", tc.desc, output.Message))
			} else {
				assert.Equal(t, tc.output.Level, output.Level, fmt.Sprintf("%s: expected Level %v got %v", tc.desc, tc.output.Level, output.Level))
				assert.Equal(t, tc.output.Message, output.Message, fmt.Sprintf("%s: expected Message %v got %v", tc.desc, tc.output.Message, output.Message))
			}
		})
	}
}

func TestWarn(t *testing.T) {
	cases := []struct {
		desc   string
		input  string
		level  string
		output logMsg
	}{
		{
			desc:   "warn log ordinary string",
			input:  "input_string",
			level:  LevelWarn,
			output: logMsg{Level: "WARN", Message: "input_string"},
		},
		{
			desc:   "warn log empty string",
			input:  "",
			level:  LevelWarn,
			output: logMsg{Level: "WARN", Message: ""},
		},
		{
			desc:   "warn ordinary string lvl not allowed",
			input:  "input_string",
			level:  LevelError,
			output: logMsg{Level: "", Message: ""},
		},
		{
			desc:   "warn empty string lvl not allowed",
			input:  "",
			level:  LevelError,
			output: logMsg{Level: "", Message: ""},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			writer := &mockWriter{}
			logger, err := mglog.New(writer, tc.level)
			assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))

			logger.Warn(tc.input)
			output, err := writer.Read()
			assert.Nil(t, err)

			if tc.level != LevelWarn && tc.level != LevelInfo && tc.level != LevelDebug {
				assert.Empty(t, output.Message, fmt.Sprintf("%s: expected no output got %s", tc.desc, output.Message))
			} else {
				assert.Equal(t, tc.output.Level, output.Level, fmt.Sprintf("%s: expected Level %v got %v", tc.desc, tc.output.Level, output.Level))
				assert.Equal(t, tc.output.Message, output.Message, fmt.Sprintf("%s: expected Message %v got %v", tc.desc, tc.output.Message, output.Message))
			}
		})
	}
}

func TestError(t *testing.T) {
	cases := []struct {
		desc   string
		input  string
		level  string
		output logMsg
	}{
		{
			desc:   "error log ordinary string",
			input:  "input_string",
			level:  LevelError,
			output: logMsg{Level: "ERROR", Message: "input_string"},
		},
		{
			desc:   "error log empty string",
			input:  "",
			level:  LevelError,
			output: logMsg{Level: "ERROR", Message: ""},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			writer := &mockWriter{}
			logger, err := mglog.New(writer, tc.level)
			assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))

			logger.Error(tc.input)
			output, err := writer.Read()
			assert.Nil(t, err)

			assert.Equal(t, tc.output.Level, output.Level, fmt.Sprintf("%s: expected Level %v got %v", tc.desc, tc.output.Level, output.Level))
			assert.Equal(t, tc.output.Message, output.Message, fmt.Sprintf("%s: expected Message %v got %v", tc.desc, tc.output.Message, output.Message))
		})
	}
}
