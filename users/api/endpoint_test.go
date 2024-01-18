// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package api_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/absmach/magistrala"
	authmocks "github.com/absmach/magistrala/auth/mocks"
	"github.com/absmach/magistrala/internal/api"
	"github.com/absmach/magistrala/internal/apiutil"
	"github.com/absmach/magistrala/internal/groups"
	"github.com/absmach/magistrala/internal/testsutil"
	mglog "github.com/absmach/magistrala/logger"
	mgclients "github.com/absmach/magistrala/pkg/clients"
	"github.com/absmach/magistrala/pkg/errors"
	svcerr "github.com/absmach/magistrala/pkg/errors/service"
	gmocks "github.com/absmach/magistrala/pkg/groups/mocks"
	"github.com/absmach/magistrala/pkg/uuid"
	httpapi "github.com/absmach/magistrala/users/api"
	"github.com/absmach/magistrala/users/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	idProvider     = uuid.New()
	secret         = "strongsecret"
	validCMetadata = mgclients.Metadata{"role": "client"}
	client         = mgclients.Client{
		ID:          testsutil.GenerateUUID(&testing.T{}),
		Name:        "clientname",
		Tags:        []string{"tag1", "tag2"},
		Credentials: mgclients.Credentials{Identity: "clientidentity@example.com", Secret: secret},
		Metadata:    validCMetadata,
		Status:      mgclients.EnabledStatus,
	}
	validToken        = "valid"
	inValidToken      = "invalid"
	inValid           = "invalid"
	validID           = "d4ebb847-5d0e-4e46-bdd9-b6aceaaa3a22"
	ErrPasswordFormat = errors.New("password does not meet the requirements")
)

const contentType = "application/json"

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

func newUsersServer() (*httptest.Server, *mocks.Service) {
	gRepo := new(gmocks.Repository)
	auth := new(authmocks.Service)

	svc := new(mocks.Service)
	gsvc := groups.NewService(gRepo, idProvider, auth)

	logger := mglog.NewMock()
	mux := chi.NewRouter()
	httpapi.MakeHandler(svc, gsvc, mux, logger, "")

	return httptest.NewServer(mux), svc
}

func toJSON(data interface{}) string {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(jsonData)
}

func TestRegisterClient(t *testing.T) {
	us, svc := newUsersServer()
	defer us.Close()

	cases := []struct {
		desc        string
		client      mgclients.Client
		token       string
		contentType string
		status      int
		err         error
	}{
		{
			desc:        "register  a new user with a valid token",
			client:      client,
			token:       validToken,
			contentType: contentType,
			status:      http.StatusCreated,
			err:         nil,
		},
		{
			desc:        "register an existing user",
			client:      client,
			token:       validToken,
			contentType: contentType,
			status:      http.StatusConflict,
			err:         errors.ErrConflict,
		},
		{
			desc:        "register a new user with an empty token",
			client:      client,
			token:       "",
			contentType: contentType,
			status:      http.StatusUnauthorized,
			err:         apiutil.ErrBearerToken,
		},
		{
			desc: "register a user with an  invalid ID",
			client: mgclients.Client{
				ID: inValid,
				Credentials: mgclients.Credentials{
					Identity: "user@example.com",
					Secret:   "12345678",
				},
			},
			token:       validToken,
			contentType: contentType,
			status:      http.StatusBadRequest,
			err:         apiutil.ErrValidation,
		},
		{
			desc: "register a user that can't be marshalled",
			client: mgclients.Client{
				Credentials: mgclients.Credentials{
					Identity: "user@example.com",
					Secret:   "12345678",
				},
				Metadata: map[string]interface{}{
					"test": make(chan int),
				},
			},
			token:       validToken,
			contentType: contentType,
			status:      http.StatusBadRequest,
			err:         apiutil.ErrValidation,
		},
		{
			desc: "register user with invalid status",
			client: mgclients.Client{
				Credentials: mgclients.Credentials{
					Identity: "newclientwithinvalidstatus@example.com",
					Secret:   secret,
				},
				Status: mgclients.AllStatus,
			},
			token:       validToken,
			contentType: contentType,
			status:      http.StatusBadRequest,
			err:         svcerr.ErrInvalidStatus,
		},
		{
			desc: "register a user with name too long",
			client: mgclients.Client{
				Name: strings.Repeat("a", 1025),
				Credentials: mgclients.Credentials{
					Identity: "newclientwithinvalidname@example.com",
					Secret:   secret,
				},
			},
			token:       validToken,
			contentType: contentType,
			status:      http.StatusBadRequest,
			err:         apiutil.ErrValidation,
		},
		{
			desc:        "register user with invalid content type",
			client:      client,
			token:       validToken,
			contentType: "application/xml",
			status:      http.StatusUnsupportedMediaType,
			err:         apiutil.ErrValidation,
		},
		{
			desc:        "register user with empty request body",
			client:      mgclients.Client{},
			token:       validToken,
			contentType: contentType,
			status:      http.StatusBadRequest,
			err:         apiutil.ErrValidation,
		},
	}

	for _, tc := range cases {
		data := toJSON(tc.client)
		req := testRequest{
			client:      us.Client(),
			method:      http.MethodPost,
			url:         fmt.Sprintf("%s/users/", us.URL),
			contentType: tc.contentType,
			token:       tc.token,
			body:        strings.NewReader(data),
		}

		repoCall := svc.On("RegisterClient", mock.Anything, tc.token, tc.client).Return(tc.client, tc.err)
		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		var errRes respBody
		err = json.NewDecoder(res.Body).Decode(&errRes)
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error while decoding response body: %s", tc.desc, err))
		if errRes.Err != "" || errRes.Message != "" {
			err = errors.Wrap(errors.New(errRes.Err), errors.New(errRes.Message))
		}
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
		repoCall.Unset()
	}
}

func TestViewClient(t *testing.T) {
	us, svc := newUsersServer()
	defer us.Close()

	cases := []struct {
		desc   string
		token  string
		id     string
		status int
		err    error
	}{
		{
			desc:   "view user with valid token",
			token:  validToken,
			id:     client.ID,
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:   "view user with invalid token",
			token:  inValidToken,
			id:     client.ID,
			status: http.StatusUnauthorized,
			err:    svcerr.ErrAuthentication,
		},
		{
			desc:   "view user with empty token",
			token:  "",
			id:     client.ID,
			status: http.StatusUnauthorized,
			err:    apiutil.ErrBearerToken,
		},
	}

	for _, tc := range cases {
		req := testRequest{
			client: us.Client(),
			method: http.MethodGet,
			url:    fmt.Sprintf("%s/users/%s", us.URL, tc.id),
			token:  tc.token,
		}

		repoCall := svc.On("ViewClient", mock.Anything, tc.token, tc.id).Return(mgclients.Client{}, tc.err)
		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		var errRes respBody
		err = json.NewDecoder(res.Body).Decode(&errRes)
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error while decoding response body: %s", tc.desc, err))
		if errRes.Err != "" || errRes.Message != "" {
			err = errors.Wrap(errors.New(errRes.Err), errors.New(errRes.Message))
		}
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
		repoCall.Unset()
	}
}

func TestViewProfile(t *testing.T) {
	us, svc := newUsersServer()
	defer us.Close()

	cases := []struct {
		desc   string
		token  string
		id     string
		status int
		err    error
	}{
		{
			desc:   "view profile with valid token",
			token:  validToken,
			id:     client.ID,
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:   "view profile with invalid token",
			token:  inValidToken,
			id:     client.ID,
			status: http.StatusUnauthorized,
			err:    svcerr.ErrAuthentication,
		},
		{
			desc:   "view profile with empty token",
			token:  "",
			id:     client.ID,
			status: http.StatusUnauthorized,
			err:    apiutil.ErrBearerToken,
		},
	}

	for _, tc := range cases {
		req := testRequest{
			client: us.Client(),
			method: http.MethodGet,
			url:    fmt.Sprintf("%s/users/profile", us.URL),
			token:  tc.token,
		}

		repoCall := svc.On("ViewProfile", mock.Anything, tc.token).Return(mgclients.Client{}, tc.err)
		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		var errRes respBody
		err = json.NewDecoder(res.Body).Decode(&errRes)
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error while decoding response body: %s", tc.desc, err))
		if errRes.Err != "" || errRes.Message != "" {
			err = errors.Wrap(errors.New(errRes.Err), errors.New(errRes.Message))
		}
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
		repoCall.Unset()
	}
}

