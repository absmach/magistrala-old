// Copyright (c) Magistrala
// SPDX-License-Identifier: Apache-2.0

package mocks

import (
	context "context"

	auth "github.com/absmach/magistrala/auth"
	"github.com/stretchr/testify/mock"
)

var _ auth.DomainsRepository = (*DomainsRepo)(nil)

type DomainsRepo struct {
	mock.Mock
}

func (m *DomainsRepo) Save(ctx context.Context, d auth.Domain) (auth.Domain, error) {
	ret := m.Called(ctx, d)

	return ret.Get(0).(auth.Domain), ret.Error(1)
}
func (m *DomainsRepo) RetrieveByID(ctx context.Context, id string) (auth.Domain, error) {
	ret := m.Called(ctx, id)
	return ret.Get(0).(auth.Domain), ret.Error(1)
}
func (m *DomainsRepo) RetrieveAllByIDs(ctx context.Context, pm auth.Page) (auth.DomainsPage, error) {
	ret := m.Called(ctx, pm)

	return ret.Get(0).(auth.DomainsPage), ret.Error(1)
}
func (m *DomainsRepo) Update(ctx context.Context, d auth.DomainReq) (auth.Domain, error) {
	ret := m.Called(ctx, d)

	return ret.Get(0).(auth.Domain), ret.Error(1)
}
func (m *DomainsRepo) Delete(ctx context.Context, id string) error {
	ret := m.Called(ctx, id)

	return ret.Error(0)
}
func (m *DomainsRepo) SavePolicyCopy(ctx context.Context, pc auth.PolicyCopy) error {
	ret := m.Called(ctx, pc)

	return ret.Error(0)
}
func (m *DomainsRepo) DeletePolicyCopy(ctx context.Context, pc auth.PolicyCopy) error {
	ret := m.Called(ctx, pc)

	return ret.Error(0)
}
