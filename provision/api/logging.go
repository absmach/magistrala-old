// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

//go:build !test

package api

import (
	"fmt"
	"time"

	mglog "github.com/absmach/magistrala/logger"
	"github.com/absmach/magistrala/provision"
)

var _ provision.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger mglog.Logger
	svc    provision.Service
}

// NewLoggingMiddleware adds logging facilities to the core service.
func NewLoggingMiddleware(svc provision.Service, logger mglog.Logger) provision.Service {
	return &loggingMiddleware{logger, svc}
}

func (lm *loggingMiddleware) Provision(token, name, externalID, externalKey string) (res provision.Result, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method provision for token: %s and things: %v took %s to complete", token, res.Things, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors", message))
	}(time.Now())

	return lm.svc.Provision(token, name, externalID, externalKey)
}

func (lm *loggingMiddleware) Cert(token, thingID, duration string) (cert string, key string, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method cert for token: %s and thing: %v took %s to complete", token, thingID, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors", message))
	}(time.Now())

	return lm.svc.Cert(token, thingID, duration)
}

func (lm *loggingMiddleware) Mapping(token string) (res map[string]interface{}, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method mapping for token: %s took %s to complete", token, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors", message))
	}(time.Now())

	return lm.svc.Mapping(token)
}
