// Copyright (c) Magistrala
// SPDX-License-Identifier: Apache-2.0

package policies

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/absmach/magistrala"
	"github.com/absmach/magistrala/auth"
	"github.com/absmach/magistrala/internal/apiutil"
	"github.com/absmach/magistrala/logger"
	"github.com/absmach/magistrala/pkg/errors"
	"github.com/go-chi/chi"
	kithttp "github.com/go-kit/kit/transport/http"
)

const contentType = "application/json"

// MakeHandler returns a HTTP handler for API endpoints.
func MakeHandler(svc auth.Service, mux *chi.Mux, logger logger.Logger) *chi.Mux {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(apiutil.LoggingErrorEncoder(logger, encodeError)),
	}
	mux.Route("/policies", func(r chi.Router) {
		r.Post("/", createPolicyEndpointHandler(svc, opts))

		r.Route("/delete", func(r chi.Router) {
			r.Post("/", deletePoliciesEndpointHandler(svc, opts))
		})
	})

	return mux
}

func createPolicyEndpointHandler(svc auth.Service, opts []kithttp.ServerOption) http.HandlerFunc {
	return kithttp.NewServer(
		(createPolicyEndpoint(svc)),
		decodePoliciesRequest,
		encodeResponse,
		opts...,
	).ServeHTTP
}

func deletePoliciesEndpointHandler(svc auth.Service, opts []kithttp.ServerOption) http.HandlerFunc {
	return kithttp.NewServer(
		(deletePoliciesEndpoint(svc)),
		decodePoliciesRequest,
		encodeResponse,
		opts...,
	).ServeHTTP
}

func decodePoliciesRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errors.ErrUnsupportedContentType
	}

	req := policiesReq{token: apiutil.ExtractBearerToken(r)}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(errors.ErrMalformedEntity, err)
	}

	return req, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", contentType)

	if ar, ok := response.(magistrala.Response); ok {
		for k, v := range ar.Headers() {
			w.Header().Set(k, v)
		}

		w.WriteHeader(ar.Code())

		if ar.Empty() {
			return nil
		}
	}

	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	switch {
	case errors.Contains(err, errors.ErrMalformedEntity),
		err == apiutil.ErrEmptyList,
		err == apiutil.ErrMissingPolicyObj,
		err == apiutil.ErrMissingPolicySub,
		err == apiutil.ErrMalformedPolicy:
		w.WriteHeader(http.StatusBadRequest)
	case errors.Contains(err, errors.ErrAuthentication),
		err == apiutil.ErrBearerToken:
		w.WriteHeader(http.StatusUnauthorized)
	case errors.Contains(err, errors.ErrNotFound):
		w.WriteHeader(http.StatusNotFound)
	case errors.Contains(err, errors.ErrConflict):
		w.WriteHeader(http.StatusConflict)
	case errors.Contains(err, errors.ErrAuthorization):
		w.WriteHeader(http.StatusForbidden)
	case errors.Contains(err, errors.ErrUnsupportedContentType):
		w.WriteHeader(http.StatusUnsupportedMediaType)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	if errorVal, ok := err.(errors.Error); ok {
		w.Header().Set("Content-Type", contentType)
		if err := json.NewEncoder(w).Encode(apiutil.ErrorRes{Err: errorVal.Msg()}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
