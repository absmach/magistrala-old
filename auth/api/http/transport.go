// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0
package http

import (
	"net/http"

	"github.com/absmach/magistrala"
	"github.com/absmach/magistrala/auth"
	"github.com/absmach/magistrala/auth/api/http/domains"
	"github.com/absmach/magistrala/auth/api/http/keys"
	"github.com/absmach/magistrala/logger"
	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MakeHandler returns a HTTP handler for API endpoints.
func MakeHandler(svc auth.Service, logger logger.Logger, instanceID string) http.Handler {
	r := chi.NewRouter()

	mux = keys.MakeHandler(svc, mux, logger)
	mux = domains.MakeHandler(svc, mux, logger)

	r.Get("/health", magistrala.Health("auth", instanceID))
	r.Handle("/metrics", promhttp.Handler())

	return r
}
