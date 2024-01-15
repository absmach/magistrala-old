// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

//go:build !test

package api

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/absmach/magistrala/pkg/messaging"
	"github.com/absmach/magistrala/twins"
)

var _ twins.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger *slog.Logger
	svc    twins.Service
}

// LoggingMiddleware adds logging facilities to the core service.
func LoggingMiddleware(svc twins.Service, logger *slog.Logger) twins.Service {
	return &loggingMiddleware{logger, svc}
}

func (lm *loggingMiddleware) AddTwin(ctx context.Context, token string, twin twins.Twin, def twins.Definition) (tw twins.Twin, err error) {
	defer func(begin time.Time) {
		message := "Method add_twin completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("method", "add_twin"),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("method", "add_twin"),
			slog.String("id", twin.ID),
			slog.String("token", token),
			slog.String("definition", fmt.Sprintf("%v", def)),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.AddTwin(ctx, token, twin, def)
}

func (lm *loggingMiddleware) UpdateTwin(ctx context.Context, token string, twin twins.Twin, def twins.Definition) (err error) {
	defer func(begin time.Time) {
		message := "Method update_twin completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("method", "update_twin"),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("method", "update_twin"),
			slog.String("id", twin.ID),
			slog.String("token", token),
			slog.String("definition", fmt.Sprintf("%v", def)),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.UpdateTwin(ctx, token, twin, def)
}

func (lm *loggingMiddleware) ViewTwin(ctx context.Context, token, twinID string) (tw twins.Twin, err error) {
	defer func(begin time.Time) {
		message := "Method view_twin completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("method", "view_twin"),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("method", "view_twin"),
			slog.String("id", twinID),
			slog.String("token", token),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.ViewTwin(ctx, token, twinID)
}

func (lm *loggingMiddleware) ListTwins(ctx context.Context, token string, offset, limit uint64, name string, metadata twins.Metadata) (page twins.Page, err error) {
	defer func(begin time.Time) {
		message := "Method list_twins completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("method", "list_twins"),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("method", "list_twins"),
			slog.String("token", token),
			slog.Uint64("offset", offset),
			slog.Uint64("limit", limit),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.ListTwins(ctx, token, offset, limit, name, metadata)
}

func (lm *loggingMiddleware) SaveStates(ctx context.Context, msg *messaging.Message) (err error) {
	defer func(begin time.Time) {
		message := "Method save_states completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("method", "save_states"),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("method", "save_states"),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.SaveStates(ctx, msg)
}

func (lm *loggingMiddleware) ListStates(ctx context.Context, token string, offset, limit uint64, twinID string) (page twins.StatesPage, err error) {
	defer func(begin time.Time) {
		message := "Method list_states completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("method", "list_states"),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("method", "list_states"),
			slog.String("token", token),
			slog.Uint64("offset", offset),
			slog.Uint64("limit", limit),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.ListStates(ctx, token, offset, limit, twinID)
}

func (lm *loggingMiddleware) RemoveTwin(ctx context.Context, token, twinID string) (err error) {
	defer func(begin time.Time) {
		message := "Method remove_twin completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("method", "remove_twin"),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("method", "remove_twin"),
			slog.String("token", token),
			slog.String("twin_ID", twinID),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.RemoveTwin(ctx, token, twinID)
}