func TestListClients(t *testing.T) {
	us, svc := newUsersServer()
	defer us.Close()

	cases := []struct {
		desc              string
		query             string
		token             string
		listUsersResponse mgclients.ClientsPage
		status            int
		err               error
	}{
		{
			desc:   "list users with valid token",
			token:  validToken,
			status: http.StatusOK,
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Total: 1,
				},
				Clients: []mgclients.Client{client},
			},
			err: nil,
		},
		{
			desc:   "list users with empty token",
			token:  "",
			status: http.StatusUnauthorized,
			err:    apiutil.ErrBearerToken,
		},
		{
			desc:   "list users with invalid token",
			token:  inValidToken,
			status: http.StatusUnauthorized,
			err:    svcerr.ErrAuthentication,
		},
		{
			desc:  "list users with offset",
			token: validToken,
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Offset: 1,
					Total:  1,
				},
				Clients: []mgclients.Client{client},
			},
			query:  "offset=1",
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:   "list users with invalid offset",
			token:  validToken,
			query:  "offset=invalid",
			status: http.StatusBadRequest,
			err:    apiutil.ErrValidation,
		},
		{
			desc:  "list users with limit",
			token: validToken,
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Limit: 1,
					Total: 1,
				},
				Clients: []mgclients.Client{client},
			},
			query:  "limit=1",
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:   "list users with invalid limit",
			token:  validToken,
			query:  "limit=invalid",
			status: http.StatusBadRequest,
			err:    apiutil.ErrValidation,
		},
		{
			desc:   "list users with limit greater than max",
			token:  validToken,
			query:  fmt.Sprintf("limit=%d", api.MaxLimitSize+1),
			status: http.StatusBadRequest,
			err:    apiutil.ErrValidation,
		},
		{
			desc:  "list users with owner_id",
			token: validToken,
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Total: 1,
				},
				Clients: []mgclients.Client{client},
			},
			query:  fmt.Sprintf("owner_id=%s", validID),
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:   "list users with duplicate owner_id",
			token:  validToken,
			query:  "owner_id=1&owner_id=2",
			status: http.StatusBadRequest,
			err:    apiutil.ErrInvalidQueryParams,
		},
		{
			desc:   "list users with invalid owner_id",
			token:  validToken,
			query:  "owner_id=invalid",
			status: http.StatusBadRequest,
			err:    apiutil.ErrValidation,
		},
		{
			desc:  "list users with name",
			token: validToken,
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Total: 1,
				},
				Clients: []mgclients.Client{client},
			},
			query:  "name=clientname",
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:   "list users with invalid name",
			token:  validToken,
			query:  "name=invalid",
			status: http.StatusBadRequest,
			err:    apiutil.ErrValidation,
		},
		{
			desc:   "list users with duplicate name",
			token:  validToken,
			query:  "name=1&name=2",
			status: http.StatusBadRequest,
			err:    apiutil.ErrInvalidQueryParams,
		},
		{
			desc:  "list users with status",
			token: validToken,
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Total: 1,
				},
				Clients: []mgclients.Client{client},
			},
			query:  "status=enabled",
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:   "list users with invalid status",
			token:  validToken,
			query:  "status=invalid",
			status: http.StatusBadRequest,
			err:    apiutil.ErrValidation,
		},
		{
			desc:   "list users with duplicate status",
			token:  validToken,
			query:  "status=enabled&status=disabled",
			status: http.StatusBadRequest,
			err:    apiutil.ErrInvalidQueryParams,
		},
		{
			desc:  "list users with tags",
			token: validToken,
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Total: 1,
				},
				Clients: []mgclients.Client{client},
			},
			query:  "tag=tag1,tag2",
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:   "list users with invalid tags",
			token:  validToken,
			query:  "tag=invalid",
			status: http.StatusBadRequest,
			err:    apiutil.ErrValidation,
		},
		{
			desc:   "list users with duplicate tags",
			token:  validToken,
			query:  "tag=tag1&tag=tag2",
			status: http.StatusBadRequest,
			err:    apiutil.ErrInvalidQueryParams,
		},
		{
			desc:  "list users with metadata",
			token: validToken,
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Total: 1,
				},
				Clients: []mgclients.Client{client},
			},
			query:  "metadata=%7B%22domain%22%3A%20%22example.com%22%7D&",
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:   "list users with invalid metadata",
			token:  validToken,
			query:  "metadata=invalid",
			status: http.StatusBadRequest,
			err:    apiutil.ErrValidation,
		},
		{
			desc:   "list users with duplicate metadata",
			token:  validToken,
			query:  "metadata=%7B%22domain%22%3A%20%22example.com%22%7D&metadata=%7B%22domain%22%3A%20%22example.com%22%7D",
			status: http.StatusBadRequest,
			err:    apiutil.ErrInvalidQueryParams,
		},
		{
			desc:  "list users with permissions",
			token: validToken,
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Total: 1,
				},
				Clients: []mgclients.Client{client},
			},
			query:  "permission=view",
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:   "list users with invalid permissions",
			token:  validToken,
			query:  "permission=invalid",
			status: http.StatusBadRequest,
			err:    apiutil.ErrValidation,
		},
		{
			desc:   "list users with duplicate permissions",
			token:  validToken,
			query:  "permission=view&permission=view",
			status: http.StatusBadRequest,
			err:    apiutil.ErrInvalidQueryParams,
		},
		{
			desc:  "list users with list perms",
			token: validToken,
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Total: 1,
				},
				Clients: []mgclients.Client{client},
			},
			query:  "list_perms=true",
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:   "list users with invalid list perms",
			token:  validToken,
			query:  "list_perms=invalid",
			status: http.StatusBadRequest,
			err:    apiutil.ErrValidation,
		},
		{
			desc:   "list users with duplicate list perms",
			token:  validToken,
			query:  "list_perms=true&list_perms=true",
			status: http.StatusBadRequest,
			err:    apiutil.ErrInvalidQueryParams,
		},
		{
			desc:  "list users with identity",
			token: validToken,
			query: fmt.Sprintf("identity=%s", client.Credentials.Identity),
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Total: 1,
				},
				Clients: []mgclients.Client{
					client,
				},
			},
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:   "list users with invalid identity",
			token:  validToken,
			query:  "identity=invalid",
			status: http.StatusBadRequest,
			err:    apiutil.ErrValidation,
		},
		{
			desc:   "list users with duplicate identity",
			token:  validToken,
			query:  "identity=1&identity=2",
			status: http.StatusBadRequest,
			err:    apiutil.ErrInvalidQueryParams,
		},
		{
			desc:  "list users with  mine visibility",
			token: validToken,
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Total: 1,
				},
				Clients: []mgclients.Client{
					client,
				},
			},
			query:  "visibility=mine",
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:  "list users with shared visisbility",
			token: validToken,
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Total: 1,
				},
				Clients: []mgclients.Client{
					client,
				},
			},
			query:  "visibility=shared",
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:  "list users with all visibility",
			token: validToken,
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Total: 1,
				},
				Clients: []mgclients.Client{
					client,
				},
			},
			query:  "visibility=all",
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:   "list users with invalid visibility",
			token:  validToken,
			query:  "visibility=invalid",
			status: http.StatusBadRequest,
			err:    apiutil.ErrValidation,
		},
		{
			desc:   "list users with duplicate visibility",
			token:  validToken,
			query:  "visibility=mine&visibility=shared",
			status: http.StatusBadRequest,
			err:    apiutil.ErrInvalidQueryParams,
		},
		{
			desc: "list users with order",
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Total: 1,
				},
				Clients: []mgclients.Client{
					client,
				},
			},
			token:  validToken,
			query:  "order=name",
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:   "list users with invalid order",
			token:  validToken,
			query:  "order=invalid",
			status: http.StatusBadRequest,
			err:    apiutil.ErrValidation,
		},
		{
			desc:   "list users with duplicate order",
			token:  validToken,
			query:  "order=name&order=name",
			status: http.StatusBadRequest,
			err:    apiutil.ErrInvalidQueryParams,
		},
		{
			desc:   "list users with invalid order direction",
			token:  validToken,
			query:  "dir=invalid",
			status: http.StatusBadRequest,
			err:    apiutil.ErrValidation,
		},
		{
			desc:   "list users with duplicate order direction",
			token:  validToken,
			query:  "dir=asc&dir=asc",
			status: http.StatusBadRequest,
			err:    apiutil.ErrInvalidQueryParams,
		},
	}

	for _, tc := range cases {
		req := testRequest{
			client:      us.Client(),
			method:      http.MethodGet,
			url:         us.URL + "/users?" + tc.query,
			contentType: contentType,
			token:       tc.token,
		}

		repoCall := svc.On("ListClients", mock.Anything, tc.token, mock.Anything, mock.Anything).Return(tc.listUsersResponse, tc.err)
		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		var bodyRes respBody
		err = json.NewDecoder(res.Body).Decode(&bodyRes)
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error while decoding response body: %s", tc.desc, err))
		if bodyRes.Err != "" || bodyRes.Message != "" {
			err = errors.Wrap(errors.New(bodyRes.Err), errors.New(bodyRes.Message))
		}
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
		repoCall.Unset()
	}
}

func TestUpdateClient(t *testing.T) {
	us, svc := newUsersServer()
	defer us.Close()

	newName := "newname"
	newMetadata := mgclients.Metadata{"newkey": "newvalue"}

	cases := []struct {
		desc           string
		id             string
		data           string
		clientResponse mgclients.Client
		token          string
		contentType    string
		status         int
		err            error
	}{
		{
			desc:        "update user with valid token",
			id:          client.ID,
			data:        fmt.Sprintf(`{"name":"%s","metadata":%s}`, newName, toJSON(newMetadata)),
			token:       validToken,
			contentType: contentType,
			clientResponse: mgclients.Client{
				ID:       client.ID,
				Name:     newName,
				Metadata: newMetadata,
			},
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:        "update user with invalid token",
			id:          client.ID,
			data:        fmt.Sprintf(`{"name":"%s","metadata":%s}`, newName, toJSON(newMetadata)),
			token:       inValidToken,
			contentType: contentType,
			status:      http.StatusUnauthorized,
			err:         svcerr.ErrAuthentication,
		},
		{
			desc:        "update user with empty token",
			id:          client.ID,
			data:        fmt.Sprintf(`{"name":"%s","metadata":%s}`, newName, toJSON(newMetadata)),
			token:       "",
			contentType: contentType,
			status:      http.StatusUnauthorized,
			err:         apiutil.ErrBearerToken,
		},
		{
			desc:        "update user with invalid id",
			id:          inValid,
			data:        fmt.Sprintf(`{"name":"%s","metadata":%s}`, newName, toJSON(newMetadata)),
			token:       validToken,
			contentType: contentType,
			status:      http.StatusForbidden,
			err:         svcerr.ErrAuthorization,
		},
		{
			desc:        "update user with invalid contentype",
			id:          client.ID,
			data:        fmt.Sprintf(`{"name":"%s","metadata":%s}`, newName, toJSON(newMetadata)),
			token:       validToken,
			contentType: "application/xml",
			status:      http.StatusUnsupportedMediaType,
			err:         apiutil.ErrValidation,
		},
		{
			desc:        "update user with malformed data",
			id:          client.ID,
			data:        fmt.Sprintf(`{"name":%s}`, "invalid"),
			token:       validToken,
			contentType: contentType,
			status:      http.StatusBadRequest,
			err:         apiutil.ErrValidation,
		},
		{
			desc:        "update user with empty id",
			id:          " ",
			data:        fmt.Sprintf(`{"name":"%s","metadata":%s}`, newName, toJSON(newMetadata)),
			token:       validToken,
			contentType: contentType,
			status:      http.StatusBadRequest,
			err:         apiutil.ErrValidation,
		},
	}

	for _, tc := range cases {
		req := testRequest{
			client:      us.Client(),
			method:      http.MethodPatch,
			url:         fmt.Sprintf("%s/users/%s", us.URL, tc.id),
			contentType: tc.contentType,
			token:       tc.token,
			body:        strings.NewReader(tc.data),
		}
		repoCall := svc.On("UpdateClient", mock.Anything, tc.token, mock.Anything).Return(tc.clientResponse, tc.err)
		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		var resBody respBody
		err = json.NewDecoder(res.Body).Decode(&resBody)
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error while decoding response body: %s", tc.desc, err))
		if resBody.Err != "" || resBody.Message != "" {
			err = errors.Wrap(errors.New(resBody.Err), errors.New(resBody.Message))
		}
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
		repoCall.Unset()
	}
}

