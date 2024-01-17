// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

//go:build !test

package api

import (
	"context"
	"fmt"
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
		message := "Method create_thing completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error", message),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors", message),
			slog.String("thing_id", mgxThing),
			slog.String("opcua_node_id", opcuaNodeID),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.CreateThing(ctx, mgxThing, opcuaNodeID)
}

func (lm loggingMiddleware) UpdateThing(ctx context.Context, mgxThing, opcuaNodeID string) (err error) {
	defer func(begin time.Time) {
		message := "Method update_thing completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error", message),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors", message),
			slog.String("thing_id", mgxThing),
			slog.String("opcua_node_id", opcuaNodeID),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.UpdateThing(ctx, mgxThing, opcuaNodeID)
}

func (lm loggingMiddleware) RemoveThing(ctx context.Context, mgxThing string) (err error) {
	defer func(begin time.Time) {
		message := "Method remove_thing completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error", message),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors", message),
			slog.String("thing_id", mgxThing),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.RemoveThing(ctx, mgxThing)
}

func (lm loggingMiddleware) CreateChannel(ctx context.Context, mgxChan, opcuaServerURI string) (err error) {
	defer func(begin time.Time) {
		message := "Method create_channel completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error", message),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors", message),
			slog.String("channel_id", mgxChan),
			slog.String("opcua_server_uri", opcuaServerURI),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.CreateChannel(ctx, mgxChan, opcuaServerURI)
}

func (lm loggingMiddleware) UpdateChannel(ctx context.Context, mgxChanID, opcuaServerURI string) (err error) {
	defer func(begin time.Time) {
		message := "Method update_channel completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error", message),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors", message),
			slog.String("channel_id", mgxChanID),
			slog.String("opcua_server_uri", opcuaServerURI),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.UpdateChannel(ctx, mgxChanID, opcuaServerURI)
}

func (lm loggingMiddleware) RemoveChannel(ctx context.Context, mgxChanID string) (err error) {
	defer func(begin time.Time) {
		message := "Method remove_channel completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error", message),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors", message),
			slog.String("channel_id", mgxChanID),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.RemoveChannel(ctx, mgxChanID)
}

func (lm loggingMiddleware) ConnectThing(ctx context.Context, mgxChanID, mgxThingID string) (err error) {
	defer func(begin time.Time) {
		message := "Method connect_thing completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error", message),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors", message),
			slog.String("channel_id", mgxChanID),
			slog.String("thing_id", mgxThingID),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.ConnectThing(ctx, mgxChanID, mgxThingID)
}

func (lm loggingMiddleware) DisconnectThing(ctx context.Context, mgxChanID, mgxThingID string) (err error) {
	defer func(begin time.Time) {
		message := "Method disconnect_thing completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error", message),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors", message),
			slog.String("channel_id", mgxChanID),
			slog.String("thing_id", mgxThingID),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.DisconnectThing(ctx, mgxChanID, mgxThingID)
}

func (lm loggingMiddleware) Browse(ctx context.Context, serverURI, namespace, identifier string) (nodes []opcua.BrowsedNode, err error) {
	defer func(begin time.Time) {
		message := "Method browse completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error", message),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors", message),
			slog.String("server_uri", serverURI),
			slog.String("namespace", namespace),
			slog.String("identifier", identifier),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.Browse(ctx, serverURI, namespace, identifier)
}
