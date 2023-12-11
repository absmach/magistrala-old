// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/absmach/magistrala/auth"
	"github.com/absmach/magistrala/internal/api"
	"github.com/absmach/magistrala/internal/apiutil"
	mglog "github.com/absmach/magistrala/logger"
	mgclients "github.com/absmach/magistrala/pkg/clients"
	"github.com/absmach/magistrala/pkg/errors"
	"github.com/absmach/magistrala/users"
	"github.com/go-chi/chi/v5"
	kithttp "github.com/go-kit/kit/transport/http"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// MakeHandler returns a HTTP handler for API endpoints.
func clientsHandler(svc users.Service, r *chi.Mux, logger mglog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(apiutil.LoggingErrorEncoder(logger, api.EncodeError)),
	}
	r.Route("/users", func(r chi.Router) {
		r.Post("/", otelhttp.NewHandler(kithttp.NewServer(
			registrationEndpoint(svc),
			decodeCreateClientReq,
			api.EncodeResponse,
			opts...,
		), "register_client").ServeHTTP)

		r.Get("/profile", otelhttp.NewHandler(kithttp.NewServer(
			viewProfileEndpoint(svc),
			decodeViewProfile,
			api.EncodeResponse,
			opts...,
		), "view_profile").ServeHTTP)

		r.Get("/{id}", otelhttp.NewHandler(kithttp.NewServer(
			viewClientEndpoint(svc),
			decodeViewClient,
			api.EncodeResponse,
			opts...,
		), "view_client").ServeHTTP)

		r.Get("/", otelhttp.NewHandler(kithttp.NewServer(
			listClientsEndpoint(svc),
			decodeListClients,
			api.EncodeResponse,
			opts...,
		), "list_clients").ServeHTTP)

		r.Get("/search", otelhttp.NewHandler(kithttp.NewServer(
			searchClientsEndpoint(svc),
			decodeListClients,
			api.EncodeResponse,
			opts...,
		), "search_clients").ServeHTTP)

		r.Patch("/secret", otelhttp.NewHandler(kithttp.NewServer(
			updateClientSecretEndpoint(svc),
			decodeUpdateClientSecret,
			api.EncodeResponse,
			opts...,
		), "update_client_secret").ServeHTTP)

		r.Patch("/{id}", otelhttp.NewHandler(kithttp.NewServer(
			updateClientEndpoint(svc),
			decodeUpdateClient,
			api.EncodeResponse,
			opts...,
		), "update_client").ServeHTTP)

		r.Patch("/{id}/tags", otelhttp.NewHandler(kithttp.NewServer(
			updateClientTagsEndpoint(svc),
			decodeUpdateClientTags,
			api.EncodeResponse,
			opts...,
		), "update_client_tags").ServeHTTP)

		r.Patch("/{id}/identity", otelhttp.NewHandler(kithttp.NewServer(
			updateClientIdentityEndpoint(svc),
			decodeUpdateClientIdentity,
			api.EncodeResponse,
			opts...,
		), "update_client_identity").ServeHTTP)

		r.Post("/password/reset-request", otelhttp.NewHandler(kithttp.NewServer(
			passwordResetRequestEndpoint(svc),
			decodePasswordResetRequest,
			api.EncodeResponse,
			opts...,
		), "password_reset_req").ServeHTTP)

		r.Put("/password/reset", otelhttp.NewHandler(kithttp.NewServer(
			passwordResetEndpoint(svc),
			decodePasswordReset,
			api.EncodeResponse,
			opts...,
		), "password_reset").ServeHTTP)

		r.Patch("/{id}/role", otelhttp.NewHandler(kithttp.NewServer(
			updateClientRoleEndpoint(svc),
			decodeUpdateClientRole,
			api.EncodeResponse,
			opts...,
		), "update_client_role").ServeHTTP)

		r.Post("/tokens/issue", otelhttp.NewHandler(kithttp.NewServer(
			issueTokenEndpoint(svc),
			decodeCredentials,
			api.EncodeResponse,
			opts...,
		), "issue_token").ServeHTTP)

		r.Post("/tokens/refresh", otelhttp.NewHandler(kithttp.NewServer(
			refreshTokenEndpoint(svc),
			decodeRefreshToken,
			api.EncodeResponse,
			opts...,
		), "refresh_token").ServeHTTP)

		r.Post("/{id}/enable", otelhttp.NewHandler(kithttp.NewServer(
			enableClientEndpoint(svc),
			decodeChangeClientStatus,
			api.EncodeResponse,
			opts...,
		), "enable_client").ServeHTTP)

		r.Post("/{id}/disable", otelhttp.NewHandler(kithttp.NewServer(
			disableClientEndpoint(svc),
			decodeChangeClientStatus,
			api.EncodeResponse,
			opts...,
		), "disable_client").ServeHTTP)
	})

	// Ideal location: users service, groups endpoint.
	// Reason for placing here :
	// SpiceDB provides list of user ids in given user_group_id
	// and users service can access spiceDB and get the user list with user_group_id.
	// Request to get list of users present in the user_group_id {groupID}
	r.Get("/groups/{groupID}/users", otelhttp.NewHandler(kithttp.NewServer(
		listMembersByGroupEndpoint(svc),
		decodeListMembersByGroup,
		api.EncodeResponse,
		opts...,
	), "list_users_by_user_group_id").ServeHTTP)

	// Ideal location: things service, channels endpoint.
	// Reason for placing here :
	// SpiceDB provides list of user ids in given channel_id
	// and users service can access spiceDB and get the user list with channel_id.
	// Request to get list of users present in the user_group_id {channelID}
	r.Get("/channels/{channelID}/users", otelhttp.NewHandler(kithttp.NewServer(
		listMembersByChannelEndpoint(svc),
		decodeListMembersByChannel,
		api.EncodeResponse,
		opts...,
	), "list_users_by_channel_id").ServeHTTP)

	r.Get("/things/{thingID}/users", otelhttp.NewHandler(kithttp.NewServer(
		listMembersByThingEndpoint(svc),
		decodeListMembersByThing,
		api.EncodeResponse,
		opts...,
	), "list_users_by_thing_id").ServeHTTP)

	r.Get("/domains/{domainID}/users", otelhttp.NewHandler(kithttp.NewServer(
		listMembersByDomainEndpoint(svc),
		decodeListMembersByDomain,
		api.EncodeResponse,
		opts...,
	), "list_users_by_domain_id").ServeHTTP)
	return r
}

