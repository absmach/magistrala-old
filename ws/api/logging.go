// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/absmach/magistrala/ws"
)

var _ ws.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger *slog.Logger
	svc    ws.Service
}

// LoggingMiddleware adds logging facilities to the websocket service.
func LoggingMiddleware(svc ws.Service, logger *slog.Logger) ws.Service {
	return &loggingMiddleware{logger, svc}
}

// Subscribe logs the subscribe request. It logs the channel and subtopic(if present) and the time it took to complete the request.
// If the request fails, it logs the error.
func (lm *loggingMiddleware) Subscribe(ctx context.Context, thingKey, chanID, subtopic string, c *ws.Client) (err error) {
	defer func(begin time.Time) {
		destChannel := chanID
		if subtopic != "" {
			destChannel = fmt.Sprintf("%s.%s", destChannel, subtopic)
		}
		message := "Method subscribe completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("method", "subscribe"),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("method", "subscribe"),
			slog.String("channel", destChannel),
			slog.String("thing_key", thingKey),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.Subscribe(ctx, thingKey, chanID, subtopic, c)
}
