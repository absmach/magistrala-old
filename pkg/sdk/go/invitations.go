// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package sdk

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/absmach/magistrala/pkg/errors"
)

const (
	invitationsEndpoint = "invitations"
	acceptEndpoint      = "accept"
)

type Invitation struct {
	InvitedBy   string    `json:"invited_by" db:"invited_by"`
	UserID      string    `json:"user_id" db:"user_id"`
	Domain      string    `json:"domain" db:"domain"`
	Token       string    `json:"token,omitempty" db:"token"`
	Relation    string    `json:"relation,omitempty" db:"relation"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" db:"updated_at,omitempty"`
	ConfirmedAt time.Time `json:"confirmed_at,omitempty" db:"confirmed_at,omitempty"`
	Resend      bool      `json:"resend,omitempty"`
}

type InvitationPage struct {
	Total       uint64       `json:"total"`
	Offset      uint64       `json:"offset"`
	Limit       uint64       `json:"limit"`
	Invitations []Invitation `json:"invitations"`
}

func (sdk mgSDK) SendInvitation(invitation Invitation, token string) (err error) {
	data, err := json.Marshal(invitation)
	if err != nil {
		return errors.NewSDKError(err)
	}

	url := sdk.invitationsURL + "/" + invitationsEndpoint

	_, _, sdkerr := sdk.processRequest(http.MethodPost, url, token, data, nil, http.StatusCreated)

	return sdkerr
}

func (sdk mgSDK) Invitation(userID, domainID, token string) (invitation Invitation, err error) {
	url := sdk.invitationsURL + "/" + invitationsEndpoint + "/" + userID + "/" + domainID

	_, body, sdkerr := sdk.processRequest(http.MethodGet, url, token, nil, nil, http.StatusOK)
	if sdkerr != nil {
		return Invitation{}, sdkerr
	}

	if err := json.Unmarshal(body, &invitation); err != nil {
		return Invitation{}, errors.NewSDKError(err)
	}

	return invitation, nil
}

func (sdk mgSDK) Invitations(pm PageMetadata, token string) (invitations InvitationPage, err error) {
	url, err := sdk.withQueryParams(sdk.invitationsURL, invitationsEndpoint, pm)
	if err != nil {
		return InvitationPage{}, errors.NewSDKError(err)
	}

	_, body, sdkerr := sdk.processRequest(http.MethodGet, url, token, nil, nil, http.StatusOK)
	if sdkerr != nil {
		return InvitationPage{}, sdkerr
	}

	var invPage InvitationPage
	if err := json.Unmarshal(body, &invPage); err != nil {
		return InvitationPage{}, errors.NewSDKError(err)
	}

	return invPage, nil
}

func (sdk mgSDK) AcceptInvitation(domain, token string) (err error) {
	req := struct {
		Domain string `json:"domain"`
	}{
		Domain: domain,
	}
	data, err := json.Marshal(req)
	if err != nil {
		return errors.NewSDKError(err)
	}

	url := sdk.invitationsURL + "/" + invitationsEndpoint + "/" + acceptEndpoint

	_, _, sdkerr := sdk.processRequest(http.MethodPost, url, token, data, nil, http.StatusOK)

	return sdkerr
}

func (sdk mgSDK) DeleteInvitation(userID, domain, token string) (err error) {
	url := sdk.invitationsURL + "/" + invitationsEndpoint + "/" + userID + "/" + domain

	_, _, sdkerr := sdk.processRequest(http.MethodDelete, url, token, nil, nil, http.StatusNoContent)

	return sdkerr
}
