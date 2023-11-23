// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package certs_test

import (
	"context"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/absmach/magistrala"
	authmocks "github.com/absmach/magistrala/auth/mocks"
	"github.com/absmach/magistrala/certs"
	"github.com/absmach/magistrala/certs/mocks"
	chmocks "github.com/absmach/magistrala/internal/groups/mocks"
	"github.com/absmach/magistrala/logger"
	"github.com/absmach/magistrala/pkg/clients"
	"github.com/absmach/magistrala/pkg/errors"
	repoerror "github.com/absmach/magistrala/pkg/errors/repository"
	svcerror "github.com/absmach/magistrala/pkg/errors/service"
	mgsdk "github.com/absmach/magistrala/pkg/sdk/go"
	"github.com/absmach/magistrala/pkg/uuid"
	"github.com/absmach/magistrala/things"
	httpapi "github.com/absmach/magistrala/things/api/http"
	thmocks "github.com/absmach/magistrala/things/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	invalid   = "invalid"
	email     = "user@example.com"
	token     = "token"
	thingsNum = 1
	thingKey  = "thingKey"
	thingID   = "1"
	ttl       = "1h"
	certNum   = 10
	validID   = "d4ebb847-5d0e-4e46-bdd9-b6aceaaa3a22"

	cfgAuthTimeout = "1s"

	caPath            = "../docker/ssl/certs/ca.crt"
	caKeyPath         = "../docker/ssl/certs/ca.key"
	cfgSignHoursValid = "24h"
	instanceID        = "5de9b29a-feb9-11ed-be56-0242ac120002"
)

func newThingsServer(svc things.Service) *httptest.Server {
	logger := logger.NewMock()
	mux := chi.NewMux()
	httpapi.MakeHandler(svc, nil, mux, logger, instanceID)
	return httptest.NewServer(mux)
}

func newService(t *testing.T) (certs.Service, *authmocks.Service, *thmocks.Repository) {
	auth := new(authmocks.Service)
	tsvc, trepo := newThingsService(auth)
	server := newThingsServer(tsvc)

	config := mgsdk.Config{
		ThingsURL: server.URL,
	}

	sdk := mgsdk.NewSDK(config)
	repo := mocks.NewCertsRepository()

	tlsCert, caCert, err := certs.LoadCertificates(caPath, caKeyPath)
	require.Nil(t, err, fmt.Sprintf("unexpected cert loading error: %s\n", err))

	authTimeout, err := time.ParseDuration(cfgAuthTimeout)
	require.Nil(t, err, fmt.Sprintf("unexpected auth timeout parsing error: %s\n", err))

	pki := mocks.NewPkiAgent(tlsCert, caCert, cfgSignHoursValid, authTimeout)

	return certs.New(auth, repo, sdk, pki), auth, trepo
}

func newThingsService(auth *authmocks.Service) (things.Service, *thmocks.Repository) {
	thingCache := thmocks.NewCache()
	idProvider := uuid.NewMock()
	cRepo := new(thmocks.Repository)
	gRepo := new(chmocks.Repository)

	return things.NewService(auth, cRepo, gRepo, thingCache, idProvider), cRepo
}

func TestIssueCert(t *testing.T) {
	svc, auth, trepo := newService(t)

	cases := []struct {
		token   string
		desc    string
		thingID string
		ttl     string
		key     string
		err     error
	}{
		{
			desc:    "issue new cert",
			token:   token,
			thingID: thingID,
			ttl:     ttl,
			err:     nil,
		},
		{
			desc:    "issue new cert for non existing thing id",
			token:   token,
			thingID: "2",
			ttl:     ttl,
			err:     certs.ErrFailedCertCreation,
		},
		{
			desc:    "issue new cert for non existing thing id",
			token:   invalid,
			thingID: thingID,
			ttl:     ttl,
			err:     errors.ErrAuthentication,
		},
	}

	for _, tc := range cases {
		repoCall := auth.On("Identify", mock.Anything, &magistrala.IdentityReq{Token: tc.token}).Return(&magistrala.IdentityRes{Id: validID}, nil)
		repoCall1 := auth.On("Authorize", mock.Anything, mock.Anything).Return(&magistrala.AuthorizeRes{Authorized: true}, tc.err)
		repoCall2 := trepo.On("RetrieveByID", mock.Anything, mock.Anything).Return(clients.Client{ID: tc.thingID, Credentials: clients.Credentials{Secret: thingKey}}, tc.err)
		c, err := svc.IssueCert(context.Background(), tc.token, tc.thingID, tc.ttl)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
		cert, _ := certs.ReadCert([]byte(c.ClientCert))
		if cert != nil {
			assert.True(t, strings.Contains(cert.Subject.CommonName, thingKey), fmt.Sprintf("%s: expected %s got %s\n", tc.desc, thingKey, cert.Subject.CommonName))
		}
		repoCall.Unset()
		repoCall1.Unset()
		repoCall2.Unset()
	}
}

