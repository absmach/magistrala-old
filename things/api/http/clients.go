// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/absmach/magistrala/internal/api"
	"github.com/absmach/magistrala/internal/apiutil"
	mglog "github.com/absmach/magistrala/logger"
	mgclients "github.com/absmach/magistrala/pkg/clients"
	"github.com/absmach/magistrala/pkg/errors"
	"github.com/absmach/magistrala/things"
	"github.com/go-chi/chi/v5"
	kithttp "github.com/go-kit/kit/transport/http"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func clientsHandler(svc things.Service, r *chi.Mux, logger mglog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(apiutil.LoggingErrorEncoder(logger, api.EncodeError)),
	}
	r.Route("/things", func(r chi.Router) {
		r.Post("/", otelhttp.NewHandler(kithttp.NewServer(
			createClientEndpoint(svc),
			decodeCreateClientReq,
			api.EncodeResponse,
			opts...,
		), "create_thing").ServeHTTP)

		r.Get("/", otelhttp.NewHandler(kithttp.NewServer(
			listClientsEndpoint(svc),
			decodeListClients,
			api.EncodeResponse,
			opts...,
		), "list_things").ServeHTTP)

		r.Post("/bulk", otelhttp.NewHandler(kithttp.NewServer(
			createClientsEndpoint(svc),
			decodeCreateClientsReq,
			api.EncodeResponse,
			opts...,
		), "create_things").ServeHTTP)

		r.Get("/{thingID}", otelhttp.NewHandler(kithttp.NewServer(
			viewClientEndpoint(svc),
			decodeViewClient,
			api.EncodeResponse,
			opts...,
		), "view_thing").ServeHTTP)

		r.Get("/{thingID}/permissions", otelhttp.NewHandler(kithttp.NewServer(
			viewClientPermsEndpoint(svc),
			decodeViewClientPerms,
			api.EncodeResponse,
			opts...,
		), "view_thing").ServeHTTP)

		r.Patch("/{thingID}", otelhttp.NewHandler(kithttp.NewServer(
			updateClientEndpoint(svc),
			decodeUpdateClient,
			api.EncodeResponse,
			opts...,
		), "update_thing").ServeHTTP)

		r.Patch("/{thingID}/tags", otelhttp.NewHandler(kithttp.NewServer(
			updateClientTagsEndpoint(svc),
			decodeUpdateClientTags,
			api.EncodeResponse,
			opts...,
		), "update_thing_tags").ServeHTTP)

		r.Patch("/{thingID}/secret", otelhttp.NewHandler(kithttp.NewServer(
			updateClientSecretEndpoint(svc),
			decodeUpdateClientCredentials,
			api.EncodeResponse,
			opts...,
		), "update_thing_credentials").ServeHTTP)

		r.Post("/{thingID}/enable", otelhttp.NewHandler(kithttp.NewServer(
			enableClientEndpoint(svc),
			decodeChangeClientStatus,
			api.EncodeResponse,
			opts...,
		), "enable_thing").ServeHTTP)

		r.Post("/{thingID}/disable", otelhttp.NewHandler(kithttp.NewServer(
			disableClientEndpoint(svc),
			decodeChangeClientStatus,
			api.EncodeResponse,
			opts...,
		), "disable_thing").ServeHTTP)

		r.Post("/{thingID}/share", otelhttp.NewHandler(kithttp.NewServer(
			thingShareEndpoint(svc),
			decodeThingShareRequest,
			api.EncodeResponse,
			opts...,
		), "thing_share").ServeHTTP)

		r.Post("/{thingID}/unshare", otelhttp.NewHandler(kithttp.NewServer(
			thingUnshareEndpoint(svc),
			decodeThingUnshareRequest,
			api.EncodeResponse,
			opts...,
		), "thing_delete_share").ServeHTTP)
	})

	// Ideal location: things service,  channels endpoint
	// Reason for placing here :
	// SpiceDB provides list of thing ids present in given channel id
	// and things service can access spiceDB and get the list of thing ids present in given channel id.
	// Request to get list of things present in channelID ({groupID}) .
	r.Get("/channels/{groupID}/things", otelhttp.NewHandler(kithttp.NewServer(
		listMembersEndpoint(svc),
		decodeListMembersRequest,
		api.EncodeResponse,
		opts...,
	), "list_things_by_channel_id").ServeHTTP)

	r.Get("/users/{userID}/things", otelhttp.NewHandler(kithttp.NewServer(
		listClientsEndpoint(svc),
		decodeListClients,
		api.EncodeResponse,
		opts...,
	), "list_user_things").ServeHTTP)
	return r
}

