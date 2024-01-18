// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"strings"
	"testing"

	"github.com/absmach/magistrala/internal/api"
	"github.com/absmach/magistrala/internal/apiutil"
	"github.com/absmach/magistrala/internal/testsutil"
	mgclients "github.com/absmach/magistrala/pkg/clients"
	"github.com/stretchr/testify/assert"
)

const (
	valid   = "valid"
	invalid = "invalid"
)

var validID = testsutil.GenerateUUID(&testing.T{})

func TestCreateClientReqValidate(t *testing.T) {
	cases := []struct {
		desc string
		req  createClientReq
		err  error
	}{
		{
			desc: "valid request",
			req: createClientReq{
				token: valid,
				client: mgclients.Client{
					ID:   validID,
					Name: valid,
					Credentials: mgclients.Credentials{
						Identity: "example@example.com",
						Secret:   valid,
					},
				},
			},
			err: nil,
		},
		{
			desc: "empty token",
			req: createClientReq{
				token: "",
				client: mgclients.Client{
					ID:   validID,
					Name: valid,
					Credentials: mgclients.Credentials{
						Identity: "example@example.com",
						Secret:   valid,
					},
				},
			},
		},
		{
			desc: "name too long",
			req: createClientReq{
				token: valid,
				client: mgclients.Client{
					ID:   validID,
					Name: strings.Repeat("a", api.MaxNameSize+1),
				},
			},
			err: apiutil.ErrNameSize,
		},
	}
	for _, tc := range cases {
		err := tc.req.validate()
		assert.Equal(t, tc.err, err)
	}
}

func TestViewClientReqValidate(t *testing.T) {
	cases := []struct {
		desc string
		req  viewClientReq
		err  error
	}{
		{
			desc: "valid request",
			req: viewClientReq{
				token: valid,
				id:    validID,
			},
			err: nil,
		},
		{
			desc: "empty token",
			req: viewClientReq{
				token: "",
				id:    validID,
			},
			err: apiutil.ErrBearerToken,
		},
		{
			desc: "empty id",
			req: viewClientReq{
				token: valid,
				id:    "",
			},
			err: apiutil.ErrMissingID,
		},
	}
	for _, c := range cases {
		err := c.req.validate()
		assert.Equal(t, c.err, err, "%s: expected %s got %s\n", c.desc, c.err, err)
	}
}

func TestViewProfileReqValidate(t *testing.T) {
	cases := []struct {
		desc string
		req  viewProfileReq
		err  error
	}{
		{
			desc: "valid request",
			req: viewProfileReq{
				token: valid,
			},
			err: nil,
		},
		{
			desc: "empty token",
			req: viewProfileReq{
				token: "",
			},
			err: apiutil.ErrBearerToken,
		},
	}
	for _, c := range cases {
		err := c.req.validate()
		assert.Equal(t, c.err, err)
	}
}

func TestListClientsReqValidate(t *testing.T) {
	cases := []struct {
		desc string
		req  listClientsReq
		err  error
	}{
		{
			desc: "valid request",
			req: listClientsReq{
				token: valid,
				limit: 10,
			},
			err: nil,
		},
		{
			desc: "empty token",
			req: listClientsReq{
				token: "",
				limit: 10,
			},
			err: apiutil.ErrBearerToken,
		},
		{
			desc: "limit too big",
			req: listClientsReq{
				token: valid,
				limit: api.MaxLimitSize + 1,
			},
			err: apiutil.ErrLimitSize,
		},
		{
			desc: "limit too small",
			req: listClientsReq{
				token: valid,
				limit: 0,
			},
			err: apiutil.ErrLimitSize,
		},
		{
			desc: "invalid visibility",
			req: listClientsReq{
				token:      valid,
				limit:      10,
				visibility: "invalid",
			},
			err: apiutil.ErrInvalidVisibilityType,
		},
		{
			desc: "invalid direction",
			req: listClientsReq{
				token: valid,
				limit: 10,
				dir:   "invalid",
			},
			err: apiutil.ErrInvalidDirection,
		},
	}
	for _, c := range cases {
		err := c.req.validate()
		assert.Equal(t, c.err, err, "%s: expected %s got %s\n", c.desc, c.err, err)
	}
}

func TestListMembersByObjectReqValidate(t *testing.T) {
	cases := []struct {
		desc string
		req  listMembersByObjectReq
		err  error
	}{
		{
			desc: "valid request",
			req: listMembersByObjectReq{
				token:      valid,
				objectKind: "group",
				objectID:   validID,
			},
			err: nil,
		},
		{
			desc: "empty token",
			req: listMembersByObjectReq{
				token:      "",
				objectKind: "group",
				objectID:   validID,
			},
			err: apiutil.ErrBearerToken,
		},
		{
			desc: "empty object kind",
			req: listMembersByObjectReq{
				token:      valid,
				objectKind: "",
				objectID:   validID,
			},
			err: apiutil.ErrMissingMemberKind,
		},
		{
			desc: "empty object id",
			req: listMembersByObjectReq{
				token:      valid,
				objectKind: "group",
				objectID:   "",
			},
			err: apiutil.ErrMissingID,
		},
	}
	for _, c := range cases {
		err := c.req.validate()
		assert.Equal(t, c.err, err)
	}
}

