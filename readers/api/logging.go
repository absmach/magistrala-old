// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

//go:build !test

package api

import (
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
		args := []interface{}{
			slog.String("duration", time.Since(begin).String()),
			slog.String("channel_ID", chanID),
			slog.Group(
				"page_metadata",
				slog.Any("offset", rpm.Offset),
				slog.Any("limit", rpm.Limit),
				slog.String("subtopic", rpm.Subtopic),
				slog.String("publisher", rpm.Publisher),
			),
		}
		if err != nil {
			args = append(args, slog.String("error", err.Error()))
			lm.logger.Warn("Read all failed to complete successfully", args...)
			return
		}
		lm.logger.Info("Read all completed successfully", args...)
	}(time.Now())

	return lm.svc.ReadAll(chanID, rpm)
}
