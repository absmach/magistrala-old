// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

//go:build !test

package api

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/absmach/magistrala/lora"
)

var _ lora.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger *slog.Logger
	svc    lora.Service
}

// LoggingMiddleware adds logging facilities to the core service.
func LoggingMiddleware(svc lora.Service, logger *slog.Logger) lora.Service {
	return &loggingMiddleware{
		logger: logger,
		svc:    svc,
	}
}

func (lm loggingMiddleware) CreateThing(ctx context.Context, thingID, loraDevEUI string) (err error) {
	defer func(begin time.Time) {
		message := "Method create_thing completed"
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
			slog.String("lora_dev_eui", loraDevEUI),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.CreateThing(ctx, thingID, loraDevEUI)
}

func (lm loggingMiddleware) UpdateThing(ctx context.Context, thingID, loraDevEUI string) (err error) {
	defer func(begin time.Time) {
		message := "Method update_thing completed"
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
			slog.String("lora_dev_eui", loraDevEUI),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.UpdateThing(ctx, thingID, loraDevEUI)
}

func (lm loggingMiddleware) RemoveThing(ctx context.Context, thingID string) (err error) {
	defer func(begin time.Time) {
		message := "Method remove_thing completed"
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
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.RemoveThing(ctx, thingID)
}

func (lm loggingMiddleware) CreateChannel(ctx context.Context, chanID, loraApp string) (err error) {
	defer func(begin time.Time) {
		message := "Method create_channel completed"
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
			slog.String("channel_id", chanID),
			slog.String("lora_app", loraApp),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.CreateChannel(ctx, chanID, loraApp)
}

func (lm loggingMiddleware) UpdateChannel(ctx context.Context, chanID, loraApp string) (err error) {
	defer func(begin time.Time) {
		message := "Method update_channel completed"
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
			slog.String("channel_id", chanID),
			slog.String("lora_app", loraApp),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.UpdateChannel(ctx, chanID, loraApp)
}

func (lm loggingMiddleware) RemoveChannel(ctx context.Context, chanID string) (err error) {
	defer func(begin time.Time) {
		message := "Method remove_channel completed"
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
			slog.String("channel_id", chanID),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.RemoveChannel(ctx, chanID)
}

func (lm loggingMiddleware) ConnectThing(ctx context.Context, chanID, thingID string) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method connect_thing mgx-%s : mgx-%s, took %s to complete", chanID, thingID, time.Since(begin))
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
			slog.String("channel_id", chanID),
			slog.String("thing_id", thingID),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.ConnectThing(ctx, chanID, thingID)
}

func (lm loggingMiddleware) DisconnectThing(ctx context.Context, chanID, thingID string) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method disconnect_thing mgx-%s : mgx-%s, took %s to complete", chanID, thingID, time.Since(begin))
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
			slog.String("channel_id", chanID),
			slog.String("thing_id", thingID),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.DisconnectThing(ctx, chanID, thingID)
}

func (lm loggingMiddleware) Publish(ctx context.Context, msg *lora.Message) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method publish application/%s/device/%s/rx took %s to complete", msg.ApplicationID, msg.DevEUI, time.Since(begin))
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
			slog.String("application_id", msg.ApplicationID),
			slog.String("device_eui", msg.DevEUI),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.Publish(ctx, msg)
}
