// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package api_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/absmach/magistrala/internal/apiutil"
	"github.com/absmach/magistrala/internal/testsutil"
	"github.com/absmach/magistrala/invitations"
	"github.com/absmach/magistrala/invitations/api"
	"github.com/absmach/magistrala/invitations/mocks"
	mglog "github.com/absmach/magistrala/logger"
	"github.com/absmach/magistrala/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	validToken      = "valid"
	validContenType = "application/json"
	validID         = testsutil.GenerateUUID(&testing.T{})
)

type testRequest struct {
	client      *http.Client
	method      string
	url         string
	token       string
	contentType string
	body        io.Reader
}

func (tr testRequest) make() (*http.Response, error) {
	req, err := http.NewRequest(tr.method, tr.url, tr.body)
	if err != nil {
		return nil, err
	}

	if tr.token != "" {
		req.Header.Set("Authorization", apiutil.BearerPrefix+tr.token)
	}

	if tr.contentType != "" {
		req.Header.Set("Content-Type", tr.contentType)
	}

	return tr.client.Do(req)
}

func newIvitationsServer() (*httptest.Server, *mocks.Service) {
	svc := new(mocks.Service)

	logger := mglog.NewMock()
	mux := api.MakeHandler(svc, logger, "test")
	return httptest.NewServer(mux), svc
}

func TestSendInvitation(t *testing.T) {
	is, svc := newIvitationsServer()

	cases := []struct {
		desc        string
		token       string
		data        string
		contentType string
		status      int
		svcErr      error
	}{
		{
			desc:        "valid request",
			token:       validToken,
			data:        fmt.Sprintf(`{"user_id": "%s", "domain": "%s", "relation": "%s"}`, validID, validID, "domain"),
			status:      http.StatusCreated,
			contentType: validContenType,
			svcErr:      nil,
		},
		{
			desc:        "invalid token",
			token:       "",
			data:        fmt.Sprintf(`{"user_id": "%s", "domain": "%s", "relation": "%s"}`, validID, validID, "domain"),
			status:      http.StatusUnauthorized,
			contentType: validContenType,
			svcErr:      nil,
		},
		{
			desc:        "invalid content type",
			token:       validToken,
			data:        fmt.Sprintf(`{"user_id": "%s", "domain": "%s", "relation": "%s"}`, validID, validID, "domain"),
			status:      http.StatusUnsupportedMediaType,
			contentType: "text/plain",
			svcErr:      nil,
		},
		{
			desc:        "invalid data",
			token:       validToken,
			data:        `data`,
			status:      http.StatusBadRequest,
			contentType: validContenType,
			svcErr:      nil,
		},
		{
			desc:        "with service error",
			token:       validToken,
			data:        fmt.Sprintf(`{"user_id": "%s", "domain": "%s", "relation": "%s"}`, validID, validID, "domain"),
			status:      http.StatusForbidden,
			contentType: validContenType,
			svcErr:      errors.ErrAuthorization,
		},
	}

	for _, tc := range cases {
		repoCall := svc.On("SendInvitation", mock.Anything, tc.token, mock.Anything).Return(tc.svcErr)
		req := testRequest{
			client:      is.Client(),
			method:      http.MethodPost,
			url:         is.URL + "/invitations",
			token:       tc.token,
			contentType: tc.contentType,
			body:        strings.NewReader(tc.data),
		}

		res, err := req.make()
		assert.Nil(t, err, tc.desc)
		assert.Equal(t, tc.status, res.StatusCode, tc.desc)
		repoCall.Unset()
	}
}