func decodeViewClient(_ context.Context, r *http.Request) (interface{}, error) {
	req := viewClientReq{
		token: apiutil.ExtractBearerToken(r),
		id:    chi.URLParam(r, "thingID"),
	}

	return req, nil
}

func decodeViewClientPerms(_ context.Context, r *http.Request) (interface{}, error) {
	req := viewClientPermsReq{
		token: apiutil.ExtractBearerToken(r),
		id:    chi.URLParam(r, "thingID"),
	}

	return req, nil
}

func decodeListClients(_ context.Context, r *http.Request) (interface{}, error) {
	var ownerID string
	s, err := apiutil.ReadStringQuery(r, api.StatusKey, api.DefClientStatus)
	if err != nil {
		return nil, errors.Wrap(apiutil.ErrValidation, err)
	}
	o, err := apiutil.ReadNumQuery[uint64](r, api.OffsetKey, api.DefOffset)
	if err != nil {
		return nil, errors.Wrap(apiutil.ErrValidation, err)
	}
	l, err := apiutil.ReadNumQuery[uint64](r, api.LimitKey, api.DefLimit)
	if err != nil {
		return nil, errors.Wrap(apiutil.ErrValidation, err)
	}
	m, err := apiutil.ReadMetadataQuery(r, api.MetadataKey, nil)
	if err != nil {
		return nil, errors.Wrap(apiutil.ErrValidation, err)
	}
	n, err := apiutil.ReadStringQuery(r, api.NameKey, "")
	if err != nil {
		return nil, errors.Wrap(apiutil.ErrValidation, err)
	}
	t, err := apiutil.ReadStringQuery(r, api.TagKey, "")
	if err != nil {
		return nil, errors.Wrap(apiutil.ErrValidation, err)
	}
	oid, err := apiutil.ReadStringQuery(r, api.OwnerKey, "")
	if err != nil {
		return nil, err
	}

	p, err := apiutil.ReadStringQuery(r, api.PermissionKey, api.DefPermission)
	if err != nil {
		return nil, errors.Wrap(apiutil.ErrValidation, err)
	}

	lp, err := apiutil.ReadBoolQuery(r, api.ListPerms, api.DefListPerms)
	if err != nil {
		return nil, errors.Wrap(apiutil.ErrValidation, err)
	}

	if oid != "" {
		ownerID = oid
	}
	st, err := mgclients.ToStatus(s)
	if err != nil {
		return nil, errors.Wrap(apiutil.ErrValidation, err)
	}
	req := listClientsReq{
		token:      apiutil.ExtractBearerToken(r),
		status:     st,
		offset:     o,
		limit:      l,
		metadata:   m,
		name:       n,
		tag:        t,
		permission: p,
		listPerms:  lp,
		userID:     chi.URLParam(r, "userID"),
		owner:      ownerID,
	}
	return req, nil
}

func decodeUpdateClient(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), api.ContentType) {
		return nil, errors.Wrap(apiutil.ErrValidation, apiutil.ErrUnsupportedContentType)
	}

	req := updateClientReq{
		token: apiutil.ExtractBearerToken(r),
		id:    chi.URLParam(r, "thingID"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(apiutil.ErrValidation, errors.Wrap(errors.ErrMalformedEntity, err))
	}

	return req, nil
}

func decodeUpdateClientTags(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), api.ContentType) {
		return nil, errors.Wrap(apiutil.ErrValidation, apiutil.ErrUnsupportedContentType)
	}

	req := updateClientTagsReq{
		token: apiutil.ExtractBearerToken(r),
		id:    chi.URLParam(r, "thingID"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(apiutil.ErrValidation, errors.Wrap(errors.ErrMalformedEntity, err))
	}

	return req, nil
}

