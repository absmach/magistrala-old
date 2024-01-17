// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

//go:build !test

package api

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/absmach/magistrala/coap"
	"github.com/absmach/magistrala/pkg/messaging"
)

var _ coap.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger *slog.Logger
	svc    coap.Service
}

// LoggingMiddleware adds logging facilities to the adapter.
func LoggingMiddleware(svc coap.Service, logger *slog.Logger) coap.Service {
	return &loggingMiddleware{logger, svc}
}

// Publish logs the publish request. It logs the channel ID, subtopic (if any) and the time it took to complete the request.
// If the request fails, it logs the error.
func (lm *loggingMiddleware) Publish(ctx context.Context, key string, msg *messaging.Message) (err error) {
	defer func(begin time.Time) {
		destChannel := msg.GetChannel()
		if msg.GetSubtopic() != "" {
			destChannel = fmt.Sprintf("%s.%s", destChannel, msg.GetSubtopic())
		}
		message := "Method publish completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("method", "publish"),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("method", "publish"),
			slog.String("channel", destChannel),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.Publish(ctx, key, msg)
}

// Subscribe logs the subscribe request. It logs the channel ID, subtopic (if any) and the time it took to complete the request.
// If the request fails, it logs the error.
func (lm *loggingMiddleware) Subscribe(ctx context.Context, key, chanID, subtopic string, c coap.Client) (err error) {
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
			slog.String("client", c.Token()),
			slog.String("channelID", chanID),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.Subscribe(ctx, key, chanID, subtopic, c)
}

// Unsubscribe logs the unsubscribe request. It logs the channel ID, subtopic (if any) and the time it took to complete the request.
// If the request fails, it logs the error.
func (lm *loggingMiddleware) Unsubscribe(ctx context.Context, key, chanID, subtopic, token string) (err error) {
	defer func(begin time.Time) {
		destChannel := chanID
		if subtopic != "" {
			destChannel = fmt.Sprintf("%s.%s", destChannel, subtopic)
		}
		message := "Method unsubscribe completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error: %s.", message, err),
				slog.String("method", "unsubscribe"),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("method", "unsubscribe"),
			slog.String("channel", destChannel),
			slog.String("token", token),
			slog.String("key", key),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.Unsubscribe(ctx, key, chanID, subtopic, token)
}