func TestUpdateClientReqValidate(t *testing.T) {
	cases := []struct {
		desc string
		req  updateClientReq
		err  error
	}{
		{
			desc: "valid request",
			req: updateClientReq{
				token: valid,
				id:    validID,
				Name:  valid,
			},
			err: nil,
		},
		{
			desc: "empty token",
			req: updateClientReq{
				token: "",
				id:    validID,
				Name:  valid,
			},
			err: apiutil.ErrBearerToken,
		},
		{
			desc: "empty id",
			req: updateClientReq{
				token: valid,
				id:    "",
				Name:  valid,
			},
			err: apiutil.ErrMissingID,
		},
	}
	for _, c := range cases {
		err := c.req.validate()
		assert.Equal(t, c.err, err, "%s: expected %s got %s\n", c.desc, c.err, err)
	}
}

func TestUpdateClientTagsReqValidate(t *testing.T) {
	cases := []struct {
		desc string
		req  updateClientTagsReq
		err  error
	}{
		{
			desc: "valid request",
			req: updateClientTagsReq{
				token: valid,
				id:    validID,
				Tags:  []string{"tag1", "tag2"},
			},
			err: nil,
		},
		{
			desc: "empty token",
			req: updateClientTagsReq{
				token: "",
				id:    validID,
				Tags:  []string{"tag1", "tag2"},
			},
			err: apiutil.ErrBearerToken,
		},
		{
			desc: "empty id",
			req: updateClientTagsReq{
				token: valid,
				id:    "",
				Tags:  []string{"tag1", "tag2"},
			},
			err: apiutil.ErrMissingID,
		},
	}
	for _, c := range cases {
		err := c.req.validate()
		assert.Equal(t, c.err, err, "%s: expected %s got %s\n", c.desc, c.err, err)
	}
}

func TestUpdateClientRoleReqValidate(t *testing.T) {
	cases := []struct {
		desc string
		req  updateClientRoleReq
		err  error
	}{
		{
			desc: "valid request",
			req: updateClientRoleReq{
				token: valid,
				id:    validID,
				Role:  "admin",
			},
			err: nil,
		},
		{
			desc: "empty token",
			req: updateClientRoleReq{
				token: "",
				id:    validID,
				Role:  "admin",
			},
			err: apiutil.ErrBearerToken,
		},
		{
			desc: "empty id",
			req: updateClientRoleReq{
				token: valid,
				id:    "",
				Role:  "admin",
			},
			err: apiutil.ErrMissingID,
		},
	}
	for _, c := range cases {
		err := c.req.validate()
		assert.Equal(t, c.err, err, "%s: expected %s got %s\n", c.desc, c.err, err)
	}
}

func TestUpdateClientIdentityReqValidate(t *testing.T) {
	cases := []struct {
		desc string
		req  updateClientIdentityReq
		err  error
	}{
		{
			desc: "valid request",
			req: updateClientIdentityReq{
				token:    valid,
				id:       validID,
				Identity: "example@example.com",
			},
			err: nil,
		},
		{
			desc: "empty token",
			req: updateClientIdentityReq{
				token:    "",
				id:       validID,
				Identity: "example@example.com",
			},
			err: apiutil.ErrBearerToken,
		},
		{
			desc: "empty id",
			req: updateClientIdentityReq{
				token:    valid,
				id:       "",
				Identity: "example@example.com",
			},
			err: apiutil.ErrMissingID,
		},
	}
	for _, c := range cases {
		err := c.req.validate()
		assert.Equal(t, c.err, err, "%s: expected %s got %s\n", c.desc, c.err, err)
	}
}

func TestUpdateClientSecretReqValidate(t *testing.T) {
	cases := []struct {
		desc string
		req  updateClientSecretReq
		err  error
	}{
		{
			desc: "valid request",
			req: updateClientSecretReq{
				token:     valid,
				OldSecret: valid,
				NewSecret: valid,
			},
			err: nil,
		},
		{
			desc: "empty token",
			req: updateClientSecretReq{
				token:     "",
				OldSecret: valid,
				NewSecret: valid,
			},
			err: apiutil.ErrBearerToken,
		},
	}
	for _, c := range cases {
		err := c.req.validate()
		assert.Equal(t, c.err, err)
	}
}

