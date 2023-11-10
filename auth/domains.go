// Copyright (c) Magistrala
// SPDX-License-Identifier: Apache-2.0

package auth

import (
	"context"
	"time"

	"github.com/absmach/magistrala/pkg/clients"
)

type DomainReq struct {
	Name     *string           `json:"name,omitempty"`
	Email    *string           `json:"email,omitempty"`
	Metadata *clients.Metadata `json:"metadata,omitempty"`
	Tags     *[]string         `json:"tags,omitempty"`
	Alias    *string           `json:"alias,omitempty"`
	Status   *clients.Status   `json:"status,omitempty"`
}
type Domain struct {
	ID         string           `json:"id"`
	Name       string           `json:"name"`
	Email      string           `json:"email"`
	Metadata   clients.Metadata `json:"metadata,omitempty"`
	Tags       []string         `json:"tags,omitempty"`
	Alias      string           `json:"alias,omitempty"`
	Status     clients.Status   `json:"status"`
	Permission string           `json:"permission,omitempty"`
	CreatedBy  string           `json:"created_by,omitempty"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedBy  string           `json:"updated_by,omitempty"`
	UpdatedAt  time.Time        `json:"updated_at,omitempty"`
}

type Page struct {
	Total      uint64           `json:"total"`
	Offset     uint64           `json:"offset"`
	Limit      uint64           `json:"limit"`
	Name       string           `json:"name,omitempty"`
	Email      string           `json:"email,omitempty"`
	Order      string           `json:"-"`
	Dir        string           `json:"-"`
	Metadata   clients.Metadata `json:"metadata,omitempty"`
	Tag        string           `json:"tag,omitempty"`
	Permission string           `json:"permission,omitempty"`
	Status     clients.Status   `json:"status,omitempty"`
	ID         string           `json:"id,omitempty"`
	IDs        []string         `json:"-"`
	Identity   string           `json:"identity,omitempty"`
	SubjectID  string           `json:"-"`
}

type DomainsPage struct {
	Page
	Domains []Domain `json:"domains,omitempty"`
}

type PolicyCopy struct {
	SubjectType     string `json:"subject_type,omitempty"`
	SubjectID       string `json:"subject_id,omitempty"`
	SubjectRelation string `json:"subject_relation,omitempty"`
	Relation        string `json:"relation,omitempty"`
	ObjectType      string `json:"object_type,omitempty"`
	ObjectID        string `json:"object_id,omitempty"`
}

type Domains interface {
	CreateDomain(ctx context.Context, token string, d Domain) (Domain, error)
	RetrieveDomain(ctx context.Context, token string, id string) (Domain, error)
	UpdateDomain(ctx context.Context, token string, id string, d DomainReq) (Domain, error)
	ChangeDomainStatus(ctx context.Context, token string, id string, d DomainReq) (Domain, error)
	ListDomains(ctx context.Context, token string, page Page) (DomainsPage, error)
	AssignUsers(ctx context.Context, token string, id string, userIds []string, relation string) error
	UnassignUsers(ctx context.Context, token string, id string, userIds []string, relation string) error
	ListUserDomains(ctx context.Context, token string, userID string, page Page) (DomainsPage, error)
}

type DomainsRepository interface {

	// Save creates db insert transaction for the given domain.
	Save(ctx context.Context, d Domain) (Domain, error)

	// RetrieveByID retrieves Domain by its unique ID.
	RetrieveByID(ctx context.Context, id string) (Domain, error)

	// RetrieveAllByIDs retrieves for given Domain IDs .
	RetrieveAllByIDs(ctx context.Context, pm Page) (DomainsPage, error)

	// Update updates the client name and metadata.
	Update(ctx context.Context, id string, userID string, d DomainReq) (Domain, error)

	// Delete
	Delete(ctx context.Context, id string) error

	// SavePolicy save policy in domains database
	SavePolicy(ctx context.Context, pc PolicyCopy) error

	// DeletePolicy delete policy from domains database
	DeletePolicy(ctx context.Context, pc PolicyCopy) error

	//ListDomains
	ListDomains(ctx context.Context, pm Page) (DomainsPage, error)
}