func TestRevokeCert(t *testing.T) {
	svc, auth, trepo := newService(t)

	repoCall := auth.On("Identify", mock.Anything, &magistrala.IdentityReq{Token: token}).Return(&magistrala.IdentityRes{Id: validID}, nil)
	repoCall1 := auth.On("Authorize", mock.Anything, mock.Anything).Return(&magistrala.AuthorizeRes{Authorized: true}, nil)
	repoCall2 := trepo.On("RetrieveByID", mock.Anything, mock.Anything).Return(clients.Client{ID: thingID, Credentials: clients.Credentials{Secret: thingKey}}, nil)
	_, err := svc.IssueCert(context.Background(), token, thingID, ttl)
	require.Nil(t, err, fmt.Sprintf("unexpected service creation error: %s\n", err))
	repoCall.Unset()
	repoCall1.Unset()
	repoCall2.Unset()

	cases := []struct {
		token   string
		desc    string
		thingID string
		err     error
	}{
		{
			desc:    "revoke cert",
			token:   token,
			thingID: thingID,
			err:     nil,
		},
		{
			desc:    "revoke cert for invalid token",
			token:   invalid,
			thingID: thingID,
			err:     errors.ErrAuthentication,
		},
		{
			desc:    "revoke cert for invalid thing id",
			token:   token,
			thingID: "2",
			err:     certs.ErrFailedCertRevocation,
		},
	}

	for _, tc := range cases {
		repoCall := auth.On("Identify", mock.Anything, &magistrala.IdentityReq{Token: tc.token}).Return(&magistrala.IdentityRes{Id: validID}, nil)
		repoCall1 := auth.On("Authorize", mock.Anything, mock.Anything).Return(&magistrala.AuthorizeRes{Authorized: true}, tc.err)
		repoCall2 := trepo.On("RetrieveByID", mock.Anything, mock.Anything).Return(clients.Client{ID: tc.thingID, Credentials: clients.Credentials{Secret: thingKey}}, tc.err)
		_, err := svc.RevokeCert(context.Background(), tc.token, tc.thingID)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
		repoCall.Unset()
		repoCall1.Unset()
		repoCall2.Unset()
	}
}

func TestListCerts(t *testing.T) {
	svc, auth, trepo := newService(t)

	for i := 0; i < certNum; i++ {
		repoCall := auth.On("Identify", mock.Anything, &magistrala.IdentityReq{Token: token}).Return(&magistrala.IdentityRes{Id: validID}, nil)
		repoCall1 := auth.On("Authorize", mock.Anything, mock.Anything).Return(&magistrala.AuthorizeRes{Authorized: true}, nil)
		repoCall2 := trepo.On("RetrieveByID", mock.Anything, mock.Anything).Return(clients.Client{ID: thingID, Credentials: clients.Credentials{Secret: thingKey}}, nil)
		_, err := svc.IssueCert(context.Background(), token, thingID, ttl)
		require.Nil(t, err, fmt.Sprintf("unexpected cert creation error: %s\n", err))
		repoCall.Unset()
		repoCall1.Unset()
		repoCall2.Unset()
	}

	cases := []struct {
		token   string
		desc    string
		thingID string
		offset  uint64
		limit   uint64
		size    uint64
		err     error
	}{
		{
			desc:    "list all certs with valid token",
			token:   token,
			thingID: thingID,
			offset:  0,
			limit:   certNum,
			size:    certNum,
			err:     nil,
		},
		{
			desc:    "list all certs with invalid token",
			token:   invalid,
			thingID: thingID,
			offset:  0,
			limit:   certNum,
			size:    0,
			err:     errors.ErrAuthentication,
		},
		{
			desc:    "list half certs with valid token",
			token:   token,
			thingID: thingID,
			offset:  certNum / 2,
			limit:   certNum,
			size:    certNum / 2,
			err:     nil,
		},
		{
			desc:    "list last cert with valid token",
			token:   token,
			thingID: thingID,
			offset:  certNum - 1,
			limit:   certNum,
			size:    1,
			err:     nil,
		},
	}

	for _, tc := range cases {
		repoCall := auth.On("Identify", mock.Anything, &magistrala.IdentityReq{Token: tc.token}).Return(&magistrala.IdentityRes{Id: validID}, nil)
		page, err := svc.ListCerts(context.Background(), tc.token, tc.thingID, tc.offset, tc.limit)
		size := uint64(len(page.Certs))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", tc.desc, tc.size, size))
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
		repoCall.Unset()
	}
}

