// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package mocks

import (
	"context"

	"github.com/absmach/magistrala/auth"
	"github.com/stretchr/testify/mock"
)

type AuthService struct {
	mock.Mock
}

func (svc *AuthService) Issue(ctx context.Context, token string, key auth.Key) (auth.Token, error) {
	ret := svc.Called(ctx, token, key)
	return ret.Get(0).(auth.Token), ret.Error(1)
}

func (svc *AuthService) Revoke(ctx context.Context, token, id string) error {
	ret := svc.Called(ctx, token, id)
	return ret.Error(0)
}

func (svc *AuthService) RetrieveKey(ctx context.Context, token, id string) (auth.Key, error) {
	ret := svc.Called(ctx, token, id)
	return ret.Get(0).(auth.Key), ret.Error(1)
}

func (svc *AuthService) Identify(ctx context.Context, token string) (auth.Key, error) {
	ret := svc.Called(ctx, token)
	return ret.Get(0).(auth.Key), ret.Error(1)
}

func (svc *AuthService) Authorize(ctx context.Context, pr auth.PolicyReq) error {
	ret := svc.Called(ctx, pr)
	return ret.Error(0)
}

func (svc *AuthService) AddPolicy(ctx context.Context, pr auth.PolicyReq) error {
	ret := svc.Called(ctx, pr)
	return ret.Error(0)
}

func (svc *AuthService) PolicyValidation(pr auth.PolicyReq) error {
	ret := svc.Called(pr)
	return ret.Error(0)
}

func (svc *AuthService) AddPolicies(ctx context.Context, prs []auth.PolicyReq) error {
	ret := svc.Called(ctx, prs)
	return ret.Error(0)
}

func (svc *AuthService) DeletePolicy(ctx context.Context, pr auth.PolicyReq) error {
	ret := svc.Called(ctx, pr)
	return ret.Error(0)
}

func (svc *AuthService) DeletePolicies(ctx context.Context, prs []auth.PolicyReq) error {
	ret := svc.Called(ctx, prs)
	return ret.Error(0)
}

func (svc *AuthService) ListObjects(ctx context.Context, pr auth.PolicyReq, nextPageToken string, limit int32) (auth.PolicyPage, error) {
	ret := svc.Called(ctx, pr, nextPageToken, limit)
	return ret.Get(0).(auth.PolicyPage), ret.Error(1)
}

func (svc *AuthService) ListAllObjects(ctx context.Context, pr auth.PolicyReq) (auth.PolicyPage, error) {
	ret := svc.Called(ctx, pr)
	return ret.Get(0).(auth.PolicyPage), ret.Error(1)
}

func (svc *AuthService) CountObjects(ctx context.Context, pr auth.PolicyReq) (int, error) {
	ret := svc.Called(ctx, pr)
	return ret.Get(0).(int), ret.Error(1)
}

func (svc *AuthService) ListSubjects(ctx context.Context, pr auth.PolicyReq, nextPageToken string, limit int32) (auth.PolicyPage, error) {
	ret := svc.Called(ctx, pr, nextPageToken, limit)
	return ret.Get(0).(auth.PolicyPage), ret.Error(1)
}

func (svc *AuthService) ListAllSubjects(ctx context.Context, pr auth.PolicyReq) (auth.PolicyPage, error) {
	ret := svc.Called(ctx, pr)
	return ret.Get(0).(auth.PolicyPage), ret.Error(1)
}

func (svc *AuthService) CountSubjects(ctx context.Context, pr auth.PolicyReq) (int, error) {
	ret := svc.Called(ctx, pr)
	return ret.Get(0).(int), ret.Error(1)
}

func (svc *AuthService) ListPermissions(ctx context.Context, pr auth.PolicyReq, filterPermisions []string) (auth.Permissions, error) {
	ret := svc.Called(ctx, pr, filterPermisions)
	return ret.Get(0).(auth.Permissions), ret.Error(1)
}

func (svc *AuthService) CreateDomain(ctx context.Context, token string, d auth.Domain) (do auth.Domain, err error) {
	ret := svc.Called(ctx, token, d)
	return ret.Get(0).(auth.Domain), ret.Error(1)
}

func (svc *AuthService) RetrieveDomain(ctx context.Context, token, id string) (auth.Domain, error) {
	ret := svc.Called(ctx, token, id)
	return ret.Get(0).(auth.Domain), ret.Error(1)
}

func (svc *AuthService) RetrieveDomainPermissions(ctx context.Context, token, id string) (auth.Permissions, error) {
	ret := svc.Called(ctx, token, id)
	return ret.Get(0).(auth.Permissions), ret.Error(1)
}

func (svc *AuthService) UpdateDomain(ctx context.Context, token, id string, d auth.DomainReq) (auth.Domain, error) {
	ret := svc.Called(ctx, token, id, d)
	return ret.Get(0).(auth.Domain), ret.Error(1)
}

func (svc *AuthService) ChangeDomainStatus(ctx context.Context, token, id string, d auth.DomainReq) (auth.Domain, error) {
	ret := svc.Called(ctx, token, id, d)
	return ret.Get(0).(auth.Domain), ret.Error(1)
}

func (svc *AuthService) ListDomains(ctx context.Context, token string, page auth.Page) (auth.DomainsPage, error) {
	ret := svc.Called(ctx, token, page)
	return ret.Get(0).(auth.DomainsPage), ret.Error(1)
}

func (svc *AuthService) AssignUsers(ctx context.Context, token, id string, userIds []string, relation string) error {
	ret := svc.Called(ctx, token, id, userIds, relation)
	return ret.Error(0)
}

func (svc *AuthService) UnassignUsers(ctx context.Context, token, id string, userIds []string, relation string) error {
	ret := svc.Called(ctx, token, id, userIds, relation)
	return ret.Error(0)
}

func (svc *AuthService) ListUserDomains(ctx context.Context, token, userID string, page auth.Page) (auth.DomainsPage, error) {
	ret := svc.Called(ctx, token, userID, page)
	return ret.Get(0).(auth.DomainsPage), ret.Error(1)
}