func TestUpdateClientTags(t *testing.T) {
	us, svc := newUsersServer()
	defer us.Close()

	newTag := "newtag"

	cases := []struct {
		desc           string
		id             string
		data           string
		contentType    string
		clientResponse mgclients.Client
		token          string
		status         int
		err            error
	}{
		{
			desc:        "update user tags with valid token",
			id:          client.ID,
			data:        fmt.Sprintf(`{"tags":["%s"]}`, newTag),
			contentType: contentType,
			clientResponse: mgclients.Client{
				ID:   client.ID,
				Tags: []string{newTag},
			},
			token:  validToken,
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:        "update user tags with empty token",
			id:          client.ID,
			data:        fmt.Sprintf(`{"tags":["%s"]}`, newTag),
			contentType: contentType,
			token:       "",
			status:      http.StatusUnauthorized,
			err:         apiutil.ErrBearerToken,
		},
		{
			desc:        "update user tags with invalid token",
			id:          client.ID,
			data:        fmt.Sprintf(`{"tags":["%s"]}`, newTag),
			contentType: contentType,
			token:       inValidToken,
			status:      http.StatusUnauthorized,
			err:         svcerr.ErrAuthentication,
		},
		{
			desc:        "update user tags with invalid id",
			id:          client.ID,
			data:        fmt.Sprintf(`{"tags":["%s"]}`, newTag),
			contentType: contentType,
			token:       validToken,
			status:      http.StatusForbidden,
			err:         svcerr.ErrAuthorization,
		},
		{
			desc:        "update user tags with invalid contentype",
			id:          client.ID,
			data:        fmt.Sprintf(`{"tags":["%s"]}`, newTag),
			contentType: "application/xml",
			token:       validToken,
			status:      http.StatusUnsupportedMediaType,
			err:         apiutil.ErrValidation,
		},
		{
			desc:        "update user tags with empty id",
			id:          "",
			data:        fmt.Sprintf(`{"tags":["%s"]}`, newTag),
			contentType: contentType,
			token:       validToken,
			status:      http.StatusBadRequest,
			err:         apiutil.ErrValidation,
		},
		{
			desc:        "update user with malfomed data",
			id:          client.ID,
			data:        fmt.Sprintf(`{"tags":%s}`, newTag),
			contentType: contentType,
			token:       validToken,
			status:      http.StatusBadRequest,
			err:         apiutil.ErrValidation,
		},
	}

	for _, tc := range cases {
		req := testRequest{
			client:      us.Client(),
			method:      http.MethodPatch,
			url:         fmt.Sprintf("%s/users/%s/tags", us.URL, tc.id),
			contentType: tc.contentType,
			token:       tc.token,
			body:        strings.NewReader(tc.data),
		}

		repoCall := svc.On("UpdateClientTags", mock.Anything, tc.token, mock.Anything).Return(tc.clientResponse, tc.err)
		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		var resBody respBody
		err = json.NewDecoder(res.Body).Decode(&resBody)
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error while decoding response body: %s", tc.desc, err))
		if resBody.Err != "" || resBody.Message != "" {
			err = errors.Wrap(errors.New(resBody.Err), errors.New(resBody.Message))
		}
		if err == nil {
			assert.Equal(t, tc.clientResponse.Tags, resBody.Tags, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.clientResponse.Tags, resBody.Tags))
		}
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
		repoCall.Unset()
	}
}

func TestUpdateClientIdentity(t *testing.T) {
	us, svc := newUsersServer()
	defer us.Close()

	cases := []struct {
		desc        string
		data        string
		client      mgclients.Client
		contentType string
		token       string
		status      int
		err         error
	}{
		{
			desc: "update user identity with valid token",
			data: fmt.Sprintf(`{"identity": "%s"}`, "newclientidentity@example.com"),
			client: mgclients.Client{
				ID: client.ID,
				Credentials: mgclients.Credentials{
					Identity: "newclientidentity@example.com",
					Secret:   "secret",
				},
			},
			contentType: contentType,
			token:       validToken,
			status:      http.StatusOK,
			err:         nil,
		},
		{
			desc: "update user identity with empty token",
			data: fmt.Sprintf(`{"identity": "%s"}`, "newclientidentity@example.com"),
			client: mgclients.Client{
				ID: client.ID,
				Credentials: mgclients.Credentials{
					Identity: "newclientidentity@example.com",
					Secret:   "secret",
				},
			},
			contentType: contentType,
			token:       "",
			status:      http.StatusUnauthorized,
			err:         apiutil.ErrBearerToken,
		},
		{
			desc: "update user identity with invalid token",
			data: fmt.Sprintf(`{"identity": "%s"}`, "newclientidentity@example.com"),
			client: mgclients.Client{
				ID: client.ID,
				Credentials: mgclients.Credentials{
					Identity: "newclientidentity@example.com",
					Secret:   "secret",
				},
			},
			contentType: contentType,
			token:       inValid,
			status:      http.StatusUnauthorized,
			err:         svcerr.ErrAuthentication,
		},
		{
			desc: "update user identity with empty id",
			data: fmt.Sprintf(`{"identity": "%s"}`, "newclientidentity@example.com"),
			client: mgclients.Client{
				ID: "",
				Credentials: mgclients.Credentials{
					Identity: "newclientidentity@example.com",
					Secret:   "secret",
				},
			},
			contentType: contentType,
			token:       validToken,
			status:      http.StatusBadRequest,
			err:         apiutil.ErrMissingID,
		},
		{
			desc: "update user identity with invalid contentype",
			data: fmt.Sprintf(`{"identity": "%s"}`, ""),
			client: mgclients.Client{
				ID: client.ID,
				Credentials: mgclients.Credentials{
					Identity: "newclientidentity@example.com",
					Secret:   "secret",
				},
			},
			contentType: "application/xml",
			token:       validToken,
			status:      http.StatusUnsupportedMediaType,
			err:         apiutil.ErrValidation,
		},
		{
			desc: "update user identity with malformed data",
			data: fmt.Sprintf(`{"identity": %s}`, "invalid"),
			client: mgclients.Client{
				ID: client.ID,
				Credentials: mgclients.Credentials{
					Identity: "",
					Secret:   "secret",
				},
			},
			contentType: contentType,
			token:       validToken,
			status:      http.StatusBadRequest,
			err:         apiutil.ErrValidation,
		},
	}

	for _, tc := range cases {
		req := testRequest{
			client:      us.Client(),
			method:      http.MethodPatch,
			url:         fmt.Sprintf("%s/users/%s/identity", us.URL, tc.client.ID),
			contentType: tc.contentType,
			token:       tc.token,
			body:        strings.NewReader(tc.data),
		}

		repoCall := svc.On("UpdateClientIdentity", mock.Anything, tc.token, mock.Anything, mock.Anything).Return(mgclients.Client{}, tc.err)
		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		var resBody respBody
		err = json.NewDecoder(res.Body).Decode(&resBody)
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error while decoding response body: %s", tc.desc, err))
		if resBody.Err != "" || resBody.Message != "" {
			err = errors.Wrap(errors.New(resBody.Err), errors.New(resBody.Message))
		}
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
		repoCall.Unset()
	}
}

func TestPasswordResetRequest(t *testing.T) {
	us, svc := newUsersServer()
	defer us.Close()

	testemail := "test@example.com"
	testhost := "example.com"

	cases := []struct {
		desc        string
		data        string
		contentType string
		status      int
		err         error
	}{
		{
			desc:        "password reset request with valid email",
			data:        fmt.Sprintf(`{"email": "%s", "host": "%s"}`, testemail, testhost),
			contentType: contentType,
			status:      http.StatusCreated,
			err:         nil,
		},
		{
			desc:        "password reset request with empty email",
			data:        fmt.Sprintf(`{"email": "%s", "host": "%s"}`, "", testhost),
			contentType: contentType,
			status:      http.StatusInternalServerError,
			err:         apiutil.ErrValidation,
		},
		{
			desc:        "password reset request with empty host",
			data:        fmt.Sprintf(`{"email": "%s", "host": "%s"}`, testemail, ""),
			contentType: contentType,
			status:      http.StatusBadRequest,
			err:         apiutil.ErrValidation,
		},
		{
			desc:        "password reset request with invalid email",
			data:        fmt.Sprintf(`{"email": "%s", "host": "%s"}`, "invalid", testhost),
			contentType: contentType,
			status:      http.StatusNotFound,
			err:         errors.ErrNotFound,
		},
		{
			desc:        "password reset with malformed data",
			data:        fmt.Sprintf(`{"email": %s, "host": %s}`, testemail, testhost),
			contentType: contentType,
			status:      http.StatusBadRequest,
			err:         apiutil.ErrValidation,
		},
		{
			desc:        "password reset with invalid contentype",
			data:        fmt.Sprintf(`{"email": "%s", "host": "%s"}`, testemail, testhost),
			contentType: "application/xml",
			status:      http.StatusUnsupportedMediaType,
			err:         apiutil.ErrValidation,
		},
	}

	for _, tc := range cases {
		req := testRequest{
			client:      us.Client(),
			method:      http.MethodPost,
			url:         fmt.Sprintf("%s/users/password/reset-request", us.URL),
			contentType: tc.contentType,
			body:        strings.NewReader(tc.data),
		}

		repoCall := svc.On("GenerateResetToken", mock.Anything, mock.Anything, mock.Anything).Return(tc.err)
		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
		repoCall.Unset()
	}
}