func TestListSerials(t *testing.T) {
	svc, auth, trepo := newService(t)

	var issuedCerts []certs.Cert
	for i := 0; i < certNum; i++ {
		repoCall := auth.On("Identify", mock.Anything, &magistrala.IdentityReq{Token: token}).Return(&magistrala.IdentityRes{Id: validID}, nil)
		repoCall1 := auth.On("Authorize", mock.Anything, mock.Anything).Return(&magistrala.AuthorizeRes{Authorized: true}, nil)
		repoCall2 := trepo.On("RetrieveByID", mock.Anything, mock.Anything).Return(clients.Client{ID: thingID, Credentials: clients.Credentials{Secret: thingKey}}, nil)
		cert, err := svc.IssueCert(context.Background(), token, thingID, ttl)
		assert.Nil(t, err, fmt.Sprintf("unexpected cert creation error: %s\n", err))
		repoCall.Unset()
		repoCall1.Unset()
		repoCall2.Unset()

		crt := certs.Cert{
			OwnerID: cert.OwnerID,
			ThingID: cert.ThingID,
			Serial:  cert.Serial,
			Expire:  cert.Expire,
		}
		issuedCerts = append(issuedCerts, crt)
	}

	cases := []struct {
		token   string
		desc    string
		thingID string
		offset  uint64
		limit   uint64
		certs   []certs.Cert
		err     error
	}{
		{
			desc:    "list all certs with valid token",
			token:   token,
			thingID: thingID,
			offset:  0,
			limit:   certNum,
			certs:   issuedCerts,
			err:     nil,
		},
		{
			desc:    "list all certs with invalid token",
			token:   invalid,
			thingID: thingID,
			offset:  0,
			limit:   certNum,
			certs:   nil,
			err:     errors.ErrAuthentication,
		},
		{
			desc:    "list half certs with valid token",
			token:   token,
			thingID: thingID,
			offset:  certNum / 2,
			limit:   certNum,
			certs:   issuedCerts[certNum/2:],
			err:     nil,
		},
		{
			desc:    "list last cert with valid token",
			token:   token,
			thingID: thingID,
			offset:  certNum - 1,
			limit:   certNum,
			certs:   []certs.Cert{issuedCerts[certNum-1]},
			err:     nil,
		},
	}

	for _, tc := range cases {
		repoCall := auth.On("Identify", mock.Anything, &magistrala.IdentityReq{Token: tc.token}).Return(&magistrala.IdentityRes{Id: validID}, nil)
		page, err := svc.ListSerials(context.Background(), tc.token, tc.thingID, tc.offset, tc.limit)
		assert.Equal(t, tc.certs, page.Certs, fmt.Sprintf("%s: expected %v got %v\n", tc.desc, tc.certs, page.Certs))
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
		repoCall.Unset()
	}
}

func TestViewCert(t *testing.T) {
	svc, auth, trepo := newService(t)

	repoCall := auth.On("Identify", mock.Anything, &magistrala.IdentityReq{Token: token}).Return(&magistrala.IdentityRes{Id: validID}, nil)
	repoCall1 := auth.On("Authorize", mock.Anything, mock.Anything).Return(&magistrala.AuthorizeRes{Authorized: true}, nil)
	repoCall2 := trepo.On("RetrieveByID", mock.Anything, mock.Anything).Return(clients.Client{ID: thingID, Credentials: clients.Credentials{Secret: thingKey}}, nil)
	ic, err := svc.IssueCert(context.Background(), token, thingID, ttl)
	require.Nil(t, err, fmt.Sprintf("unexpected cert creation error: %s\n", err))
	repoCall.Unset()
	repoCall1.Unset()
	repoCall2.Unset()

	cert := certs.Cert{
		ThingID:    thingID,
		ClientCert: ic.ClientCert,
		Serial:     ic.Serial,
		Expire:     ic.Expire,
	}

	cases := []struct {
		token    string
		desc     string
		serialID string
		cert     certs.Cert
		err      error
	}{
		{
			desc:     "list cert with valid token and serial",
			token:    token,
			serialID: cert.Serial,
			cert:     cert,
			err:      nil,
		},
		{
			desc:     "list cert with invalid token",
			token:    invalid,
			serialID: cert.Serial,
			cert:     certs.Cert{},
			err:      errors.ErrAuthentication,
		},
		{
			desc:     "list cert with invalid serial",
			token:    token,
			serialID: invalid,
			cert:     certs.Cert{},
			err:      errors.ErrNotFound,
		},
	}

	for _, tc := range cases {
		repoCall := auth.On("Identify", mock.Anything, &magistrala.IdentityReq{Token: tc.token}).Return(&magistrala.IdentityRes{Id: validID}, nil)
		cert, err := svc.ViewCert(context.Background(), tc.token, tc.serialID)
		assert.Equal(t, tc.cert, cert, fmt.Sprintf("%s: expected %v got %v\n", tc.desc, tc.cert, cert))
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
		repoCall.Unset()
	}
}