func decodeUpdateClientCredentials(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), api.ContentType) {
		return nil, errors.Wrap(apiutil.ErrValidation, apiutil.ErrUnsupportedContentType)
	}

	req := updateClientCredentialsReq{
		token: apiutil.ExtractBearerToken(r),
		id:    chi.URLParam(r, "thingID"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(apiutil.ErrValidation, errors.Wrap(errors.ErrMalformedEntity, err))
	}

	return req, nil
}

func decodeCreateClientReq(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), api.ContentType) {
		return nil, errors.Wrap(apiutil.ErrValidation, apiutil.ErrUnsupportedContentType)
	}

	var c mgclients.Client
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		return nil, errors.Wrap(apiutil.ErrValidation, errors.Wrap(errors.ErrMalformedEntity, err))
	}
	req := createClientReq{
		client: c,
		token:  apiutil.ExtractBearerToken(r),
	}

	return req, nil
}

func decodeCreateClientsReq(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), api.ContentType) {
		return nil, errors.Wrap(apiutil.ErrValidation, apiutil.ErrUnsupportedContentType)
	}

	c := createClientsReq{token: apiutil.ExtractBearerToken(r)}
	if err := json.NewDecoder(r.Body).Decode(&c.Clients); err != nil {
		return nil, errors.Wrap(apiutil.ErrValidation, errors.Wrap(errors.ErrMalformedEntity, err))
	}

	return c, nil
}

func decodeChangeClientStatus(_ context.Context, r *http.Request) (interface{}, error) {
	req := changeClientStatusReq{
		token: apiutil.ExtractBearerToken(r),
		id:    chi.URLParam(r, "thingID"),
	}

	return req, nil
}

func decodeListMembersRequest(_ context.Context, r *http.Request) (interface{}, error) {
	s, err := apiutil.ReadStringQuery(r, api.StatusKey, api.DefClientStatus)
	if err != nil {
		return nil, errors.Wrap(apiutil.ErrValidation, err)
	}
	o, err := apiutil.ReadNumQuery[uint64](r, api.OffsetKey, api.DefOffset)
	if err != nil {
		return nil, errors.Wrap(apiutil.ErrValidation, err)
	}
	l, err := apiutil.ReadNumQuery[uint64](r, api.LimitKey, api.DefLimit)
	if err != nil {
		return nil, errors.Wrap(apiutil.ErrValidation, err)
	}
	m, err := apiutil.ReadMetadataQuery(r, api.MetadataKey, nil)
	if err != nil {
		return nil, errors.Wrap(apiutil.ErrValidation, err)
	}
	st, err := mgclients.ToStatus(s)
	if err != nil {
		return nil, errors.Wrap(apiutil.ErrValidation, err)
	}
	p, err := apiutil.ReadStringQuery(r, api.PermissionKey, api.DefPermission)
	if err != nil {
		return nil, errors.Wrap(apiutil.ErrValidation, err)
	}

	lp, err := apiutil.ReadBoolQuery(r, api.ListPerms, api.DefListPerms)
	if err != nil {
		return nil, errors.Wrap(apiutil.ErrValidation, err)
	}
	req := listMembersReq{
		token: apiutil.ExtractBearerToken(r),
		Page: mgclients.Page{
			Status:     st,
			Offset:     o,
			Limit:      l,
			Permission: p,
			Metadata:   m,
			ListPerms:  lp,
		},
		groupID: chi.URLParam(r, "groupID"),
	}
	return req, nil
}

func decodeThingShareRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), api.ContentType) {
		return nil, errors.Wrap(apiutil.ErrValidation, apiutil.ErrUnsupportedContentType)
	}

	req := thingShareRequest{
		token:   apiutil.ExtractBearerToken(r),
		thingID: chi.URLParam(r, "thingID"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(apiutil.ErrValidation, errors.Wrap(errors.ErrMalformedEntity, err))
	}

	return req, nil
}

func decodeThingUnshareRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), api.ContentType) {
		return nil, errors.Wrap(apiutil.ErrValidation, apiutil.ErrUnsupportedContentType)
	}

	req := thingUnshareRequest{
		token:   apiutil.ExtractBearerToken(r),
		thingID: chi.URLParam(r, "thingID"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(apiutil.ErrValidation, errors.Wrap(errors.ErrMalformedEntity, err))
	}

	return req, nil
}
