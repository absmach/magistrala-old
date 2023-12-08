// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/absmach/magistrala"
	"github.com/absmach/magistrala/internal/api"
	"github.com/absmach/magistrala/internal/apiutil"
	"github.com/absmach/magistrala/invitations"
	mglog "github.com/absmach/magistrala/logger"
	"github.com/absmach/magistrala/pkg/errors"
	"github.com/go-chi/chi/v5"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const (
	userIDKey    = "user_id"
	domainKey    = "domain"
	invitedByKey = "invited_by"
	relationKey  = "relation"
)

func MakeHandler(svc invitations.Service, logger mglog.Logger, instanceID string) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(apiutil.LoggingErrorEncoder(logger, api.EncodeError)),
	}

	mux := chi.NewRouter()

	mux.Route("/invitations", func(r chi.Router) {
		r.Post("/", otelhttp.NewHandler(kithttp.NewServer(
			sendInvitationEndpoint(svc),
			decodeSendInvitationReq,
			api.EncodeResponse,
			opts...,
		), "send_invitation").ServeHTTP)
		r.Get("/", otelhttp.NewHandler(kithttp.NewServer(
			listInvitationsEndpoint(svc),
			decodeListInvitationsReq,
			api.EncodeResponse,
			opts...,
		), "list_invitations").ServeHTTP)
		r.Route("/{user_id}/{domain}", func(r chi.Router) {
			r.Get("/", otelhttp.NewHandler(kithttp.NewServer(
				viewInvitationEndpoint(svc),
				decodeInvitationReq,
				api.EncodeResponse,
				opts...,
			), "view_invitations").ServeHTTP)
			r.Delete("/", otelhttp.NewHandler(kithttp.NewServer(
				deleteInvitationEndpoint(svc),
				decodeInvitationReq,
				api.EncodeResponse,
				opts...,
			), "delete_invitation").ServeHTTP)
		})
		r.Post("/accept", otelhttp.NewHandler(kithttp.NewServer(
			acceptInvitationEndpoint(svc),
			decodeAcceptInvitationReq,
			api.EncodeResponse,
			opts...,
		), "accept_invitation").ServeHTTP)
	})

	mux.Get("/health", magistrala.Health("invitations", instanceID))
	mux.Handle("/metrics", promhttp.Handler())

	return mux
}

func decodeSendInvitationReq(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), api.ContentType) {
		return nil, errors.Wrap(apiutil.ErrValidation, apiutil.ErrUnsupportedContentType)
	}

	var req sendInvitationReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(apiutil.ErrValidation, errors.Wrap(err, errors.ErrMalformedEntity))
	}
	req.token = apiutil.ExtractBearerToken(r)

	return req, nil
}

func decodeListInvitationsReq(_ context.Context, r *http.Request) (interface{}, error) {
	offset, err := apiutil.ReadNumQuery[uint64](r, api.OffsetKey, api.DefOffset)
	if err != nil {
		return nil, errors.Wrap(apiutil.ErrValidation, err)
	}
	limit, err := apiutil.ReadNumQuery[uint64](r, api.LimitKey, api.DefLimit)
	if err != nil {
		return nil, errors.Wrap(apiutil.ErrValidation, err)
	}
	userID, err := apiutil.ReadStringQuery(r, userIDKey, "")
	if err != nil {
		return nil, errors.Wrap(apiutil.ErrValidation, err)
	}
	invitedBy, err := apiutil.ReadStringQuery(r, invitedByKey, "")
	if err != nil {
		return nil, errors.Wrap(apiutil.ErrValidation, err)
	}
	relation, err := apiutil.ReadStringQuery(r, relationKey, "")
	if err != nil {
		return nil, errors.Wrap(apiutil.ErrValidation, err)
	}
	domain, err := apiutil.ReadStringQuery(r, domainKey, "")
	if err != nil {
		return nil, errors.Wrap(apiutil.ErrValidation, err)
	}

	req := listInvitationsReq{
		token: apiutil.ExtractBearerToken(r),
		Page: invitations.Page{
			Offset:    offset,
			Limit:     limit,
			InvitedBy: invitedBy,
			UserID:    userID,
			Relation:  relation,
			Domain:    domain,
		},
	}

	return req, nil
}

func decodeAcceptInvitationReq(_ context.Context, r *http.Request) (interface{}, error) {
	return acceptInvitationReq{
		token: apiutil.ExtractBearerToken(r),
	}, nil
}

func decodeInvitationReq(_ context.Context, r *http.Request) (interface{}, error) {
	req := invitationReq{
		token:  apiutil.ExtractBearerToken(r),
		userID: chi.URLParam(r, "user_id"),
		domain: chi.URLParam(r, "domain"),
	}

	return req, nil
}