func TestPasswordReset(t *testing.T) {
	us, svc := newUsersServer()
	defer us.Close()

	strongPass := "StrongPassword"

	cases := []struct {
		desc        string
		data        string
		token       string
		contentType string
		status      int
		err         error
	}{
		{
			desc:        "password reset with valid token",
			data:        fmt.Sprintf(`{"token": "%s", "password": "%s", "confirm_password": "%s"}`, validToken, strongPass, strongPass),
			token:       validToken,
			contentType: contentType,
			status:      http.StatusCreated,
			err:         nil,
		},
		{
			desc:        "password reset with invalid token",
			data:        fmt.Sprintf(`{"token": "%s", "password": "%s", "confirm_password": "%s"}`, inValidToken, strongPass, strongPass),
			token:       inValidToken,
			contentType: contentType,
			status:      http.StatusUnauthorized,
			err:         svcerr.ErrAuthentication,
		},
		{
			desc:        "password reset to weak password",
			data:        fmt.Sprintf(`{"token": "%s", "password": "%s", "confirm_password": "%s"}`, validToken, "weak", "weak"),
			token:       validToken,
			contentType: contentType,
			status:      http.StatusInternalServerError,
			err:         ErrPasswordFormat,
		},
		{
			desc:        "password reset with empty token",
			data:        fmt.Sprintf(`{"token": "%s", "password": "%s", "confirm_password": "%s"}`, "", strongPass, strongPass),
			token:       "",
			contentType: contentType,
			status:      http.StatusUnauthorized,
			err:         apiutil.ErrBearerToken,
		},
		{
			desc:        "password reset with empty password",
			data:        fmt.Sprintf(`{"token": "%s", "password": "%s", "confirm_password": "%s"}`, validToken, "", ""),
			token:       validToken,
			contentType: contentType,
			status:      http.StatusBadRequest,
			err:         apiutil.ErrValidation,
		},
		{
			desc:        "password reset with malformed data",
			data:        fmt.Sprintf(`{"token": "%s", "password": %s, "confirm_password": %s}`, validToken, strongPass, strongPass),
			token:       validToken,
			contentType: contentType,
			status:      http.StatusBadRequest,
			err:         apiutil.ErrValidation,
		},
		{
			desc:   "password reset with invalid contentype",
			data:   fmt.Sprintf(`{"token": "%s", "password": "%s", "confirm_password": "%s"}`, validToken, strongPass, strongPass),
			token:  validToken,
			status: http.StatusUnsupportedMediaType,
			err:    apiutil.ErrValidation,
		},
	}

	for _, tc := range cases {
		req := testRequest{
			client:      us.Client(),
			method:      http.MethodPut,
			url:         fmt.Sprintf("%s/users/password/reset", us.URL),
			contentType: tc.contentType,
			token:       tc.token,
			body:        strings.NewReader(tc.data),
		}

		repoCall := svc.On("ResetSecret", mock.Anything, mock.Anything, mock.Anything).Return(tc.err)
		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
		repoCall.Unset()
	}
}

func TestUpdateClientRole(t *testing.T) {
	us, svc := newUsersServer()
	defer us.Close()

	cases := []struct {
		desc        string
		data        string
		clientID    string
		token       string
		contentType string
		status      int
		err         error
	}{
		{
			desc:        "update client role with valid token",
			data:        fmt.Sprintf(`{"role": "%s"}`, "admin"),
			clientID:    client.ID,
			token:       validToken,
			contentType: contentType,
			status:      http.StatusOK,
			err:         nil,
		},
		{
			desc:        "update client role with invalid token",
			data:        fmt.Sprintf(`{"role": "%s"}`, "admin"),
			clientID:    client.ID,
			token:       inValidToken,
			contentType: contentType,
			status:      http.StatusUnauthorized,
			err:         svcerr.ErrAuthentication,
		},
		{
			desc:        "update client role with invalid id",
			data:        fmt.Sprintf(`{"role": "%s"}`, "admin"),
			clientID:    inValid,
			token:       validToken,
			contentType: contentType,
			status:      http.StatusForbidden,
			err:         svcerr.ErrAuthorization,
		},
		{
			desc:        "update client role with empty token",
			data:        fmt.Sprintf(`{"role": "%s"}`, "admin"),
			clientID:    client.ID,
			token:       "",
			contentType: contentType,
			status:      http.StatusUnauthorized,
			err:         apiutil.ErrBearerToken,
		},
		{
			desc:        "update client with invalid role",
			data:        fmt.Sprintf(`{"role": "%s"}`, "invalid"),
			clientID:    client.ID,
			token:       validToken,
			contentType: contentType,
			status:      http.StatusInternalServerError,
			err:         svcerr.ErrInvalidRole,
		},
		{
			desc:        "update client with invalid contentype",
			data:        fmt.Sprintf(`{"role": "%s"}`, "admin"),
			clientID:    client.ID,
			token:       validToken,
			contentType: "application/xml",
			status:      http.StatusUnsupportedMediaType,
			err:         apiutil.ErrValidation,
		},
		{
			desc:        "update client with malformed data",
			data:        fmt.Sprintf(`{"role": %s}`, "admin"),
			clientID:    client.ID,
			token:       validToken,
			contentType: contentType,
			status:      http.StatusBadRequest,
			err:         apiutil.ErrValidation,
		},
	}

	for _, tc := range cases {
		req := testRequest{
			client:      us.Client(),
			method:      http.MethodPatch,
			url:         fmt.Sprintf("%s/users/%s/role", us.URL, tc.clientID),
			contentType: tc.contentType,
			token:       tc.token,
			body:        strings.NewReader(tc.data),
		}

		repoCall := svc.On("UpdateClientRole", mock.Anything, tc.token, mock.Anything).Return(mgclients.Client{}, tc.err)
		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		var resBody respBody
		err = json.NewDecoder(res.Body).Decode(&resBody)
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error while decoding response body: %s", tc.desc, err))
		if resBody.Err != "" || resBody.Message != "" {
			err = errors.Wrap(errors.New(resBody.Err), errors.New(resBody.Message))
		}
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
		repoCall.Unset()
	}
}

func TestUpdateClientSecret(t *testing.T) {
	us, svc := newUsersServer()
	defer us.Close()

	cases := []struct {
		desc        string
		data        string
		client      mgclients.Client
		contentType string
		token       string
		status      int
		err         error
	}{
		{
			desc: "update user secret with valid token",
			data: fmt.Sprintf(`{"secret": "%s"}`, "strongersecret"),
			client: mgclients.Client{
				ID: client.ID,
				Credentials: mgclients.Credentials{
					Identity: "clientname",
					Secret:   "strongersecret",
				},
			},
			contentType: contentType,
			token:       validToken,
			status:      http.StatusOK,
			err:         nil,
		},
		{
			desc: "update user secret with empty token",
			data: fmt.Sprintf(`{"secret": "%s"}`, "strongersecret"),
			client: mgclients.Client{
				ID: client.ID,
				Credentials: mgclients.Credentials{
					Identity: "clientname",
					Secret:   "strongersecret",
				},
			},
			contentType: contentType,
			token:       "",
			status:      http.StatusUnauthorized,
			err:         apiutil.ErrBearerToken,
		},
		{
			desc: "update user secret with invalid token",
			data: fmt.Sprintf(`{"secret": "%s"}`, "strongersecret"),
			client: mgclients.Client{
				ID: client.ID,
				Credentials: mgclients.Credentials{
					Identity: "clientname",
					Secret:   "strongersecret",
				},
			},
			contentType: contentType,
			token:       inValid,
			status:      http.StatusUnauthorized,
			err:         svcerr.ErrAuthentication,
		},

		{
			desc: "update user secret with empty secret",
			data: fmt.Sprintf(`{"secret": "%s"}`, ""),
			client: mgclients.Client{
				ID: client.ID,
				Credentials: mgclients.Credentials{
					Identity: "clientname",
					Secret:   "",
				},
			},
			contentType: contentType,
			token:       validToken,
			status:      http.StatusBadRequest,
			err:         apiutil.ErrBearerKey,
		},
		{
			desc: "update user secret with invalid contentype",
			data: fmt.Sprintf(`{"secret": "%s"}`, ""),
			client: mgclients.Client{
				ID: client.ID,
				Credentials: mgclients.Credentials{
					Identity: "clientname",
					Secret:   "",
				},
			},
			contentType: "application/xml",
			token:       validToken,
			status:      http.StatusUnsupportedMediaType,
			err:         apiutil.ErrValidation,
		},
		{
			desc: "update user secret with malformed data",
			data: fmt.Sprintf(`{"secret": %s}`, "invalid"),
			client: mgclients.Client{
				ID: client.ID,
				Credentials: mgclients.Credentials{
					Identity: "clientname",
					Secret:   "",
				},
			},
			contentType: contentType,
			token:       validToken,
			status:      http.StatusBadRequest,
			err:         apiutil.ErrValidation,
		},
	}

	for _, tc := range cases {
		req := testRequest{
			client:      us.Client(),
			method:      http.MethodPatch,
			url:         fmt.Sprintf("%s/users/secret", us.URL),
			contentType: tc.contentType,
			token:       tc.token,
			body:        strings.NewReader(tc.data),
		}

		repoCall := svc.On("UpdateClientSecret", mock.Anything, tc.token, mock.Anything, mock.Anything).Return(tc.client, tc.err)

		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		var resBody respBody
		err = json.NewDecoder(res.Body).Decode(&resBody)
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error while decoding response body: %s", tc.desc, err))
		if resBody.Err != "" || resBody.Message != "" {
			err = errors.Wrap(errors.New(resBody.Err), errors.New(resBody.Message))
		}
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
		repoCall.Unset()
	}
}

