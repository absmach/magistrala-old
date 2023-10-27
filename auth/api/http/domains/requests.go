// Copyright (c) Magistrala
// SPDX-License-Identifier: Apache-2.0

package domains

import (
	"github.com/absmach/magistrala/internal/apiutil"
)

type createDomainReq struct {
	token       string
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	Alias       string                 `json:"alias,omitempty"`
}

func (req createDomainReq) validate() error {
	if req.token == "" {
		return apiutil.ErrBearerToken
	}

	if req.Name == "" {
		return apiutil.ErrMissingName
	}
	return nil
}

type viewDomainRequest struct {
	token    string
	domainID string
}

func (req viewDomainRequest) validate() error {
	if req.token == "" {
		return apiutil.ErrBearerToken
	}

	if req.domainID == "" {
		return apiutil.ErrMissingID
	}

	return nil
}

type updateDomainReq struct {
	token       string
	domainID    string
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	Alias       string                 `json:"alias,omitempty"`
}

func (req updateDomainReq) validate() error {
	if req.token == "" {
		return apiutil.ErrBearerToken
	}

	if req.domainID == "" {
		return apiutil.ErrMissingID
	}

	return nil
}

type assignUsersReq struct {
	token    string
	domainID string
	UserIDs  []string `json:"user_ids"`
	Relation string   `json:"relation"`
}

func (req assignUsersReq) validate() error {
	if req.token == "" {
		return apiutil.ErrBearerToken
	}

	if req.domainID == "" {
		return apiutil.ErrMissingID
	}

	if len(req.UserIDs) == 0 {
		return apiutil.ErrMalformedPolicy
	}

	if req.Relation == "" {
		return apiutil.ErrMalformedPolicy
	}

	return nil
}

type unassignUsersReq struct {
	token    string
	domainID string
	UserIDs  []string `json:"user_ids"`
	Relation string   `json:"relation"`
}

func (req unassignUsersReq) validate() error {
	if req.token == "" {
		return apiutil.ErrBearerToken
	}

	if req.domainID == "" {
		return apiutil.ErrMissingID
	}

	if len(req.UserIDs) == 0 {
		return apiutil.ErrMalformedPolicy
	}

	if req.Relation == "" {
		return apiutil.ErrMalformedPolicy
	}

	return nil
}