func decodeViewClient(_ context.Context, r *http.Request) (interface{}, error) {
	req := viewClientReq{
		token: apiutil.ExtractBearerToken(r),
		id:    chi.URLParam(r, "id"),
	}

	return req, nil
}

func decodeViewProfile(_ context.Context, r *http.Request) (interface{}, error) {
	req := viewProfileReq{token: apiutil.ExtractBearerToken(r)}

	return req, nil
}

func decodeListClients(_ context.Context, r *http.Request) (interface{}, error) {
	var sharedID, ownerID string
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
	i, err := apiutil.ReadStringQuery(r, api.IdentityKey, "")
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
	visibility, err := apiutil.ReadStringQuery(r, api.VisibilityKey, "")
	if err != nil {
		return nil, errors.Wrap(apiutil.ErrValidation, err)
	}
	order, err := apiutil.ReadStringQuery(r, api.OrderKey, api.DefOrder)
	if err != nil {
		return nil, errors.Wrap(apiutil.ErrValidation, err)
	}
	dir, err := apiutil.ReadStringQuery(r, api.DirKey, api.DefDir)
	if err != nil {
		return nil, errors.Wrap(apiutil.ErrValidation, err)
	}
	switch visibility {
	case api.MyVisibility:
		ownerID = api.MyVisibility
	case api.SharedVisibility:
		sharedID = api.MyVisibility
	case api.AllVisibility:
		sharedID = api.MyVisibility
		ownerID = api.MyVisibility
	}
	if oid != "" {
		ownerID = oid
	}
	st, err := mgclients.ToStatus(s)
	if err != nil {
		return nil, errors.Wrap(apiutil.ErrValidation, err)
	}
	req := listClientsReq{
		token:    apiutil.ExtractBearerToken(r),
		status:   st,
		offset:   o,
		limit:    l,
		metadata: m,
		name:     n,
		identity: i,
		tag:      t,
		sharedBy: sharedID,
		owner:    ownerID,
		order:    order,
		dir:      dir,
	}

	return req, nil
}