func TestIssueToken(t *testing.T) {
	us, svc := newUsersServer()
	defer us.Close()

	validIdentity := "valid"

	cases := []struct {
		desc        string
		data        string
		contentType string
		status      int
		err         error
	}{
		{
			desc:        "issue token with valid identity and secret",
			data:        fmt.Sprintf(`{"identity": "%s", "secret": "%s", "domainID": "%s"}`, validIdentity, secret, validID),
			contentType: contentType,
			status:      http.StatusCreated,
			err:         nil,
		},
		{
			desc:        "issue token with empty identity",
			data:        fmt.Sprintf(`{"identity": "%s", "secret": "%s", "domainID": "%s"}`, "", secret, validID),
			contentType: contentType,
			status:      http.StatusInternalServerError,
			err:         apiutil.ErrValidation,
		},
		{
			desc:        "issue token with empty secret",
			data:        fmt.Sprintf(`{"identity": "%s", "secret": "%s", "domainID": "%s"}`, validIdentity, "", validID),
			contentType: contentType,
			status:      http.StatusBadRequest,
			err:         apiutil.ErrValidation,
		},
		{
			desc:        "issue token with empty domain",
			data:        fmt.Sprintf(`{"identity": "%s", "secret": "%s", "domainID": "%s"}`, validIdentity, secret, ""),
			contentType: contentType,
			status:      http.StatusBadRequest,
			err:         apiutil.ErrValidation,
		},
		{
			desc:        "issue token with invalid identity",
			data:        fmt.Sprintf(`{"identity": "%s", "secret": "%s", "domainID": "%s"}`, "invalid", secret, validID),
			contentType: contentType,
			status:      http.StatusUnauthorized,
			err:         svcerr.ErrAuthentication,
		},
		{
			desc:        "issues token with malformed data",
			data:        fmt.Sprintf(`{"identity": %s, "secret": %s, "domainID": %s}`, validIdentity, secret, validID),
			contentType: contentType,
			status:      http.StatusBadRequest,
			err:         apiutil.ErrValidation,
		},
		{
			desc:        "issue token with invalid contentype",
			data:        fmt.Sprintf(`{"identity": "%s", "secret": "%s", "domainID": "%s"}`, "invalid", secret, validID),
			contentType: "application/xml",
			status:      http.StatusUnsupportedMediaType,
			err:         apiutil.ErrValidation,
		},
	}

	for _, tc := range cases {
		req := testRequest{
			client:      us.Client(),
			method:      http.MethodPost,
			url:         fmt.Sprintf("%s/users/tokens/issue", us.URL),
			contentType: tc.contentType,
			body:        strings.NewReader(tc.data),
		}

		repoCall := svc.On("IssueToken", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&magistrala.Token{}, tc.err)
		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		if tc.err != nil {
			var resBody respBody
			err = json.NewDecoder(res.Body).Decode(&resBody)
			assert.Nil(t, err, fmt.Sprintf("%s: unexpected error while decoding response body: %s", tc.desc, err))
			if resBody.Err != "" || resBody.Message != "" {
				err = errors.Wrap(errors.New(resBody.Err), errors.New(resBody.Message))
			}
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
		}
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
		repoCall.Unset()
	}
}

func TestRefreshToken(t *testing.T) {
	us, svc := newUsersServer()
	defer us.Close()

	cases := []struct {
		desc        string
		data        string
		contentType string
		status      int
		err         error
	}{
		{
			desc:        "refresh token with valid token",
			data:        fmt.Sprintf(`{"refresh_token": "%s", "domain_id": "%s"}`, validToken, validID),
			contentType: contentType,
			status:      http.StatusCreated,
			err:         nil,
		},
		{
			desc:        "refresh token with invalid token",
			data:        fmt.Sprintf(`{"refresh_token": "%s", "domain_id": "%s"}`, inValidToken, validID),
			contentType: contentType,
			status:      http.StatusUnauthorized,
			err:         svcerr.ErrAuthentication,
		},
		{
			desc:        "refresh token with empty token",
			data:        fmt.Sprintf(`{"refresh_token": "%s", "domain_id": "%s"}`, "", validID),
			contentType: contentType,
			status:      http.StatusUnauthorized,
			err:         apiutil.ErrValidation,
		},
		{
			desc:        "refresh token with invalid domain",
			data:        fmt.Sprintf(`{"refresh_token": "%s", "domain_id": "%s"}`, validToken, "invalid"),
			contentType: contentType,
			status:      http.StatusUnauthorized,
			err:         svcerr.ErrAuthentication,
		},
		{
			desc:        "refresh token with malformed data",
			data:        fmt.Sprintf(`{"refresh_token": %s, "domain_id": %s}`, validToken, validID),
			contentType: contentType,
			status:      http.StatusBadRequest,
			err:         apiutil.ErrValidation,
		},
		{
			desc:        "refresh token with invalid contentype",
			data:        fmt.Sprintf(`{"refresh_token": "%s", "domain_id": "%s"}`, validToken, validID),
			contentType: "application/xml",
			status:      http.StatusUnsupportedMediaType,
			err:         apiutil.ErrValidation,
		},
	}

	for _, tc := range cases {
		req := testRequest{
			client:      us.Client(),
			method:      http.MethodPost,
			url:         fmt.Sprintf("%s/users/tokens/refresh", us.URL),
			contentType: tc.contentType,
			body:        strings.NewReader(tc.data),
		}

		repoCall := svc.On("RefreshToken", mock.Anything, mock.Anything, mock.Anything).Return(&magistrala.Token{}, tc.err)
		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		if tc.err != nil {
			var resBody respBody
			err = json.NewDecoder(res.Body).Decode(&resBody)
			assert.Nil(t, err, fmt.Sprintf("%s: unexpected error while decoding response body: %s", tc.desc, err))
			if resBody.Err != "" || resBody.Message != "" {
				err = errors.Wrap(errors.New(resBody.Err), errors.New(resBody.Message))
			}
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
		}
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
		repoCall.Unset()
	}
}

func TestEnableClient(t *testing.T) {
	us, svc := newUsersServer()
	defer us.Close()
	cases := []struct {
		desc     string
		client   mgclients.Client
		response mgclients.Client
		token    string
		status   int
		err      error
	}{
		{
			desc:   "enable client with valid token",
			client: client,
			response: mgclients.Client{
				ID:     client.ID,
				Status: mgclients.EnabledStatus,
			},
			token:  validToken,
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:   "enable client with invalid token",
			client: client,
			token:  inValidToken,
			status: http.StatusUnauthorized,
			err:    svcerr.ErrAuthentication,
		},
		{
			desc: "enable client with empty id",
			client: mgclients.Client{
				ID: "",
			},
			token:  validToken,
			status: http.StatusBadRequest,
			err:    apiutil.ErrMissingID,
		},
		{
			desc: "enable client with invalid id",
			client: mgclients.Client{
				ID: "invalid",
			},
			token:  validToken,
			status: http.StatusForbidden,
			err:    svcerr.ErrAuthorization,
		},
	}

	for _, tc := range cases {
		data := toJSON(tc.client)
		req := testRequest{
			client:      us.Client(),
			method:      http.MethodPost,
			url:         fmt.Sprintf("%s/users/%s/enable", us.URL, tc.client.ID),
			contentType: contentType,
			token:       tc.token,
			body:        strings.NewReader(data),
		}

		repoCall := svc.On("EnableClient", mock.Anything, mock.Anything, mock.Anything).Return(mgclients.Client{}, tc.err)
		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		if tc.err != nil {
			var resBody respBody
			err = json.NewDecoder(res.Body).Decode(&resBody)
			assert.Nil(t, err, fmt.Sprintf("%s: unexpected error while decoding response body: %s", tc.desc, err))
			if resBody.Err != "" || resBody.Message != "" {
				err = errors.Wrap(errors.New(resBody.Err), errors.New(resBody.Message))
			}
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
		}
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
		repoCall.Unset()
	}
}

func TestDisableClient(t *testing.T) {
	us, svc := newUsersServer()
	defer us.Close()

	cases := []struct {
		desc     string
		client   mgclients.Client
		response mgclients.Client
		token    string
		status   int
		err      error
	}{
		{
			desc:   "disable user with valid token",
			client: client,
			response: mgclients.Client{
				ID:     client.ID,
				Status: mgclients.DisabledStatus,
			},
			token:  validToken,
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:   "disable user with invalid token",
			client: client,
			token:  inValidToken,
			status: http.StatusUnauthorized,
			err:    svcerr.ErrAuthentication,
		},
		{
			desc: "disable user with empty id",
			client: mgclients.Client{
				ID: "",
			},
			token:  validToken,
			status: http.StatusBadRequest,
			err:    apiutil.ErrMissingID,
		},
		{
			desc: "disable user with invalid id",
			client: mgclients.Client{
				ID: "invalid",
			},
			token:  validToken,
			status: http.StatusForbidden,
			err:    svcerr.ErrAuthorization,
		},
	}

	for _, tc := range cases {
		data := toJSON(tc.client)
		req := testRequest{
			client:      us.Client(),
			method:      http.MethodPost,
			url:         fmt.Sprintf("%s/users/%s/disable", us.URL, tc.client.ID),
			contentType: contentType,
			token:       tc.token,
			body:        strings.NewReader(data),
		}

		repoCall := svc.On("DisableClient", mock.Anything, mock.Anything, mock.Anything).Return(mgclients.Client{}, tc.err)
		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
		repoCall.Unset()
	}
}

