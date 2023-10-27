// Copyright (c) Magistrala
// SPDX-License-Identifier: Apache-2.0

package domains

import (
	"net/http"

	"github.com/absmach/magistrala/auth"
	"github.com/absmach/magistrala/internal/api"
	"github.com/absmach/magistrala/internal/apiutil"
	"github.com/absmach/magistrala/logger"
	"github.com/go-chi/chi/v5"
	kithttp "github.com/go-kit/kit/transport/http"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func groupsHandler(svc auth.Service, r *chi.Mux, logger logger.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(apiutil.LoggingErrorEncoder(logger, api.EncodeError)),
	}
	r.Route("/domains", func(r chi.Router) {
		r.Post("/", otelhttp.NewHandler(kithttp.NewServer(
			createDomainEndpoint(svc),
			decodeCreateDomainRequest,
			api.EncodeResponse,
			opts...,
		), "create_domain").ServeHTTP)

		r.Get("/{domainID}", otelhttp.NewHandler(kithttp.NewServer(
			viewDomainEndpoint(svc),
			decodeViewDomainRequest,
			api.EncodeResponse,
			opts...,
		), "view_domain").ServeHTTP)

		r.Put("/{domainID}", otelhttp.NewHandler(kithttp.NewServer(
			updateDomainEndpoint(svc),
			decodeUpdateDomainRequest,
			api.EncodeResponse,
			opts...,
		), "update_domain").ServeHTTP)

		r.Get("/", otelhttp.NewHandler(kithttp.NewServer(
			listDomainsEndpoint(svc),
			decodeListDomainRequest,
			api.EncodeResponse,
			opts...,
		), "list_domains").ServeHTTP)

		r.Post("/{domainID}/enable", otelhttp.NewHandler(kithttp.NewServer(
			enableDomainEndpoint(svc),
			decodeEnableDomainRequest,
			api.EncodeResponse,
			opts...,
		), "enable_domain").ServeHTTP)

		r.Post("/{domainID}/disable", otelhttp.NewHandler(kithttp.NewServer(
			disableDomainEndpoint(svc),
			decodeDisableDomainRequest,
			api.EncodeResponse,
			opts...,
		), "disable_domain").ServeHTTP)

		r.Post("/{domainID}/users/assign", otelhttp.NewHandler(kithttp.NewServer(
			assignDomainUsersEndpoint(svc),
			decodeAssignUsersRequest,
			api.EncodeResponse,
			opts...,
		), "assign_domain_users").ServeHTTP)

		r.Post("/{domainID}/users/unassign", otelhttp.NewHandler(kithttp.NewServer(
			unassignDomainUsersEndpoint(svc),
			decodeUnassignUsersRequest,
			api.EncodeResponse,
			opts...,
		), "unassign_domain_users").ServeHTTP)

	})

	r.Get("/users/{userID}/domains", otelhttp.NewHandler(kithttp.NewServer(
		listUsersDomains(svc),
		decodeUsersDomainRequest,
		api.EncodeResponse,
		opts...,
	), "list_channel_by_user_id").ServeHTTP)

	return r
}