func decodeUpdateClient(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), api.ContentType) {
		return nil, errors.Wrap(apiutil.ErrValidation, apiutil.ErrUnsupportedContentType)
	}

	req := updateClientReq{
		token: apiutil.ExtractBearerToken(r),
		id:    chi.URLParam(r, "id"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(apiutil.ErrValidation, errors.Wrap(err, errors.ErrMalformedEntity))
	}

	return req, nil
}

func decodeUpdateClientTags(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), api.ContentType) {
		return nil, errors.Wrap(apiutil.ErrValidation, apiutil.ErrUnsupportedContentType)
	}

	req := updateClientTagsReq{
		token: apiutil.ExtractBearerToken(r),
		id:    chi.URLParam(r, "id"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(apiutil.ErrValidation, errors.Wrap(err, errors.ErrMalformedEntity))
	}

	return req, nil
}

func decodeUpdateClientIdentity(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), api.ContentType) {
		return nil, errors.Wrap(apiutil.ErrValidation, apiutil.ErrUnsupportedContentType)
	}

	req := updateClientIdentityReq{
		token: apiutil.ExtractBearerToken(r),
		id:    chi.URLParam(r, "id"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(apiutil.ErrValidation, errors.Wrap(err, errors.ErrMalformedEntity))
	}

	return req, nil
}

func decodeUpdateClientSecret(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), api.ContentType) {
		return nil, errors.Wrap(apiutil.ErrValidation, apiutil.ErrUnsupportedContentType)
	}

	req := updateClientSecretReq{
		token: apiutil.ExtractBearerToken(r),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(apiutil.ErrValidation, errors.Wrap(err, errors.ErrMalformedEntity))
	}

	return req, nil
}

func decodePasswordResetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), api.ContentType) {
		return nil, apiutil.ErrUnsupportedContentType
	}

	var req passwResetReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(apiutil.ErrValidation, errors.Wrap(err, errors.ErrMalformedEntity))
	}

	req.Host = r.Header.Get("Referer")
	return req, nil
}

func decodePasswordReset(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), api.ContentType) {
		return nil, errors.Wrap(apiutil.ErrValidation, apiutil.ErrUnsupportedContentType)
	}

	var req resetTokenReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(apiutil.ErrValidation, errors.Wrap(err, errors.ErrMalformedEntity))
	}

	return req, nil
}

func decodeUpdateClientRole(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), api.ContentType) {
		return nil, errors.Wrap(apiutil.ErrValidation, apiutil.ErrUnsupportedContentType)
	}

	req := updateClientRoleReq{
		token: apiutil.ExtractBearerToken(r),
		id:    chi.URLParam(r, "id"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(apiutil.ErrValidation, errors.Wrap(err, errors.ErrMalformedEntity))
	}
	var err error
	req.role, err = mgclients.ToRole(req.Role)
	return req, err
}

func decodeCredentials(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), api.ContentType) {
		return nil, errors.Wrap(apiutil.ErrValidation, apiutil.ErrUnsupportedContentType)
	}

	req := loginClientReq{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(apiutil.ErrValidation, errors.Wrap(err, errors.ErrMalformedEntity))
	}

	return req, nil
}

func decodeRefreshToken(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), api.ContentType) {
		return nil, errors.Wrap(apiutil.ErrValidation, apiutil.ErrUnsupportedContentType)
	}
	req := tokenReq{RefreshToken: apiutil.ExtractBearerToken(r)}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(apiutil.ErrValidation, errors.Wrap(err, errors.ErrMalformedEntity))
	}
	return req, nil
}

func decodeCreateClientReq(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), api.ContentType) {
		return nil, errors.Wrap(apiutil.ErrValidation, apiutil.ErrUnsupportedContentType)
	}

	var c mgclients.Client
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		return nil, errors.Wrap(apiutil.ErrValidation, errors.Wrap(err, errors.ErrMalformedEntity))
	}
	req := createClientReq{
		client: c,
		token:  apiutil.ExtractBearerToken(r),
	}

	return req, nil
}

func decodeChangeClientStatus(_ context.Context, r *http.Request) (interface{}, error) {
	req := changeClientStatusReq{
		token: apiutil.ExtractBearerToken(r),
		id:    chi.URLParam(r, "id"),
	}

	return req, nil
}