func TestChangeClientStatusReqValidate(t *testing.T) {
	cases := []struct {
		desc string
		req  changeClientStatusReq
		err  error
	}{
		{
			desc: "valid request",
			req: changeClientStatusReq{
				token: valid,
				id:    validID,
			},
			err: nil,
		},
		{
			desc: "empty token",
			req: changeClientStatusReq{
				token: "",
				id:    validID,
			},
			err: apiutil.ErrBearerToken,
		},
		{
			desc: "empty id",
			req: changeClientStatusReq{
				token: valid,
				id:    "",
			},
			err: apiutil.ErrMissingID,
		},
	}
	for _, c := range cases {
		err := c.req.validate()
		assert.Equal(t, c.err, err, "%s: expected %s got %s\n", c.desc, c.err, err)
	}
}

func TestLoginClientReqValidate(t *testing.T) {
	cases := []struct {
		desc string
		req  loginClientReq
		err  error
	}{
		{
			desc: "valid request",
			req: loginClientReq{
				Identity: "eaxmple,example.com",
				Secret:   valid,
			},
			err: nil,
		},
		{
			desc: "empty identity",
			req: loginClientReq{
				Identity: "",
				Secret:   valid,
			},
			err: apiutil.ErrMissingIdentity,
		},
		{
			desc: "empty secret",
			req: loginClientReq{
				Identity: "eaxmple,example.com",
				Secret:   "",
			},
			err: apiutil.ErrMissingSecret,
		},
	}
	for _, c := range cases {
		err := c.req.validate()
		assert.Equal(t, c.err, err, "%s: expected %s got %s\n", c.desc, c.err, err)
	}
}

func TestTokenReqValidate(t *testing.T) {
	cases := []struct {
		desc string
		req  tokenReq
		err  error
	}{
		{
			desc: "valid request",
			req: tokenReq{
				RefreshToken: valid,
			},
			err: nil,
		},
		{
			desc: "empty token",
			req: tokenReq{
				RefreshToken: "",
			},
			err: apiutil.ErrBearerToken,
		},
	}
	for _, c := range cases {
		err := c.req.validate()
		assert.Equal(t, c.err, err, "%s: expected %s got %s\n", c.desc, c.err, err)
	}
}

func TestPasswResetReqValidate(t *testing.T) {
	cases := []struct {
		desc string
		req  passwResetReq
		err  error
	}{
		{
			desc: "valid request",
			req: passwResetReq{
				Email: "example@example.com",
				Host:  "example.com",
			},
			err: nil,
		},
		{
			desc: "empty email",
			req: passwResetReq{
				Email: "",
				Host:  "example.com",
			},
			err: apiutil.ErrMissingEmail,
		},
		{
			desc: "empty host",
			req: passwResetReq{
				Email: "example@example.com",
				Host:  "",
			},
			err: apiutil.ErrMissingHost,
		},
	}
	for _, c := range cases {
		err := c.req.validate()
		assert.Equal(t, c.err, err, "%s: expected %s got %s\n", c.desc, c.err, err)
	}
}

func TestResetTokenReqValidate(t *testing.T) {
	cases := []struct {
		desc string
		req  resetTokenReq
		err  error
	}{
		{
			desc: "valid request",
			req: resetTokenReq{
				Token:    valid,
				Password: valid,
				ConfPass: valid,
			},
			err: nil,
		},
		{
			desc: "empty token",
			req: resetTokenReq{
				Token:    "",
				Password: valid,
				ConfPass: valid,
			},
			err: apiutil.ErrBearerToken,
		},
		{
			desc: "empty password",
			req: resetTokenReq{
				Token:    valid,
				Password: "",
				ConfPass: valid,
			},
			err: apiutil.ErrMissingPass,
		},
		{
			desc: "empty confpass",
			req: resetTokenReq{
				Token:    valid,
				Password: valid,
				ConfPass: "",
			},
			err: apiutil.ErrMissingConfPass,
		},
		{
			desc: "mismatching password and confpass",
			req: resetTokenReq{
				Token:    valid,
				Password: "valid2",
				ConfPass: valid,
			},
			err: apiutil.ErrInvalidResetPass,
		},
	}
	for _, c := range cases {
		err := c.req.validate()
		assert.Equal(t, c.err, err)
	}
}