func TestListInvitation(t *testing.T) {
	is, svc := newIvitationsServer()

	cases := []struct {
		desc        string
		token       string
		query       string
		contentType string
		status      int
		svcErr      error
	}{
		{
			desc:        "valid request",
			token:       validToken,
			status:      http.StatusOK,
			contentType: validContenType,
			svcErr:      nil,
		},
		{
			desc:        "invalid token",
			token:       "",
			status:      http.StatusUnauthorized,
			contentType: validContenType,
			svcErr:      nil,
		},
		{
			desc:        "with offset",
			token:       validToken,
			query:       "offset=1",
			status:      http.StatusOK,
			contentType: validContenType,
			svcErr:      nil,
		},
		{
			desc:        "with invalid offset",
			token:       validToken,
			query:       "offset=invalid",
			status:      http.StatusInternalServerError,
			contentType: validContenType,
			svcErr:      nil,
		},
		{
			desc:        "with limit",
			token:       validToken,
			query:       "limit=1",
			status:      http.StatusOK,
			contentType: validContenType,
			svcErr:      nil,
		},
		{
			desc:        "with invalid limit",
			token:       validToken,
			query:       "limit=invalid",
			status:      http.StatusInternalServerError,
			contentType: validContenType,
			svcErr:      nil,
		},
		{
			desc:        "with user_id",
			token:       validToken,
			query:       fmt.Sprintf("user_id=%s", validID),
			status:      http.StatusOK,
			contentType: validContenType,
			svcErr:      nil,
		},
		{
			desc:        "with duplicate user_id",
			token:       validToken,
			query:       "user_id=1&user_id=2",
			status:      http.StatusInternalServerError,
			contentType: validContenType,
			svcErr:      nil,
		},
		{
			desc:        "with invited_by",
			token:       validToken,
			query:       fmt.Sprintf("invited_by=%s", validID),
			status:      http.StatusOK,
			contentType: validContenType,
			svcErr:      nil,
		},
		{
			desc:        "with duplicate invited_by",
			token:       validToken,
			query:       "invited_by=1&invited_by=2",
			status:      http.StatusInternalServerError,
			contentType: validContenType,
			svcErr:      nil,
		},
		{
			desc:        "with relation",
			token:       validToken,
			query:       fmt.Sprintf("relation=%s", "relation"),
			status:      http.StatusOK,
			contentType: validContenType,
			svcErr:      nil,
		},
		{
			desc:        "with duplicate relation",
			token:       validToken,
			query:       "relation=1&relation=2",
			status:      http.StatusInternalServerError,
			contentType: validContenType,
			svcErr:      nil,
		},
		{
			desc:        "with domain",
			token:       validToken,
			query:       fmt.Sprintf("domain=%s", validID),
			status:      http.StatusOK,
			contentType: validContenType,
			svcErr:      nil,
		},
		{
			desc:        "with duplicate domain",
			token:       validToken,
			query:       "domain=1&domain=2",
			status:      http.StatusInternalServerError,
			contentType: validContenType,
			svcErr:      nil,
		},
		{
			desc:        "with service error",
			token:       validToken,
			status:      http.StatusForbidden,
			contentType: validContenType,
			svcErr:      errors.ErrAuthorization,
		},
	}

	for _, tc := range cases {
		repoCall := svc.On("ListInvitations", mock.Anything, tc.token, mock.Anything).Return(invitations.InvitationPage{}, tc.svcErr)
		req := testRequest{
			client:      is.Client(),
			method:      http.MethodGet,
			url:         is.URL + "/invitations?" + tc.query,
			token:       tc.token,
			contentType: tc.contentType,
		}

		res, err := req.make()
		assert.Nil(t, err, tc.desc)
		assert.Equal(t, tc.status, res.StatusCode, tc.desc)
		repoCall.Unset()
	}
}

func TestViewInvitation(t *testing.T) {
	is, svc := newIvitationsServer()

	cases := []struct {
		desc        string
		token       string
		domain      string
		userID      string
		contentType string
		status      int
		svcErr      error
	}{
		{
			desc:        "valid request",
			token:       validToken,
			userID:      validID,
			domain:      validID,
			status:      http.StatusOK,
			contentType: validContenType,
			svcErr:      nil,
		},
		{
			desc:        "invalid token",
			token:       "",
			userID:      validID,
			domain:      validID,
			status:      http.StatusUnauthorized,
			contentType: validContenType,
			svcErr:      nil,
		},
		{
			desc:        "with service error",
			token:       validToken,
			userID:      validID,
			domain:      validID,
			status:      http.StatusForbidden,
			contentType: validContenType,
			svcErr:      errors.ErrAuthorization,
		},
		{
			desc:        "with empty user_id",
			token:       validToken,
			userID:      "",
			domain:      validID,
			status:      http.StatusBadRequest,
			contentType: validContenType,
			svcErr:      nil,
		},
		{
			desc:        "with empty domain",
			token:       validToken,
			userID:      validID,
			domain:      "",
			status:      http.StatusNotFound,
			contentType: validContenType,
			svcErr:      nil,
		},
		{
			desc:        "with empty user_id and domain",
			token:       validToken,
			userID:      "",
			domain:      "",
			status:      http.StatusNotFound,
			contentType: validContenType,
			svcErr:      nil,
		},
	}

	for _, tc := range cases {
		repoCall := svc.On("ViewInvitation", mock.Anything, tc.token, tc.userID, tc.domain).Return(invitations.Invitation{}, tc.svcErr)
		req := testRequest{
			client:      is.Client(),
			method:      http.MethodGet,
			url:         is.URL + "/invitations/" + tc.userID + "/" + tc.domain,
			token:       tc.token,
			contentType: tc.contentType,
		}

		res, err := req.make()
		assert.Nil(t, err, tc.desc)
		assert.Equal(t, tc.status, res.StatusCode, tc.desc)
		repoCall.Unset()
	}
}

