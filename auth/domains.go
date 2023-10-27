// Copyright (c) Magistrala
// SPDX-License-Identifier: Apache-2.0

package auth

import (
	"context"
	"time"

	"github.com/absmach/magistrala/pkg/clients"
)

type DomainReq struct {
	Name        *string           `json:"name"`
	Description *string           `json:"description,omitempty"`
	Metadata    *clients.Metadata `json:"metadata,omitempty"`
	Tags        *[]string         `json:"tags,omitempty"`
	Alias       *string           `json:"alias,omitempty"`
	Status      *clients.Status   `json:"status,omitempty"`
}
type Domain struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description,omitempty"`
	Metadata    clients.Metadata `json:"metadata,omitempty"`
	Tags        []string         `json:"tags,omitempty"`
	Alias       string           `json:"alias,omitempty"`
	Status      clients.Status   `json:"status"`
	CreatedBy   string           `json:"created_by,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedBy   string           `json:"updated_by,omitempty"`
	UpdatedAt   time.Time        `json:"updated_at,omitempty"`
}

type Page struct {
	Total      uint64           `json:"total"`
	Offset     uint64           `json:"offset"`
	Limit      uint64           `json:"limit"`
	Name       string           `json:"name,omitempty"`
	Order      string           `json:"order,omitempty"`
	Dir        string           `json:"dir,omitempty"`
	Metadata   clients.Metadata `json:"metadata,omitempty"`
	Owner      string           `json:"owner,omitempty"`
	Tag        string           `json:"tag,omitempty"`
	Permission string           `json:"permission,omitempty"`
	Status     clients.Status   `json:"status,omitempty"`
	IDs        []string         `json:"ids,omitempty"`
	Identity   string           `json:"identity,omitempty"`
}

type DomainsPage struct {
	Page
	Domains []Domain
}
type Domains interface {
	CreateDomain(ctx context.Context, token string, d Domain) (Domain, error)
	ViewDomain(ctx context.Context, token string, id string) (Domain, error)
	UpdateDomain(ctx context.Context, token string, id string, d Domain) (Domain, error)
	ListDomains(ctx context.Context, token string) (DomainsPage, error)
	AssignUsers(ctx context.Context, token string, id string, userIds []string, relation string) error
	UnassignUsers(ctx context.Context, token string, id string, userIds []string, relation string) error
}
