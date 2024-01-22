// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

//go:build !test

package api

import (
	"context"
	"log/slog"
	"time"

	"github.com/absmach/magistrala/opcua"
)

var _ opcua.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger *slog.Logger
	svc    opcua.Service
}

// LoggingMiddleware adds logging facilities to the core service.
func LoggingMiddleware(svc opcua.Service, logger *slog.Logger) opcua.Service {
	return &loggingMiddleware{
		logger: logger,
		svc:    svc,
	}
}

func (lm loggingMiddleware) CreateThing(ctx context.Context, mgxThing, opcuaNodeID string) (err error) {
	defer func(begin time.Time) {
		args := []interface{}{
			slog.String("duration", time.Since(begin).String()),
			slog.Group(
				"thing",
				slog.String("thing_id", mgxThing),
				slog.String("opcua_node_id", opcuaNodeID),
			),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Create thing failed to complete successfully", args...)
			return
		}
		lm.logger.Info("Create thing completed successfully", args...)
	}(time.Now())

	return lm.svc.CreateThing(ctx, mgxThing, opcuaNodeID)
}

func (lm loggingMiddleware) UpdateThing(ctx context.Context, mgxThing, opcuaNodeID string) (err error) {
	defer func(begin time.Time) {
		args := []interface{}{
			slog.String("duration", time.Since(begin).String()),
			slog.Group(
				"thing",
				slog.String("thing_id", mgxThing),
				slog.String("opcua_node_id", opcuaNodeID),
			),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Update thing failed to complete successfully", args...)
			return
		}
		lm.logger.Info("Update thing completed successfully", args...)
	}(time.Now())

	return lm.svc.UpdateThing(ctx, mgxThing, opcuaNodeID)
}

func (lm loggingMiddleware) RemoveThing(ctx context.Context, mgxThing string) (err error) {
	defer func(begin time.Time) {
		args := []interface{}{
			slog.String("duration", time.Since(begin).String()),
			slog.String("thing_id", mgxThing),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Remove thing failed to complete successfully", args...)
			return
		}
		lm.logger.Info("Remove thing completed successfully", args...)
	}(time.Now())

	return lm.svc.RemoveThing(ctx, mgxThing)
}

func (lm loggingMiddleware) CreateChannel(ctx context.Context, mgxChan, opcuaServerURI string) (err error) {
	defer func(begin time.Time) {
		args := []interface{}{
			slog.String("duration", time.Since(begin).String()),
			slog.Group(
				"channel",
				slog.String("channel_id", mgxChan),
				slog.String("opcua_server_uri", opcuaServerURI),
			),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Create channel failed to complete successfully", args...)
			return
		}
		lm.logger.Info("Create channel completed successfully", args...)
	}(time.Now())

	return lm.svc.CreateChannel(ctx, mgxChan, opcuaServerURI)
}

func (lm loggingMiddleware) UpdateChannel(ctx context.Context, mgxChanID, opcuaServerURI string) (err error) {
	defer func(begin time.Time) {
		args := []interface{}{
			slog.String("duration", time.Since(begin).String()),
			slog.Group(
				"channel",
				slog.String("channel_id", mgxChanID),
				slog.String("opcua_server_uri", opcuaServerURI),
			),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Update channel failed to complete successfully", args...)
			return
		}
		lm.logger.Info("Update channel completed successfully", args...)
	}(time.Now())

	return lm.svc.UpdateChannel(ctx, mgxChanID, opcuaServerURI)
}

func (lm loggingMiddleware) RemoveChannel(ctx context.Context, mgxChanID string) (err error) {
	defer func(begin time.Time) {
		args := []interface{}{
			slog.String("duration", time.Since(begin).String()),
			slog.String("channel_id", mgxChanID),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Remove channel failed to complete successfully", args...)
			return
		}
		lm.logger.Info("Remove channel completed successfully", args...)
	}(time.Now())

	return lm.svc.RemoveChannel(ctx, mgxChanID)
}

func (lm loggingMiddleware) ConnectThing(ctx context.Context, mgxChanID, mgxThingID string) (err error) {
	defer func(begin time.Time) {
		args := []interface{}{
			slog.String("duration", time.Since(begin).String()),
			slog.String("channel_id", mgxChanID),
			slog.String("thing_id", mgxThingID),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Connect thing failed to complete successfully", args...)
			return
		}
		lm.logger.Info("Connect thing completed successfully", args...)
	}(time.Now())

	return lm.svc.ConnectThing(ctx, mgxChanID, mgxThingID)
}

func (lm loggingMiddleware) DisconnectThing(ctx context.Context, mgxChanID, mgxThingID string) (err error) {
	defer func(begin time.Time) {
		args := []interface{}{
			slog.String("duration", time.Since(begin).String()),
			slog.String("channel_id", mgxChanID),
			slog.String("thing_id", mgxThingID),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Disconnect thing failed to complete successfully", args...)
			return
		}
		lm.logger.Info("Disconnect thing completed successfully", args...)
	}(time.Now())

	return lm.svc.DisconnectThing(ctx, mgxChanID, mgxThingID)
}

func (lm loggingMiddleware) Browse(ctx context.Context, serverURI, namespace, identifier string) (nodes []opcua.BrowsedNode, err error) {
	defer func(begin time.Time) {
		args := []interface{}{
			slog.String("duration", time.Since(begin).String()),
			slog.String("server_uri", serverURI),
			slog.String("namespace", namespace),
			slog.String("identifier", identifier),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Browse failed to complete successfully", args...)
			return
		}
		lm.logger.Info("Browse completed successfully", args...)
	}(time.Now())

	return lm.svc.Browse(ctx, serverURI, namespace, identifier)
}
