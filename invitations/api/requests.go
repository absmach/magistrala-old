// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"errors"

	"github.com/absmach/magistrala/internal/apiutil"
	"github.com/absmach/magistrala/invitations"
)

const maxLimitSize = 100

var errMissingDomain = errors.New("missing domain")

type sendInvitationReq struct {
	token    string
	UserID   string `json:"user_id,omitempty"`
	Domain   string `json:"domain,omitempty"`
	Relation string `json:"relation,omitempty"`
	Resend   bool   `json:"resend,omitempty"`
}

func (req *sendInvitationReq) validate() error {
	if req.token == "" {
		return apiutil.ErrBearerToken
	}
	if req.UserID == "" {
		return apiutil.ErrMissingID
	}
	if req.Domain == "" {
		return errMissingDomain
	}
	if err := invitations.CheckRelation(req.Relation); err != nil {
		return err
	}

	return nil
}

type listInvitationsReq struct {
	token string
	invitations.Page
}

func (req *listInvitationsReq) validate() error {
	if req.token == "" {
		return apiutil.ErrBearerToken
	}
	if req.Page.Limit > maxLimitSize || req.Page.Limit < 1 {
		return apiutil.ErrLimitSize
	}

	return nil
}

type acceptInvitationReq struct {
	token string
}

func (req *acceptInvitationReq) validate() error {
	if req.token == "" {
		return apiutil.ErrBearerToken
	}

	return nil
}

type invitationReq struct {
	token  string
	userID string
	domain string
}

func (req *invitationReq) validate() error {
	if req.token == "" {
		return apiutil.ErrBearerToken
	}
	if req.userID == "" {
		return apiutil.ErrMissingID
	}
	if req.domain == "" {
		return errMissingDomain
	}

	return nil
}
