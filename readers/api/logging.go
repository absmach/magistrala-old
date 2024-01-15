// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

//go:build !test

package api

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/absmach/magistrala/readers"
)

var _ readers.MessageRepository = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger *slog.Logger
	svc    readers.MessageRepository
}

// LoggingMiddleware adds logging facilities to the core service.
func LoggingMiddleware(svc readers.MessageRepository, logger *slog.Logger) readers.MessageRepository {
	return &loggingMiddleware{
		logger: logger,
		svc:    svc,
	}
}

func (lm *loggingMiddleware) ReadAll(chanID string, rpm readers.PageMetadata) (page readers.MessagesPage, err error) {
	defer func(begin time.Time) {
		message := "Method completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("method", "read_all"),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("method", "read_all"),
			slog.String("channel_ID", chanID),
			slog.String("query", fmt.Sprintf("%v", rpm)),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.ReadAll(chanID, rpm)
}