func TestListUsersByUserGroupId(t *testing.T) {
	us, svc := newUsersServer()
	defer us.Close()

	cases := []struct {
		desc              string
		token             string
		groupID           string
		page              mgclients.Page
		status            int
		query             string
		listUsersResponse mgclients.ClientsPage
		err               error
	}{
		{
			desc:    "list users with valid token",
			token:   validToken,
			groupID: validID,
			status:  http.StatusOK,
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Total: 1,
				},
				Clients: []mgclients.Client{client},
			},
			err: nil,
		},
		{
			desc:    "list users with empty id",
			token:   validToken,
			groupID: "",
			status:  http.StatusBadRequest,
			err:     apiutil.ErrMissingID,
		},
		{
			desc:    "list users with empty token",
			token:   "",
			groupID: validID,
			status:  http.StatusUnauthorized,
			err:     apiutil.ErrBearerToken,
		},
		{
			desc:    "list users with invalid token",
			token:   inValidToken,
			groupID: validID,
			status:  http.StatusUnauthorized,
			err:     svcerr.ErrAuthentication,
		},
		{
			desc:    "list users with offset",
			token:   validToken,
			groupID: validID,
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Offset: 1,
					Total:  1,
				},
				Clients: []mgclients.Client{client},
			},
			query:  "offset=1",
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:    "list users with invalid offset",
			token:   validToken,
			groupID: validID,
			query:   "offset=invalid",
			status:  http.StatusBadRequest,
			err:     apiutil.ErrValidation,
		},
		{
			desc:    "list users with limit",
			token:   validToken,
			groupID: validID,
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Limit: 1,
					Total: 1,
				},
				Clients: []mgclients.Client{client},
			},
			query:  "limit=1",
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:    "list users with invalid limit",
			token:   validToken,
			groupID: validID,
			query:   "limit=invalid",
			status:  http.StatusBadRequest,
			err:     apiutil.ErrValidation,
		},
		{
			desc:    "list users with limit greater than max",
			token:   validToken,
			groupID: validID,
			query:   fmt.Sprintf("limit=%d", api.MaxLimitSize+1),
			status:  http.StatusBadRequest,
			err:     apiutil.ErrValidation,
		},
		{
			desc:    "list users with owner_id",
			token:   validToken,
			groupID: validID,
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Total: 1,
				},
				Clients: []mgclients.Client{client},
			},
			query:  fmt.Sprintf("owner_id=%s", validID),
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:    "list users with duplicate owner_id",
			token:   validToken,
			groupID: validID,
			query:   "owner_id=1&owner_id=2",
			status:  http.StatusBadRequest,
			err:     apiutil.ErrInvalidQueryParams,
		},
		{
			desc:    "list users with invalid owner_id",
			token:   validToken,
			groupID: validID,
			query:   "owner_id=invalid",
			status:  http.StatusBadRequest,
			err:     apiutil.ErrValidation,
		},
		{
			desc:    "list users with name",
			token:   validToken,
			groupID: validID,
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Total: 1,
				},
				Clients: []mgclients.Client{client},
			},
			query:  "name=clientname",
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:    "list users with invalid name",
			token:   validToken,
			groupID: validID,
			query:   "name=invalid",
			status:  http.StatusBadRequest,
			err:     apiutil.ErrValidation,
		},
		{
			desc:    "list users with duplicate name",
			token:   validToken,
			groupID: validID,
			query:   "name=1&name=2",
			status:  http.StatusBadRequest,
			err:     apiutil.ErrInvalidQueryParams,
		},
		{
			desc:    "list users with status",
			token:   validToken,
			groupID: validID,
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Total: 1,
				},
				Clients: []mgclients.Client{client},
			},
			query:  "status=enabled",
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:    "list users with invalid status",
			token:   validToken,
			groupID: validID,
			query:   "status=invalid",
			status:  http.StatusBadRequest,
			err:     apiutil.ErrValidation,
		},
		{
			desc:    "list users with duplicate status",
			token:   validToken,
			groupID: validID,
			query:   "status=enabled&status=disabled",
			status:  http.StatusBadRequest,
			err:     apiutil.ErrInvalidQueryParams,
		},
		{
			desc:    "list users with tags",
			token:   validToken,
			groupID: validID,
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Total: 1,
				},
				Clients: []mgclients.Client{client},
			},
			query:  "tag=tag1,tag2",
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:    "list users with invalid tags",
			token:   validToken,
			groupID: validID,
			query:   "tag=invalid",
			status:  http.StatusBadRequest,
			err:     apiutil.ErrValidation,
		},
		{
			desc:    "list users with duplicate tags",
			token:   validToken,
			groupID: validID,
			query:   "tag=tag1&tag=tag2",
			status:  http.StatusBadRequest,
			err:     apiutil.ErrInvalidQueryParams,
		},
		{
			desc:    "list users with metadata",
			token:   validToken,
			groupID: validID,
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Total: 1,
				},
				Clients: []mgclients.Client{client},
			},
			query:  "metadata=%7B%22domain%22%3A%20%22example.com%22%7D&",
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:    "list users with invalid metadata",
			token:   validToken,
			groupID: validID,
			query:   "metadata=invalid",
			status:  http.StatusBadRequest,
			err:     apiutil.ErrValidation,
		},
		{
			desc:    "list users with duplicate metadata",
			token:   validToken,
			groupID: validID,
			query:   "metadata=%7B%22domain%22%3A%20%22example.com%22%7D&metadata=%7B%22domain%22%3A%20%22example.com%22%7D",
			status:  http.StatusBadRequest,
			err:     apiutil.ErrInvalidQueryParams,
		},
		{
			desc:    "list users with permissions",
			token:   validToken,
			groupID: validID,
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Total: 1,
				},
				Clients: []mgclients.Client{client},
			},
			query:  "permission=view",
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:    "list users with invalid permissions",
			token:   validToken,
			groupID: validID,
			query:   "permission=invalid",
			status:  http.StatusBadRequest,
			err:     apiutil.ErrValidation,
		},
		{
			desc:    "list users with duplicate permissions",
			token:   validToken,
			groupID: validID,
			query:   "permission=view&permission=view",
			status:  http.StatusBadRequest,
			err:     apiutil.ErrInvalidQueryParams,
		},
		{
			desc:    "list users with identity",
			token:   validToken,
			groupID: validID,
			query:   fmt.Sprintf("identity=%s", client.Credentials.Identity),
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Total: 1,
				},
				Clients: []mgclients.Client{
					client,
				},
			},
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:    "list users with invalid identity",
			token:   validToken,
			groupID: validID,
			query:   "identity=invalid",
			status:  http.StatusBadRequest,
			err:     apiutil.ErrValidation,
		},
		{
			desc:    "list users with duplicate identity",
			token:   validToken,
			groupID: validID,
			query:   "identity=1&identity=2",
			status:  http.StatusBadRequest,
			err:     apiutil.ErrInvalidQueryParams,
		},
	}

	for _, tc := range cases {
		req := testRequest{
			client: us.Client(),
			method: http.MethodGet,
			url:    fmt.Sprintf("%s/groups/%s/users?", us.URL, tc.groupID) + tc.query,
			token:  tc.token,
		}

		repoCall := svc.On("ListMembers", mock.Anything, tc.token, mock.Anything, mock.Anything, mock.Anything).Return(
			mgclients.MembersPage{
				Page:    tc.listUsersResponse.Page,
				Members: tc.listUsersResponse.Clients,
			},
			tc.err)
		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
		repoCall.Unset()
	}
}