func TestAssignUsersRequestValidate(t *testing.T) {
	cases := []struct {
		desc string
		req  assignUsersReq
		err  error
	}{
		{
			desc: "valid request",
			req: assignUsersReq{
				token:    valid,
				groupID:  validID,
				UserIDs:  []string{validID},
				Relation: valid,
			},
			err: nil,
		},
		{
			desc: "empty token",
			req: assignUsersReq{
				token:    "",
				groupID:  validID,
				UserIDs:  []string{validID},
				Relation: valid,
			},
			err: apiutil.ErrBearerToken,
		},
		{
			desc: "empty id",
			req: assignUsersReq{
				token:    valid,
				groupID:  "",
				UserIDs:  []string{validID},
				Relation: valid,
			},
			err: apiutil.ErrMissingID,
		},
		{
			desc: "empty users",
			req: assignUsersReq{
				token:    valid,
				groupID:  validID,
				UserIDs:  []string{},
				Relation: valid,
			},
			err: apiutil.ErrEmptyList,
		},
		{
			desc: "empty relation",
			req: assignUsersReq{
				token:    valid,
				groupID:  validID,
				UserIDs:  []string{validID},
				Relation: "",
			},
			err: apiutil.ErrMissingRelation,
		},
	}
	for _, c := range cases {
		err := c.req.validate()
		assert.Equal(t, c.err, err, "%s: expected %s got %s\n", c.desc, c.err, err)
	}
}

func TestUnassignUsersRequestValidate(t *testing.T) {
	cases := []struct {
		desc string
		req  unassignUsersReq
		err  error
	}{
		{
			desc: "valid request",
			req: unassignUsersReq{
				token:    valid,
				groupID:  validID,
				UserIDs:  []string{validID},
				Relation: valid,
			},
			err: nil,
		},
		{
			desc: "empty token",
			req: unassignUsersReq{
				token:    "",
				groupID:  validID,
				UserIDs:  []string{validID},
				Relation: valid,
			},
			err: apiutil.ErrBearerToken,
		},
		{
			desc: "empty id",
			req: unassignUsersReq{
				token:    valid,
				groupID:  "",
				UserIDs:  []string{validID},
				Relation: valid,
			},
			err: apiutil.ErrMissingID,
		},
		{
			desc: "empty users",
			req: unassignUsersReq{
				token:    valid,
				groupID:  validID,
				UserIDs:  []string{},
				Relation: valid,
			},
			err: apiutil.ErrEmptyList,
		},
		{
			desc: "empty relation",
			req: unassignUsersReq{
				token:    valid,
				groupID:  validID,
				UserIDs:  []string{validID},
				Relation: "",
			},
			err: apiutil.ErrMissingRelation,
		},
	}
	for _, c := range cases {
		err := c.req.validate()
		assert.Equal(t, c.err, err, "%s: expected %s got %s\n", c.desc, c.err, err)
	}
}

func TestAssignGroupsRequestValidate(t *testing.T) {
	cases := []struct {
		desc string
		req  assignGroupsReq
		err  error
	}{
		{
			desc: "valid request",
			req: assignGroupsReq{
				token:    valid,
				groupID:  validID,
				GroupIDs: []string{validID},
			},
			err: nil,
		},
		{
			desc: "empty token",
			req: assignGroupsReq{
				token:    "",
				groupID:  validID,
				GroupIDs: []string{validID},
			},
			err: apiutil.ErrBearerToken,
		},
		{
			desc: "empty group id",
			req: assignGroupsReq{
				token:    valid,
				groupID:  "",
				GroupIDs: []string{validID},
			},
			err: apiutil.ErrMissingID,
		},
		{
			desc: "empty user group ids",
			req: assignGroupsReq{
				token:    valid,
				groupID:  validID,
				GroupIDs: []string{},
			},
			err: apiutil.ErrEmptyList,
		},
	}
	for _, c := range cases {
		err := c.req.validate()
		assert.Equal(t, c.err, err, "%s: expected %s got %s\n", c.desc, c.err, err)
	}
}

func TestUnassignGroupsRequestValidate(t *testing.T) {
	cases := []struct {
		desc string
		req  unassignGroupsReq
		err  error
	}{
		{
			desc: "valid request",
			req: unassignGroupsReq{
				token:    valid,
				groupID:  validID,
				GroupIDs: []string{validID},
			},
			err: nil,
		},
		{
			desc: "empty token",
			req: unassignGroupsReq{
				token:    "",
				groupID:  validID,
				GroupIDs: []string{validID},
			},
			err: apiutil.ErrBearerToken,
		},
		{
			desc: "empty group id",
			req: unassignGroupsReq{
				token:    valid,
				groupID:  "",
				GroupIDs: []string{validID},
			},
			err: apiutil.ErrMissingID,
		},
		{
			desc: "empty user group ids",
			req: unassignGroupsReq{
				token:    valid,
				groupID:  validID,
				GroupIDs: []string{},
			},
			err: apiutil.ErrEmptyList,
		},
	}
	for _, c := range cases {
		err := c.req.validate()
		assert.Equal(t, c.err, err, "%s: expected %s got %s\n", c.desc, c.err, err)
	}
}
