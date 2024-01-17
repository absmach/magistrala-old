// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

//go:build !test

package api

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/absmach/magistrala/provision"
)

var _ provision.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger *slog.Logger
	svc    provision.Service
}

// NewLoggingMiddleware adds logging facilities to the core service.
func NewLoggingMiddleware(svc provision.Service, logger *slog.Logger) provision.Service {
	return &loggingMiddleware{logger, svc}
}

func (lm *loggingMiddleware) Provision(token, name, externalID, externalKey string) (res provision.Result, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method provision for things: %v took completed", res.Things)
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors", message),
			slog.String("token", token),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.Provision(token, name, externalID, externalKey)
}

func (lm *loggingMiddleware) Cert(token, thingID, duration string) (cert, key string, err error) {
	defer func(begin time.Time) {
		message := "Method cert completed"
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
			slog.String("token", token),
			slog.String("thing_id", thingID),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.Cert(token, thingID, duration)
}

func (lm *loggingMiddleware) Mapping(token string) (res map[string]interface{}, err error) {
	defer func(begin time.Time) {
		message := "Method mapping completed"
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
			slog.String("token", token),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.Mapping(token)
}