func decodeListMembersByGroup(_ context.Context, r *http.Request) (interface{}, error) {
	page, err := queryPageParams(r, api.DefPermission)
	if err != nil {
		return nil, err
	}
	req := listMembersByObjectReq{
		token:    apiutil.ExtractBearerToken(r),
		Page:     page,
		objectID: chi.URLParam(r, "groupID"),
	}

	return req, nil
}

func decodeListMembersByChannel(_ context.Context, r *http.Request) (interface{}, error) {
	page, err := queryPageParams(r, api.DefPermission)
	if err != nil {
		return nil, err
	}
	req := listMembersByObjectReq{
		token:    apiutil.ExtractBearerToken(r),
		Page:     page,
		objectID: chi.URLParam(r, "channelID"),
	}

	return req, nil
}

func decodeListMembersByThing(_ context.Context, r *http.Request) (interface{}, error) {
	page, err := queryPageParams(r, api.DefPermission)
	if err != nil {
		return nil, err
	}
	req := listMembersByObjectReq{
		token:    apiutil.ExtractBearerToken(r),
		Page:     page,
		objectID: chi.URLParam(r, "thingID"),
	}

	return req, nil
}

func decodeListMembersByDomain(_ context.Context, r *http.Request) (interface{}, error) {
	page, err := queryPageParams(r, auth.MembershipPermission)
	if err != nil {
		return nil, err
	}
	// For domains default permission in membership, In "queryPageParams" default is view,
	// so overwriting the permission given by queryPageParams function with default membership permission.
	p, err := apiutil.ReadStringQuery(r, api.PermissionKey, auth.MembershipPermission)
	if err != nil {
		return mgclients.Page{}, errors.Wrap(apiutil.ErrValidation, err)
	}
	page.Permission = p
	req := listMembersByObjectReq{
		token:    apiutil.ExtractBearerToken(r),
		Page:     page,
		objectID: chi.URLParam(r, "domainID"),
	}

	return req, nil
}

func queryPageParams(r *http.Request, defPermission string) (mgclients.Page, error) {
	s, err := apiutil.ReadStringQuery(r, api.StatusKey, api.DefClientStatus)
	if err != nil {
		return mgclients.Page{}, errors.Wrap(apiutil.ErrValidation, err)
	}
	o, err := apiutil.ReadNumQuery[uint64](r, api.OffsetKey, api.DefOffset)
	if err != nil {
		return mgclients.Page{}, errors.Wrap(apiutil.ErrValidation, err)
	}
	l, err := apiutil.ReadNumQuery[uint64](r, api.LimitKey, api.DefLimit)
	if err != nil {
		return mgclients.Page{}, errors.Wrap(apiutil.ErrValidation, err)
	}
	m, err := apiutil.ReadMetadataQuery(r, api.MetadataKey, nil)
	if err != nil {
		return mgclients.Page{}, errors.Wrap(apiutil.ErrValidation, err)
	}
	n, err := apiutil.ReadStringQuery(r, api.NameKey, "")
	if err != nil {
		return mgclients.Page{}, errors.Wrap(apiutil.ErrValidation, err)
	}
	i, err := apiutil.ReadStringQuery(r, api.IdentityKey, "")
	if err != nil {
		return mgclients.Page{}, errors.Wrap(apiutil.ErrValidation, err)
	}
	t, err := apiutil.ReadStringQuery(r, api.TagKey, "")
	if err != nil {
		return mgclients.Page{}, errors.Wrap(apiutil.ErrValidation, err)
	}
	oid, err := apiutil.ReadStringQuery(r, api.OwnerKey, "")
	if err != nil {
		return mgclients.Page{}, errors.Wrap(apiutil.ErrValidation, err)
	}
	st, err := mgclients.ToStatus(s)
	if err != nil {
		return mgclients.Page{}, errors.Wrap(apiutil.ErrValidation, err)
	}
	p, err := apiutil.ReadStringQuery(r, api.PermissionKey, defPermission)
	if err != nil {
		return mgclients.Page{}, errors.Wrap(apiutil.ErrValidation, err)
	}
	return mgclients.Page{
		Status:     st,
		Offset:     o,
		Limit:      l,
		Metadata:   m,
		Identity:   i,
		Name:       n,
		Tag:        t,
		Owner:      oid,
		Permission: p,
	}, nil
}
