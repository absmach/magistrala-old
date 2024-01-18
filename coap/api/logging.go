// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

//go:build !test

package api

import (
	"context"
	"fmt"
	"time"

	"github.com/absmach/magistrala/coap"
	mglog "github.com/absmach/magistrala/logger"
	"github.com/absmach/magistrala/pkg/messaging"
)

var _ coap.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger mglog.Logger
	svc    coap.Service
}

// LoggingMiddleware adds logging facilities to the adapter.
func LoggingMiddleware(svc coap.Service, logger mglog.Logger) coap.Service {
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
		message := fmt.Sprintf("Method publish to %s took %s to complete", destChannel, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
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
		message := fmt.Sprintf("Method subscribe to %s for client %s took %s to complete", destChannel, c.Token(), time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
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
		message := fmt.Sprintf("Method unsubscribe for the client %s from the channel %s took %s to complete", token, destChannel, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.Unsubscribe(ctx, key, chanID, subtopic, token)
}
