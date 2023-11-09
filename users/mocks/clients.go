// Copyright (c) Magistrala
// SPDX-License-Identifier: Apache-2.0

package mocks

import (
	"context"

	mgclients "github.com/absmach/magistrala/pkg/clients"
	repoerror "github.com/absmach/magistrala/pkg/errors/repository"
	"github.com/absmach/magistrala/users/postgres"
	"github.com/stretchr/testify/mock"
)

const WrongID = "wrongID"

var _ postgres.Repository = (*Repository)(nil)

type Repository struct {
	mock.Mock
}

func (m *Repository) ChangeStatus(ctx context.Context, client mgclients.Client) (mgclients.Client, error) {
	ret := m.Called(ctx, client)

	if client.ID == WrongID {
		return mgclients.Client{}, repoerror.ErrNotFound
	}

	if client.Status != mgclients.EnabledStatus && client.Status != mgclients.DisabledStatus {
		return mgclients.Client{}, repoerror.ErrMalformedEntity
	}

	return ret.Get(0).(mgclients.Client), ret.Error(1)
}

func (m *Repository) Members(ctx context.Context, groupID string, pm mgclients.Page) (mgclients.MembersPage, error) {
	ret := m.Called(ctx, groupID, pm)
	if groupID == WrongID {
		return mgclients.MembersPage{}, repoerror.ErrNotFound
	}

	return ret.Get(0).(mgclients.MembersPage), ret.Error(1)
}

func (m *Repository) RetrieveAll(ctx context.Context, pm mgclients.Page) (mgclients.ClientsPage, error) {
	ret := m.Called(ctx, pm)

	return ret.Get(0).(mgclients.ClientsPage), ret.Error(1)
}

func (m *Repository) RetrieveByID(ctx context.Context, id string) (mgclients.Client, error) {
	ret := m.Called(ctx, id)

	if id == WrongID {
		return mgclients.Client{}, repoerror.ErrNotFound
	}

	return ret.Get(0).(mgclients.Client), ret.Error(1)
}

func (m *Repository) RetrieveByIdentity(ctx context.Context, identity string) (mgclients.Client, error) {
	ret := m.Called(ctx, identity)

	if identity == "" {
		return mgclients.Client{}, repoerror.ErrMalformedEntity
	}

	return ret.Get(0).(mgclients.Client), ret.Error(1)
}

func (m *Repository) Save(ctx context.Context, client mgclients.Client) (mgclients.Client, error) {
	ret := m.Called(ctx, client)
	if client.Owner == WrongID {
		return mgclients.Client{}, repoerror.ErrMalformedEntity
	}
	if client.Credentials.Secret == "" {
		return mgclients.Client{}, repoerror.ErrMalformedEntity
	}

	return client, ret.Error(1)
}

func (m *Repository) Update(ctx context.Context, client mgclients.Client) (mgclients.Client, error) {
	ret := m.Called(ctx, client)

	if client.ID == WrongID {
		return mgclients.Client{}, repoerror.ErrNotFound
	}
	return ret.Get(0).(mgclients.Client), ret.Error(1)
}

func (m *Repository) UpdateIdentity(ctx context.Context, client mgclients.Client) (mgclients.Client, error) {
	ret := m.Called(ctx, client)

	if client.ID == WrongID {
		return mgclients.Client{}, repoerror.ErrNotFound
	}
	if client.Credentials.Identity == "" {
		return mgclients.Client{}, repoerror.ErrMalformedEntity
	}

	return ret.Get(0).(mgclients.Client), ret.Error(1)
}

func (m *Repository) UpdateSecret(ctx context.Context, client mgclients.Client) (mgclients.Client, error) {
	ret := m.Called(ctx, client)

	if client.ID == WrongID {
		return mgclients.Client{}, repoerror.ErrNotFound
	}
	if client.Credentials.Secret == "" {
		return mgclients.Client{}, repoerror.ErrMalformedEntity
	}

	return ret.Get(0).(mgclients.Client), ret.Error(1)
}

func (m *Repository) UpdateTags(ctx context.Context, client mgclients.Client) (mgclients.Client, error) {
	ret := m.Called(ctx, client)

	if client.ID == WrongID {
		return mgclients.Client{}, repoerror.ErrNotFound
	}

	return ret.Get(0).(mgclients.Client), ret.Error(1)
}

func (m *Repository) UpdateOwner(ctx context.Context, client mgclients.Client) (mgclients.Client, error) {
	ret := m.Called(ctx, client)

	if client.ID == WrongID {
		return mgclients.Client{}, repoerror.ErrNotFound
	}

	return ret.Get(0).(mgclients.Client), ret.Error(1)
}

func (m *Repository) RetrieveBySecret(ctx context.Context, key string) (mgclients.Client, error) {
	ret := m.Called(ctx, key)

	if key == "" {
		return mgclients.Client{}, repoerror.ErrMalformedEntity
	}

	return ret.Get(0).(mgclients.Client), ret.Error(1)
}

func (m *Repository) IsOwner(ctx context.Context, clientID string, ownerID string) error {
	ret := m.Called(ctx, clientID, ownerID)

	if clientID == WrongID {
		return repoerror.ErrNotFound
	}

	return ret.Error(0)
}

func (m *Repository) RetrieveAllByIDs(ctx context.Context, pm mgclients.Page) (mgclients.ClientsPage, error) {
	ret := m.Called(ctx, pm)

	return ret.Get(0).(mgclients.ClientsPage), ret.Error(1)
}