func TestListUsersByChannelID(t *testing.T) {
	us, svc := newUsersServer()
	defer us.Close()

	cases := []struct {
		desc              string
		token             string
		groupID           string
		page              mgclients.Page
		status            int
		query             string
		listUsersResponse mgclients.ClientsPage
		err               error
	}{
		{
			desc:   "list users with valid token",
			token:  validToken,
			status: http.StatusOK,
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Total: 1,
				},
				Clients: []mgclients.Client{client},
			},
			err: nil,
		},
		{
			desc:   "list users with empty token",
			token:  "",
			status: http.StatusUnauthorized,
			err:    apiutil.ErrBearerToken,
		},
		{
			desc:   "list users with invalid token",
			token:  inValidToken,
			status: http.StatusUnauthorized,
			err:    svcerr.ErrAuthentication,
		},
		{
			desc:  "list users with offset",
			token: validToken,
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Offset: 1,
					Total:  1,
				},
				Clients: []mgclients.Client{client},
			},
			query:  "offset=1",
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:   "list users with invalid offset",
			token:  validToken,
			query:  "offset=invalid",
			status: http.StatusBadRequest,
			err:    apiutil.ErrValidation,
		},
		{
			desc:  "list users with limit",
			token: validToken,
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Limit: 1,
					Total: 1,
				},
				Clients: []mgclients.Client{client},
			},
			query:  "limit=1",
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:   "list users with invalid limit",
			token:  validToken,
			query:  "limit=invalid",
			status: http.StatusBadRequest,
			err:    apiutil.ErrValidation,
		},
		{
			desc:   "list users with limit greater than max",
			token:  validToken,
			query:  fmt.Sprintf("limit=%d", api.MaxLimitSize+1),
			status: http.StatusBadRequest,
			err:    apiutil.ErrValidation,
		},
		{
			desc:  "list users with owner_id",
			token: validToken,
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Total: 1,
				},
				Clients: []mgclients.Client{client},
			},
			query:  fmt.Sprintf("owner_id=%s", validID),
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:   "list users with duplicate owner_id",
			token:  validToken,
			query:  "owner_id=1&owner_id=2",
			status: http.StatusBadRequest,
			err:    apiutil.ErrInvalidQueryParams,
		},
		{
			desc:   "list users with invalid owner_id",
			token:  validToken,
			query:  "owner_id=invalid",
			status: http.StatusBadRequest,
			err:    apiutil.ErrValidation,
		},
		{
			desc:  "list users with name",
			token: validToken,
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Total: 1,
				},
				Clients: []mgclients.Client{client},
			},
			query:  "name=clientname",
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:   "list users with invalid name",
			token:  validToken,
			query:  "name=invalid",
			status: http.StatusBadRequest,
			err:    apiutil.ErrValidation,
		},
		{
			desc:   "list users with duplicate name",
			token:  validToken,
			query:  "name=1&name=2",
			status: http.StatusBadRequest,
			err:    apiutil.ErrInvalidQueryParams,
		},
		{
			desc:  "list users with status",
			token: validToken,
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Total: 1,
				},
				Clients: []mgclients.Client{client},
			},
			query:  "status=enabled",
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:   "list users with invalid status",
			token:  validToken,
			query:  "status=invalid",
			status: http.StatusBadRequest,
			err:    apiutil.ErrValidation,
		},
		{
			desc:   "list users with duplicate status",
			token:  validToken,
			query:  "status=enabled&status=disabled",
			status: http.StatusBadRequest,
			err:    apiutil.ErrInvalidQueryParams,
		},
		{
			desc:  "list users with tags",
			token: validToken,
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Total: 1,
				},
				Clients: []mgclients.Client{client},
			},
			query:  "tag=tag1,tag2",
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:   "list users with invalid tags",
			token:  validToken,
			query:  "tag=invalid",
			status: http.StatusBadRequest,
			err:    apiutil.ErrValidation,
		},
		{
			desc:   "list users with duplicate tags",
			token:  validToken,
			query:  "tag=tag1&tag=tag2",
			status: http.StatusBadRequest,
			err:    apiutil.ErrInvalidQueryParams,
		},
		{
			desc:  "list users with metadata",
			token: validToken,
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Total: 1,
				},
				Clients: []mgclients.Client{client},
			},
			query:  "metadata=%7B%22domain%22%3A%20%22example.com%22%7D&",
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:   "list users with invalid metadata",
			token:  validToken,
			query:  "metadata=invalid",
			status: http.StatusBadRequest,
			err:    apiutil.ErrValidation,
		},
		{
			desc:   "list users with duplicate metadata",
			token:  validToken,
			query:  "metadata=%7B%22domain%22%3A%20%22example.com%22%7D&metadata=%7B%22domain%22%3A%20%22example.com%22%7D",
			status: http.StatusBadRequest,
			err:    apiutil.ErrInvalidQueryParams,
		},
		{
			desc:  "list users with permissions",
			token: validToken,
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Total: 1,
				},
				Clients: []mgclients.Client{client},
			},
			query:  "permission=view",
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:   "list users with invalid permissions",
			token:  validToken,
			query:  "permission=invalid",
			status: http.StatusBadRequest,
			err:    apiutil.ErrValidation,
		},
		{
			desc:   "list users with duplicate permissions",
			token:  validToken,
			query:  "permission=view&permission=view",
			status: http.StatusBadRequest,
			err:    apiutil.ErrInvalidQueryParams,
		},
		{
			desc:  "list users with identity",
			token: validToken,
			query: fmt.Sprintf("identity=%s", client.Credentials.Identity),
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Total: 1,
				},
				Clients: []mgclients.Client{
					client,
				},
			},
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:   "list users with invalid identity",
			token:  validToken,
			query:  "identity=invalid",
			status: http.StatusBadRequest,
			err:    apiutil.ErrValidation,
		},
		{
			desc:   "list users with duplicate identity",
			token:  validToken,
			query:  "identity=1&identity=2",
			status: http.StatusBadRequest,
			err:    apiutil.ErrInvalidQueryParams,
		},
		{
			desc:   "list users with list_perms",
			token:  validToken,
			query:  "list_perms=true",
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:   "list users with invalid list_perms",
			token:  validToken,
			query:  "list_perms=invalid",
			status: http.StatusBadRequest,
			err:    apiutil.ErrValidation,
		},
		{
			desc:   "list users with duplicate list_perms",
			token:  validToken,
			query:  "list_perms=true&list_perms=false",
			status: http.StatusBadRequest,
			err:    apiutil.ErrValidation,
		},
	}

	for _, tc := range cases {
		req := testRequest{
			client: us.Client(),
			method: http.MethodGet,
			url:    fmt.Sprintf("%s/channels/%s/users?", us.URL, validID) + tc.query,
			token:  tc.token,
		}

		repoCall := svc.On("ListMembers", mock.Anything, tc.token, mock.Anything, mock.Anything, mock.Anything).Return(
			mgclients.MembersPage{
				Page:    tc.listUsersResponse.Page,
				Members: tc.listUsersResponse.Clients,
			},
			tc.err)
		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
		repoCall.Unset()
	}
}

func TestListUsersByDomainID(t *testing.T) {
	us, svc := newUsersServer()
	defer us.Close()

	cases := []struct {
		desc              string
		token             string
		groupID           string
		page              mgclients.Page
		status            int
		query             string
		listUsersResponse mgclients.ClientsPage
		err               error
	}{
		{
			desc:   "list users with valid token",
			token:  validToken,
			status: http.StatusOK,
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Total: 1,
				},
				Clients: []mgclients.Client{client},
			},
			err: nil,
		},
		{
			desc:   "list users with empty token",
			token:  "",
			status: http.StatusUnauthorized,
			err:    apiutil.ErrBearerToken,
		},
		{
			desc:   "list users with invalid token",
			token:  inValidToken,
			status: http.StatusUnauthorized,
			err:    svcerr.ErrAuthentication,
		},
		{
			desc:  "list users with offset",
			token: validToken,
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Offset: 1,
					Total:  1,
				},
				Clients: []mgclients.Client{client},
			},
			query:  "offset=1",
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:   "list users with invalid offset",
			token:  validToken,
			query:  "offset=invalid",
			status: http.StatusBadRequest,
			err:    apiutil.ErrValidation,
		},
		{
			desc:  "list users with limit",
			token: validToken,
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Limit: 1,
					Total: 1,
				},
				Clients: []mgclients.Client{client},
			},
			query:  "limit=1",
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:   "list users with invalid limit",
			token:  validToken,
			query:  "limit=invalid",
			status: http.StatusBadRequest,
			err:    apiutil.ErrValidation,
		},
		{
			desc:   "list users with limit greater than max",
			token:  validToken,
			query:  fmt.Sprintf("limit=%d", api.MaxLimitSize+1),
			status: http.StatusBadRequest,
			err:    apiutil.ErrValidation,
		},
		{
			desc:  "list users with owner_id",
			token: validToken,
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Total: 1,
				},
				Clients: []mgclients.Client{client},
			},
			query:  fmt.Sprintf("owner_id=%s", validID),
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:   "list users with duplicate owner_id",
			token:  validToken,
			query:  "owner_id=1&owner_id=2",
			status: http.StatusBadRequest,
			err:    apiutil.ErrInvalidQueryParams,
		},
		{
			desc:   "list users with invalid owner_id",
			token:  validToken,
			query:  "owner_id=invalid",
			status: http.StatusBadRequest,
			err:    apiutil.ErrValidation,
		},
		{
			desc:  "list users with name",
			token: validToken,
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Total: 1,
				},
				Clients: []mgclients.Client{client},
			},
			query:  "name=clientname",
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:   "list users with invalid name",
			token:  validToken,
			query:  "name=invalid",
			status: http.StatusBadRequest,
			err:    apiutil.ErrValidation,
		},
		{
			desc:   "list users with duplicate name",
			token:  validToken,
			query:  "name=1&name=2",
			status: http.StatusBadRequest,
			err:    apiutil.ErrInvalidQueryParams,
		},
		{
			desc:  "list users with status",
			token: validToken,
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Total: 1,
				},
				Clients: []mgclients.Client{client},
			},
			query:  "status=enabled",
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:   "list users with invalid status",
			token:  validToken,
			query:  "status=invalid",
			status: http.StatusBadRequest,
			err:    apiutil.ErrValidation,
		},
		{
			desc:   "list users with duplicate status",
			token:  validToken,
			query:  "status=enabled&status=disabled",
			status: http.StatusBadRequest,
			err:    apiutil.ErrInvalidQueryParams,
		},
		{
			desc:  "list users with tags",
			token: validToken,
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Total: 1,
				},
				Clients: []mgclients.Client{client},
			},
			query:  "tag=tag1,tag2",
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:   "list users with invalid tags",
			token:  validToken,
			query:  "tag=invalid",
			status: http.StatusBadRequest,
			err:    apiutil.ErrValidation,
		},
		{
			desc:   "list users with duplicate tags",
			token:  validToken,
			query:  "tag=tag1&tag=tag2",
			status: http.StatusBadRequest,
			err:    apiutil.ErrInvalidQueryParams,
		},
		{
			desc:  "list users with metadata",
			token: validToken,
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Total: 1,
				},
				Clients: []mgclients.Client{client},
			},
			query:  "metadata=%7B%22domain%22%3A%20%22example.com%22%7D&",
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:   "list users with invalid metadata",
			token:  validToken,
			query:  "metadata=invalid",
			status: http.StatusBadRequest,
			err:    apiutil.ErrValidation,
		},
		{
			desc:   "list users with duplicate metadata",
			token:  validToken,
			query:  "metadata=%7B%22domain%22%3A%20%22example.com%22%7D&metadata=%7B%22domain%22%3A%20%22example.com%22%7D",
			status: http.StatusBadRequest,
			err:    apiutil.ErrInvalidQueryParams,
		},
		{
			desc:  "list users with permissions",
			token: validToken,
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Total: 1,
				},
				Clients: []mgclients.Client{client},
			},
			query:  "permission=membership",
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:   "list users with invalid permissions",
			token:  validToken,
			query:  "permission=invalid",
			status: http.StatusBadRequest,
			err:    apiutil.ErrValidation,
		},
		{
			desc:   "list users with duplicate permissions",
			token:  validToken,
			query:  "permission=view&permission=view",
			status: http.StatusBadRequest,
			err:    apiutil.ErrInvalidQueryParams,
		},
		{
			desc:  "list users with identity",
			token: validToken,
			query: fmt.Sprintf("identity=%s", client.Credentials.Identity),
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Total: 1,
				},
				Clients: []mgclients.Client{
					client,
				},
			},
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:   "list users with invalid identity",
			token:  validToken,
			query:  "identity=invalid",
			status: http.StatusBadRequest,
			err:    apiutil.ErrValidation,
		},
		{
			desc:   "list users with duplicate identity",
			token:  validToken,
			query:  "identity=1&identity=2",
			status: http.StatusBadRequest,
			err:    apiutil.ErrInvalidQueryParams,
		},
		{
			desc:  "list users wiith list permissions",
			token: validToken,
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Total: 1,
				},
				Clients: []mgclients.Client{
					client,
				},
			},
			query:  "list_perms=true",
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:   "list users with invalid list_perms",
			token:  validToken,
			query:  "list_perms=invalid",
			status: http.StatusBadRequest,
			err:    apiutil.ErrValidation,
		},
		{
			desc:   "list users with duplicate list_perms",
			token:  validToken,
			query:  "list_perms=true&list_perms=false",
			status: http.StatusBadRequest,
			err:    apiutil.ErrValidation,
		},
	}

	for _, tc := range cases {
		req := testRequest{
			client: us.Client(),
			method: http.MethodGet,
			url:    fmt.Sprintf("%s/domains/%s/users?", us.URL, validID) + tc.query,
			token:  tc.token,
		}

		repoCall := svc.On("ListMembers", mock.Anything, tc.token, mock.Anything, mock.Anything, mock.Anything).Return(
			mgclients.MembersPage{
				Page:    tc.listUsersResponse.Page,
				Members: tc.listUsersResponse.Clients,
			},
			tc.err)
		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		assert.Equal(t, tc.status, res.StatusCode)
		repoCall.Unset()
	}
}

