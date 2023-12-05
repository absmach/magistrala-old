// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package invitations

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/absmach/magistrala/auth"
)

var (
	errMissingRelation = errors.New("missing relation")
	errInvalidRelation = errors.New("invalid relation")
)

// Invitation is an invitation to join a domain.
type Invitation struct {
	InvitedBy   string    `json:"invited_by" db:"invited_by"`
	UserID      string    `json:"user_id" db:"user_id"`
	Domain      string    `json:"domain" db:"domain"`
	Token       string    `json:"token,omitempty" db:"token"`
	Relation    string    `json:"relation,omitempty" db:"relation"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" db:"updated_at,omitempty"`
	ConfirmedAt time.Time `json:"confirmed_at,omitempty" db:"confirmed_at,omitempty"`
}

// Page is a page of invitations.
type Page struct {
	Offset    uint64 `json:"offset" db:"offset"`
	Limit     uint64 `json:"limit" db:"limit"`
	InvitedBy string `json:"invited_by,omitempty" db:"invited_by,omitempty"`
	UserID    string `json:"user_id,omitempty" db:"user_id,omitempty"`
	Domain    string `json:"domain,omitempty" db:"domain,omitempty"`
	Relation  string `json:"relation,omitempty" db:"relation,omitempty"`
}

// InvitationPage is a page of invitations.
type InvitationPage struct {
	Total       uint64       `json:"total"`
	Offset      uint64       `json:"offset"`
	Limit       uint64       `json:"limit"`
	Invitations []Invitation `json:"invitations"`
}

func (page InvitationPage) MarshalJSON() ([]byte, error) {
	type Alias InvitationPage
	a := struct {
		Alias
	}{
		Alias: Alias(page),
	}

	if a.Invitations == nil {
		a.Invitations = make([]Invitation, 0)
	}

	return json.Marshal(a)
}

// Service is an interface that defines methods for managing invitations.
type Service interface {
	// SendInvitation sends an invitation to the email address associated with the given user.
	SendInvitation(ctx context.Context, token, host string, invitation Invitation) (err error)

	// ViewInvitation returns an invitation.
	ViewInvitation(ctx context.Context, token string, userID, domain string) (invitation Invitation, err error)

	// ListInvitations returns a list of invitations.
	ListInvitations(ctx context.Context, token string, page Page) (invitations InvitationPage, err error)

	// AcceptInvitation accepts an invitation by adding the user to the domain.
	AcceptInvitation(ctx context.Context, token string) (userID string, domains []string, err error)

	// DeleteInvitation deletes an invitation.
	DeleteInvitation(ctx context.Context, token string, userID, domain string) (err error)
}

type Repository interface {
	// Create creates an invitation.
	Create(ctx context.Context, invitation Invitation) (err error)

	// Retrieve returns an invitation.
	Retrieve(ctx context.Context, userID, domainID string) (Invitation, error)

	// RetrieveAll returns a list of invitations based on the given page.
	RetrieveAll(ctx context.Context, withToken bool, page Page) (invitations InvitationPage, err error)

	// UpdateToken updates an invitation by setting the token.
	UpdateToken(ctx context.Context, invitation Invitation) (err error)

	// UpdateConfirmation updates an invitation by setting the confirmation time.
	UpdateConfirmation(ctx context.Context, invitation Invitation) (err error)

	// Delete deletes an invitation.
	Delete(ctx context.Context, userID, domainID string) (err error)
}

// CheckRelation checks if the given relation is valid.
// It returns an error if the relation is empty or invalid.
func CheckRelation(relation string) error {
	if relation == "" {
		return errMissingRelation
	}
	if relation != auth.AdministratorRelation &&
		relation != auth.EditorRelation &&
		relation != auth.ViewerRelation &&
		relation != auth.MemberRelation &&
		relation != auth.DomainRelation &&
		relation != auth.ParentGroupRelation &&
		relation != auth.RoleGroupRelation &&
		relation != auth.GroupRelation &&
		relation != auth.PlatformRelation {
		return errInvalidRelation
	}

	return nil
}
