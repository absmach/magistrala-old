// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

//go:build !test

package api

import (
	"context"
	"log/slog"
	"time"

	"github.com/absmach/magistrala/consumers/notifiers"
)

var _ notifiers.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger *slog.Logger
	svc    notifiers.Service
}

// LoggingMiddleware adds logging facilities to the core service.
func LoggingMiddleware(svc notifiers.Service, logger *slog.Logger) notifiers.Service {
	return &loggingMiddleware{logger, svc}
}

// CreateSubscription logs the create_subscription request. It logs subscription ID and topic and the time it took to complete the request.
// If the request fails, it logs the error.
func (lm *loggingMiddleware) CreateSubscription(ctx context.Context, token string, sub notifiers.Subscription) (id string, err error) {
	defer func(begin time.Time) {
		args := []interface{}{
			slog.String("duration", time.Since(begin).String()),
			slog.Group(
				"subscription",
				slog.String("topic", sub.Topic),
				slog.String("id", id),
			),
		}
		if err != nil {
			args = append(args, slog.String("error", err.Error()))
			lm.logger.Warn("Create subscription failed to complete successfully", args...)
			return
		}
		lm.logger.Info("Create subscription completed successfully", args...)
	}(time.Now())

	return lm.svc.CreateSubscription(ctx, token, sub)
}

// ViewSubscription logs the view_subscription request. It logs subscription topic and id and the time it took to complete the request.
// If the request fails, it logs the error.
func (lm *loggingMiddleware) ViewSubscription(ctx context.Context, token, topic string) (sub notifiers.Subscription, err error) {
	defer func(begin time.Time) {
		args := []interface{}{
			slog.String("duration", time.Since(begin).String()),
			slog.Group(
				"subscription",
				slog.String("topic", sub.Topic),
				slog.String("id", sub.ID),
			),
		}
		if err != nil {
			args = append(args, slog.String("error", err.Error()))
			lm.logger.Warn("View subscription failed to complete successfully", args...)
			return
		}
		lm.logger.Info("View subscription completed successfully", args...)
	}(time.Now())

	return lm.svc.ViewSubscription(ctx, token, topic)
}

// ListSubscriptions logs the list_subscriptions request. It logs page metadata and subscription topic and the time it took to complete the request.
// If the request fails, it logs the error.
func (lm *loggingMiddleware) ListSubscriptions(ctx context.Context, token string, pm notifiers.PageMetadata) (res notifiers.Page, err error) {
	defer func(begin time.Time) {
		args := []interface{}{
			slog.String("duration", time.Since(begin).String()),
			slog.Group(
				"page_metadata",
				slog.String("topic", pm.Topic),
				slog.Any("limit", pm.Limit),
				slog.Any("offset", pm.Offset),
			),
		}
		if err != nil {
			args = append(args, slog.String("error", err.Error()))
			lm.logger.Warn("List subscriptions failed to complete successfully", args...)
			return
		}
		lm.logger.Info("List subscriptions completed successfully", args...)
	}(time.Now())

	return lm.svc.ListSubscriptions(ctx, token, pm)
}

// RemoveSubscription logs the remove_subscription request. It logs subscription ID and the time it took to complete the request.
// If the request fails, it logs the error.
func (lm *loggingMiddleware) RemoveSubscription(ctx context.Context, token, id string) (err error) {
	defer func(begin time.Time) {
		args := []interface{}{
			slog.String("duration", time.Since(begin).String()),
			slog.String("id", id),
		}
		if err != nil {
			args = append(args, slog.String("error", err.Error()))
			lm.logger.Warn("Remove subscription failed to complete successfully", args...)
			return
		}
		lm.logger.Info("Remove subscription completed successfully", args...)
	}(time.Now())

	return lm.svc.RemoveSubscription(ctx, token, id)
}

// ConsumeBlocking logs the consume_blocking request. It logs the time it took to complete the request.
// If the request fails, it logs the error.
func (lm *loggingMiddleware) ConsumeBlocking(ctx context.Context, msg interface{}) (err error) {
	defer func(begin time.Time) {
		args := []interface{}{
			slog.String("duration", time.Since(begin).String()),
		}
		if err != nil {
			args = append(args, slog.String("error", err.Error()))
			lm.logger.Warn("Consume blocking failed to complete successfully", args...)
			return
		}
		lm.logger.Info("Consume blocking completed successfully", args...)
	}(time.Now())

	return lm.svc.ConsumeBlocking(ctx, msg)
}
