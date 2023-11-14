// Copyright (c) Magistrala
// SPDX-License-Identifier: Apache-2.0

package policies_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/absmach/magistrala/auth"
	httpapi "github.com/absmach/magistrala/auth/api/http"
	"github.com/absmach/magistrala/auth/jwt"
	"github.com/absmach/magistrala/auth/mocks"
	"github.com/absmach/magistrala/internal/apiutil"
	"github.com/absmach/magistrala/logger"
	"github.com/absmach/magistrala/pkg/uuid"
	"github.com/stretchr/testify/assert"
)

const (
	secret          = "secret"
	contentType     = "application/json"
	id              = uuid.Prefix + "-000000000001"
	email           = "user@example.com"
	unauthzID       = uuid.Prefix + "-000000000002"
	unauthzEmail    = "unauthz@example.com"
	loginDuration   = 30 * time.Minute
	refreshDuration = 24 * time.Hour
)

type testRequest struct {
	client      *http.Client
	method      string
	url         string
	contentType string
	token       string
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

	req.Header.Set("Referer", "http://localhost")
	return tr.client.Do(req)
}

func newService() auth.Service {
	krepo := new(mocks.Keys)
	prepo := new(mocks.PolicyAgent)
	drepo := new(mocks.DomainsRepo)
	idProvider := uuid.NewMock()

	t := jwt.New([]byte(secret))

	return auth.New(krepo, drepo, idProvider, t, prepo, loginDuration, refreshDuration)
}

func newServer(svc auth.Service) *httptest.Server {
	logger := logger.NewMock()
	mux := httpapi.MakeHandler(svc, logger, "")
	return httptest.NewServer(mux)
}

func toJSON(data interface{}) string {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(jsonData)
}

type addPolicyRequest struct {
	SubjectIDs []string `json:"subjects"`
	Policies   []string `json:"policies"`
	Object     string   `json:"object"`
}

func TestAddPolicies(t *testing.T) {
	svc := newService()
	token, err := svc.Issue(context.Background(), "", auth.Key{Type: auth.AccessKey, IssuedAt: time.Now(), Subject: id})
	assert.Nil(t, err, fmt.Sprintf("Issuing user key expected to succeed: %s", err))

	userLoginToken, err := svc.Issue(context.Background(), "", auth.Key{Type: auth.AccessKey, IssuedAt: time.Now(), Subject: unauthzID})
	assert.Nil(t, err, fmt.Sprintf("Issuing unauthorized user's key expected to succeed: %s", err))

	ts := newServer(svc)
	defer ts.Close()
	client := ts.Client()

	policies := []auth.PolicyReq{
		{
			Subject:     "user1",
			SubjectType: auth.UserType,
			Relation:    auth.ViewerRelation,
			Object:      "thing1",
			ObjectType:  auth.ThingType,
		},
	}
	cases := []struct {
		desc   string
		token  string
		ct     string
		status int
		req    string
	}{
		{
			desc:   "Add policies ",
			token:  token.AccessToken,
			ct:     contentType,
			status: http.StatusCreated,
			req:    toJSON(policies),
		},
		{
			desc:   "Add policies ",
			token:  userLoginToken.AccessToken,
			ct:     contentType,
			status: http.StatusCreated,
			req:    toJSON(policies),
		},
	}

	for _, tc := range cases {
		req := testRequest{
			client:      client,
			method:      http.MethodPost,
			url:         fmt.Sprintf("%s/policies", ts.URL),
			contentType: tc.ct,
			token:       tc.token,
			body:        strings.NewReader(tc.req),
		}

		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
	}
}

func TestDeletePolicies(t *testing.T) {
	svc := newService()
	token, err := svc.Issue(context.Background(), "", auth.Key{Type: auth.AccessKey, IssuedAt: time.Now(), Subject: id})
	assert.Nil(t, err, fmt.Sprintf("Issuing user key expected to succeed: %s", err))

	userLoginToken, err := svc.Issue(context.Background(), "", auth.Key{Type: auth.AccessKey, IssuedAt: time.Now(), Subject: unauthzID})
	assert.Nil(t, err, fmt.Sprintf("Issuing unauthorized user's key expected to succeed: %s", err))

	ts := newServer(svc)
	defer ts.Close()
	client := ts.Client()

	policies := []auth.PolicyReq{
		{
			Subject:     "user1",
			SubjectType: auth.UserType,
			Relation:    auth.ViewerRelation,
			Object:      "thing1",
			ObjectType:  auth.ThingType,
		},
	}

	err = svc.AddPolicies(context.Background(), policies)
	assert.Nil(t, err, fmt.Sprintf("Adding policies expected to succeed: %s", err))

	cases := []struct {
		desc   string
		token  string
		ct     string
		req    string
		status int
	}{
		{
			desc:   "Delete policies with valid access",
			token:  token.AccessToken,
			ct:     contentType,
			status: http.StatusForbidden,
			req:    toJSON(policies),
		},
		{
			desc:   "Delete policies with unauthorized access",
			token:  userLoginToken.AccessToken,
			ct:     contentType,
			status: http.StatusForbidden,
			req:    toJSON(policies),
		},
	}

	for _, tc := range cases {
		req := testRequest{
			client:      client,
			method:      http.MethodPut,
			url:         fmt.Sprintf("%s/policies", ts.URL),
			contentType: tc.ct,
			token:       tc.token,
			body:        strings.NewReader(tc.req),
		}

		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
	}
}
