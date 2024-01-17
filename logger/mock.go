// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package logger

import (
	"bytes"
	"log/slog"
)

var _ Logger = (*loggerMock)(nil)

type loggerMock struct{}

// NewMock returns wrapped logger mock.
func NewMock() *slog.Logger {
	buf := &bytes.Buffer{}

	return slog.New(slog.NewJSONHandler(buf, nil))
}

// NewKitMock returns wrapped go kit logger mock.
func NewKitMock() Logger {
	return &loggerMock{}
}

func (l loggerMock) Debug(msg string) {
}

func (l loggerMock) Info(msg string) {
}

func (l loggerMock) Warn(msg string) {
}

func (l loggerMock) Error(msg string) {
}

func (l loggerMock) Fatal(msg string) {
}