func TestListUsersByThingID(t *testing.T) {
	us, svc := newUsersServer()
	defer us.Close()

	cases := []struct {
		desc              string
		token             string
		groupID           string
		page              mgclients.Page
		status            int
		query             string
		listUsersResponse mgclients.ClientsPage
		err               error
	}{
		{
			desc:   "list users with valid token",
			token:  validToken,
			status: http.StatusOK,
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Total: 1,
				},
				Clients: []mgclients.Client{client},
			},
			err: nil,
		},
		{
			desc:   "list users with empty token",
			token:  "",
			status: http.StatusUnauthorized,
			err:    apiutil.ErrBearerToken,
		},
		{
			desc:   "list users with invalid token",
			token:  inValidToken,
			status: http.StatusUnauthorized,
			err:    svcerr.ErrAuthentication,
		},
		{
			desc:  "list users with offset",
			token: validToken,
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Offset: 1,
					Total:  1,
				},
				Clients: []mgclients.Client{client},
			},
			query:  "offset=1",
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:   "list users with invalid offset",
			token:  validToken,
			query:  "offset=invalid",
			status: http.StatusBadRequest,
			err:    apiutil.ErrValidation,
		},
		{
			desc:  "list users with limit",
			token: validToken,
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Limit: 1,
					Total: 1,
				},
				Clients: []mgclients.Client{client},
			},
			query:  "limit=1",
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:   "list users with invalid limit",
			token:  validToken,
			query:  "limit=invalid",
			status: http.StatusBadRequest,
			err:    apiutil.ErrValidation,
		},
		{
			desc:   "list users with limit greater than max",
			token:  validToken,
			query:  fmt.Sprintf("limit=%d", api.MaxLimitSize+1),
			status: http.StatusBadRequest,
			err:    apiutil.ErrValidation,
		},
		{
			desc:  "list users with owner_id",
			token: validToken,
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Total: 1,
				},
				Clients: []mgclients.Client{client},
			},
			query:  fmt.Sprintf("owner_id=%s", validID),
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:   "list users with duplicate owner_id",
			token:  validToken,
			query:  "owner_id=1&owner_id=2",
			status: http.StatusBadRequest,
			err:    apiutil.ErrInvalidQueryParams,
		},
		{
			desc:   "list users with invalid owner_id",
			token:  validToken,
			query:  "owner_id=invalid",
			status: http.StatusBadRequest,
			err:    apiutil.ErrValidation,
		},
		{
			desc:  "list users with name",
			token: validToken,
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Total: 1,
				},
				Clients: []mgclients.Client{client},
			},
			query:  "name=clientname",
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:   "list users with invalid name",
			token:  validToken,
			query:  "name=invalid",
			status: http.StatusBadRequest,
			err:    apiutil.ErrValidation,
		},
		{
			desc:   "list users with duplicate name",
			token:  validToken,
			query:  "name=1&name=2",
			status: http.StatusBadRequest,
			err:    apiutil.ErrInvalidQueryParams,
		},
		{
			desc:  "list users with status",
			token: validToken,
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Total: 1,
				},
				Clients: []mgclients.Client{client},
			},
			query:  "status=enabled",
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:   "list users with invalid status",
			token:  validToken,
			query:  "status=invalid",
			status: http.StatusBadRequest,
			err:    apiutil.ErrValidation,
		},
		{
			desc:   "list users with duplicate status",
			token:  validToken,
			query:  "status=enabled&status=disabled",
			status: http.StatusBadRequest,
			err:    apiutil.ErrInvalidQueryParams,
		},
		{
			desc:  "list users with tags",
			token: validToken,
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Total: 1,
				},
				Clients: []mgclients.Client{client},
			},
			query:  "tag=tag1,tag2",
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:   "list users with invalid tags",
			token:  validToken,
			query:  "tag=invalid",
			status: http.StatusBadRequest,
			err:    apiutil.ErrValidation,
		},
		{
			desc:   "list users with duplicate tags",
			token:  validToken,
			query:  "tag=tag1&tag=tag2",
			status: http.StatusBadRequest,
			err:    apiutil.ErrInvalidQueryParams,
		},
		{
			desc:  "list users with metadata",
			token: validToken,
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Total: 1,
				},
				Clients: []mgclients.Client{client},
			},
			query:  "metadata=%7B%22domain%22%3A%20%22example.com%22%7D&",
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:   "list users with invalid metadata",
			token:  validToken,
			query:  "metadata=invalid",
			status: http.StatusBadRequest,
			err:    apiutil.ErrValidation,
		},
		{
			desc:   "list users with duplicate metadata",
			token:  validToken,
			query:  "metadata=%7B%22domain%22%3A%20%22example.com%22%7D&metadata=%7B%22domain%22%3A%20%22example.com%22%7D",
			status: http.StatusBadRequest,
			err:    apiutil.ErrInvalidQueryParams,
		},
		{
			desc:  "list users with permissions",
			token: validToken,
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Total: 1,
				},
				Clients: []mgclients.Client{client},
			},
			query:  "permission=view",
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:   "list users with invalid permissions",
			token:  validToken,
			query:  "permission=invalid",
			status: http.StatusBadRequest,
			err:    apiutil.ErrValidation,
		},
		{
			desc:   "list users with duplicate permissions",
			token:  validToken,
			query:  "permission=view&permission=view",
			status: http.StatusBadRequest,
			err:    apiutil.ErrInvalidQueryParams,
		},
		{
			desc:  "list users with identity",
			token: validToken,
			query: fmt.Sprintf("identity=%s", client.Credentials.Identity),
			listUsersResponse: mgclients.ClientsPage{
				Page: mgclients.Page{
					Total: 1,
				},
				Clients: []mgclients.Client{
					client,
				},
			},
			status: http.StatusOK,
			err:    nil,
		},
		{
			desc:   "list users with invalid identity",
			token:  validToken,
			query:  "identity=invalid",
			status: http.StatusBadRequest,
			err:    apiutil.ErrValidation,
		},
		{
			desc:   "list users with duplicate identity",
			token:  validToken,
			query:  "identity=1&identity=2",
			status: http.StatusBadRequest,
			err:    apiutil.ErrInvalidQueryParams,
		},
	}

	for _, tc := range cases {
		req := testRequest{
			client: us.Client(),
			method: http.MethodGet,
			url:    fmt.Sprintf("%s/things/%s/users?", us.URL, validID) + tc.query,
			token:  tc.token,
		}

		repoCall := svc.On("ListMembers", mock.Anything, tc.token, mock.Anything, mock.Anything, mock.Anything).Return(
			mgclients.MembersPage{
				Page:    tc.listUsersResponse.Page,
				Members: tc.listUsersResponse.Clients,
			},
			tc.err)
		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		assert.Equal(t, tc.status, res.StatusCode)
		repoCall.Unset()
	}
}

type respBody struct {
	Err     string           `json:"error"`
	Message string           `json:"message"`
	Total   int              `json:"total"`
	ID      string           `json:"id"`
	Tags    []string         `json:"tags"`
	Role    mgclients.Role   `json:"role"`
	Status  mgclients.Status `json:"status"`
}