func TestDeleteInvitation(t *testing.T) {
	is, svc := newIvitationsServer()

	cases := []struct {
		desc        string
		token       string
		domain      string
		userID      string
		contentType string
		status      int
		svcErr      error
	}{
		{
			desc:        "valid request",
			token:       validToken,
			userID:      validID,
			domain:      validID,
			status:      http.StatusNoContent,
			contentType: validContenType,
			svcErr:      nil,
		},
		{
			desc:        "invalid token",
			token:       "",
			userID:      validID,
			domain:      validID,
			status:      http.StatusUnauthorized,
			contentType: validContenType,
			svcErr:      nil,
		},
		{
			desc:        "with service error",
			token:       validToken,
			userID:      validID,
			domain:      validID,
			status:      http.StatusForbidden,
			contentType: validContenType,
			svcErr:      errors.ErrAuthorization,
		},
		{
			desc:        "with empty user_id",
			token:       validToken,
			userID:      "",
			domain:      validID,
			status:      http.StatusBadRequest,
			contentType: validContenType,
			svcErr:      nil,
		},
		{
			desc:        "with empty domain",
			token:       validToken,
			userID:      validID,
			domain:      "",
			status:      http.StatusNotFound,
			contentType: validContenType,
			svcErr:      nil,
		},
		{
			desc:        "with empty user_id and domain",
			token:       validToken,
			userID:      "",
			domain:      "",
			status:      http.StatusNotFound,
			contentType: validContenType,
			svcErr:      nil,
		},
	}

	for _, tc := range cases {
		repoCall := svc.On("DeleteInvitation", mock.Anything, tc.token, tc.userID, tc.domain).Return(tc.svcErr)
		req := testRequest{
			client:      is.Client(),
			method:      http.MethodDelete,
			url:         is.URL + "/invitations/" + tc.userID + "/" + tc.domain,
			token:       tc.token,
			contentType: tc.contentType,
		}

		res, err := req.make()
		assert.Nil(t, err, tc.desc)
		assert.Equal(t, tc.status, res.StatusCode, tc.desc)
		repoCall.Unset()
	}
}

func TestAcceptInvitation(t *testing.T) {
	is, svc := newIvitationsServer()

	cases := []struct {
		desc        string
		token       string
		contentType string
		status      int
		svcErr      error
	}{
		{
			desc:        "valid request",
			token:       validToken,
			status:      http.StatusOK,
			contentType: validContenType,
			svcErr:      nil,
		},
		{
			desc:        "invalid token",
			token:       "",
			status:      http.StatusUnauthorized,
			contentType: validContenType,
			svcErr:      nil,
		},
		{
			desc:        "with service error",
			token:       validToken,
			status:      http.StatusForbidden,
			contentType: validContenType,
			svcErr:      errors.ErrAuthorization,
		},
	}

	for _, tc := range cases {
		repoCall := svc.On("AcceptInvitation", mock.Anything, tc.token).Return([]string{}, tc.svcErr)
		req := testRequest{
			client:      is.Client(),
			method:      http.MethodPost,
			url:         is.URL + "/invitations/accept",
			token:       tc.token,
			contentType: tc.contentType,
		}

		res, err := req.make()
		assert.Nil(t, err, tc.desc)
		assert.Equal(t, tc.status, res.StatusCode, tc.desc)
		repoCall.Unset()
	}
}
