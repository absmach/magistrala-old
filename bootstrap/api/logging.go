// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

//go:build !test

package api

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/absmach/magistrala/bootstrap"
)

var _ bootstrap.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger *slog.Logger
	svc    bootstrap.Service
}

// LoggingMiddleware adds logging facilities to the bootstrap service.
func LoggingMiddleware(svc bootstrap.Service, logger *slog.Logger) bootstrap.Service {
	return &loggingMiddleware{logger, svc}
}

// Add logs the add request. It logs the thing ID and the time it took to complete the request.
// If the request fails, it logs the error.
func (lm *loggingMiddleware) Add(ctx context.Context, token string, cfg bootstrap.Config) (saved bootstrap.Config, err error) {
	defer func(begin time.Time) {
		message :="Method add completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("thing_id", saved.ThingID),
			slog.String("token", token),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.Add(ctx, token, cfg)
}

// View logs the view request. It logs the thing ID and the time it took to complete the request.
// If the request fails, it logs the error.
func (lm *loggingMiddleware) View(ctx context.Context, token, id string) (saved bootstrap.Config, err error) {
	defer func(begin time.Time) {
		message :="Method view completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("thing_id", saved.ThingID),
			slog.String("token", token),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.View(ctx, token, id)
}

// Update logs the update request. It logs token, bootstrap thing ID and the time it took to complete the request.
// If the request fails, it logs the error.
func (lm *loggingMiddleware) Update(ctx context.Context, token string, cfg bootstrap.Config) (err error) {
	defer func(begin time.Time) {
		message :="Method update completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error: %s.", message, err),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("thing_id", cfg.ThingID),
			slog.String("token", token),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.Update(ctx, token, cfg)
}

// UpdateCert logs the update_cert request. It logs token, thing ID and the time it took to complete the request.
// If the request fails, it logs the error.
func (lm *loggingMiddleware) UpdateCert(ctx context.Context, token, thingID, clientCert, clientKey, caCert string) (cfg bootstrap.Config, err error) {
	defer func(begin time.Time) {	
		message:= "Method update_cert completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("thing_id", thingID),
			slog.String("token", token),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.UpdateCert(ctx, token, thingID, clientCert, clientKey, caCert)
}

// UpdateConnections logs the update_connections request. It logs token, bootstrap ID and the time it took to complete the request.
// If the request fails, it logs the error.
func (lm *loggingMiddleware) UpdateConnections(ctx context.Context, token, id string, connections []string) (err error) {
	defer func(begin time.Time) {
		message :="Method update_connections completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error: %s.", message, err),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("thing_id", id),
			slog.String("token", token),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.UpdateConnections(ctx, token, id, connections)
}

// List logs the list request. It logs token, offset, limit and the time it took to complete the request.
// If the request fails, it logs the error.
func (lm *loggingMiddleware) List(ctx context.Context, token string, filter bootstrap.Filter, offset, limit uint64) (res bootstrap.ConfigsPage, err error) {
	defer func(begin time.Time) {
		message :="Method list completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("token", token),
			slog.Uint64("offset", offset),
			slog.Uint64("limit", limit),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.List(ctx, token, filter, offset, limit)
}

// Remove logs the remove request. It logs token, bootstrap ID and the time it took to complete the request.
// If the request fails, it logs the error.
func (lm *loggingMiddleware) Remove(ctx context.Context, token, id string) (err error) {
	defer func(begin time.Time) {
		message :="Method remove completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error: %s.", message, err),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("thing_id", id),
			slog.String("token", token),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.Remove(ctx, token, id)
}

func (lm *loggingMiddleware) Bootstrap(ctx context.Context, externalKey, externalID string, secure bool) (cfg bootstrap.Config, err error) {
	defer func(begin time.Time) {
		message :="Method bootstrap completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("external_id", externalID),
			slog.String("external_key", externalKey),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.Bootstrap(ctx, externalKey, externalID, secure)
}

func (lm *loggingMiddleware) ChangeState(ctx context.Context, token, id string, state bootstrap.State) (err error) {
	defer func(begin time.Time) {
		message :="Method change_state completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("thing_id", id),
			slog.String("token", token),
			slog.String("state", state.String()),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.ChangeState(ctx, token, id, state)
}

func (lm *loggingMiddleware) UpdateChannelHandler(ctx context.Context, channel bootstrap.Channel) (err error) {
	defer func(begin time.Time) {
		message :="Method update_channel_handler completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("channel_id", channel.ID),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.UpdateChannelHandler(ctx, channel)
}

func (lm *loggingMiddleware) RemoveConfigHandler(ctx context.Context, id string) (err error) {
	defer func(begin time.Time) {
		message :="Method remove_config_handler completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error: %s.", message, err),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("config_id", id),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.RemoveConfigHandler(ctx, id)
}

func (lm *loggingMiddleware) RemoveChannelHandler(ctx context.Context, id string) (err error) {
	defer func(begin time.Time) {
		message :="Method remove_channel_handler completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("channel_id", id),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.RemoveChannelHandler(ctx, id)
}

func (lm *loggingMiddleware) DisconnectThingHandler(ctx context.Context, channelID, thingID string) (err error) {
	defer func(begin time.Time) {
		message :="Method disconnect_thing_handler completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("channel_id", channelID),
			slog.String("thing_id", thingID),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.DisconnectThingHandler(ctx, channelID, thingID)
}
