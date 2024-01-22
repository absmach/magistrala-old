// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

//go:build !test

package handler

import (
	"context"
	"log/slog"
	"time"

	"github.com/absmach/mproxy/pkg/session"
)

var _ session.Handler = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger *slog.Logger
	svc    session.Handler
}

// AuthConnect implements session.Handler.
func (lm *loggingMiddleware) AuthConnect(ctx context.Context) (err error) {
	defer func(begin time.Time) {
		args := []interface{}{
			slog.String("duration", time.Since(begin).String()),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("AuthConnect failed to complete successfully", args...)
			return
		}
		lm.logger.Info("AuthConnect completed successfully", args...)
	}(time.Now())

	return lm.svc.AuthConnect(ctx)
}

// AuthPublish implements session.Handler.
func (lm *loggingMiddleware) AuthPublish(ctx context.Context, topic *string, payload *[]byte) (err error) {
	defer func(begin time.Time) {
		args := []interface{}{
			slog.String("duration", time.Since(begin).String()),
			slog.String("topic", *topic),
			slog.Any("payload", payload),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("AuthPublish failed to complete successfully", args...)
			return
		}
		lm.logger.Info("AuthPublish completed successfully", args...)
	}(time.Now())

	return lm.svc.AuthPublish(ctx, topic, payload)
}

// AuthSubscribe implements session.Handler.
func (lm *loggingMiddleware) AuthSubscribe(ctx context.Context, topics *[]string) (err error) {
	defer func(begin time.Time) {
		args := []interface{}{
			slog.String("duration", time.Since(begin).String()),
			slog.Any("topics", topics),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("AuthSubscribe failed to complete successfully", args...)
			return
		}
		lm.logger.Info("AuthSubscribe completed successfully", args...)
	}(time.Now())

	return lm.svc.AuthSubscribe(ctx, topics)
}

// Connect implements session.Handler.
func (lm *loggingMiddleware) Connect(ctx context.Context) (err error) {
	defer func(begin time.Time) {
		args := []interface{}{
			slog.String("duration", time.Since(begin).String()),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Connect failed to complete successfully", args...)
			return
		}
		lm.logger.Info("Connect completed successfully", args...)
	}(time.Now())

	return lm.svc.Connect(ctx)
}

// Disconnect implements session.Handler.
func (lm *loggingMiddleware) Disconnect(ctx context.Context) (err error) {
	defer func(begin time.Time) {
		args := []interface{}{
			slog.String("duration", time.Since(begin).String()),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Disconnect failed to complete successfully", args...)
			return
		}
		lm.logger.Info("Disconnect completed successfully", args...)
	}(time.Now())

	return lm.svc.Disconnect(ctx)
}

// Publish logs the publish request. It logs the time it took to complete the request.
// If the request fails, it logs the error.
func (lm *loggingMiddleware) Publish(ctx context.Context, topic *string, payload *[]byte) (err error) {
	defer func(begin time.Time) {
		args := []interface{}{
			slog.String("duration", time.Since(begin).String()),
			slog.Any("topic", topic),
			slog.Any("payload", payload),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Publish failed to complete successfully", args...)
			return
		}
		lm.logger.Info("Publish completed successfully", args...)
	}(time.Now())

	return lm.svc.Publish(ctx, topic, payload)
}

// Subscribe implements session.Handler.
func (lm *loggingMiddleware) Subscribe(ctx context.Context, topics *[]string) (err error) {
	defer func(begin time.Time) {
		args := []interface{}{
			slog.String("duration", time.Since(begin).String()),
			slog.Any("topics", topics),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Subscribe failed to complete successfully", args...)
			return
		}
		lm.logger.Info("Subscribe completed successfully", args...)
	}(time.Now())

	return lm.svc.Subscribe(ctx, topics)
}

// Unsubscribe implements session.Handler.
func (lm *loggingMiddleware) Unsubscribe(ctx context.Context, topics *[]string) (err error) {
	defer func(begin time.Time) {
		args := []interface{}{
			slog.String("duration", time.Since(begin).String()),
			slog.Any("topics", topics),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Unsubscribe failed to complete successfully", args...)
			return
		}
		lm.logger.Info("Unsubscribe completed successfully", args...)
	}(time.Now())

	return lm.svc.Unsubscribe(ctx, topics)
}

// LoggingMiddleware adds logging facilities to the adapter.
func LoggingMiddleware(svc session.Handler, logger *slog.Logger) session.Handler {
	return &loggingMiddleware{logger, svc}
}
