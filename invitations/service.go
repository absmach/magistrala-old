// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package invitations

import (
	"context"
	"errors"
	"time"

	"github.com/absmach/magistrala"
	"github.com/absmach/magistrala/auth"
	svcerr "github.com/absmach/magistrala/pkg/errors/service"
	mgsdk "github.com/absmach/magistrala/pkg/sdk/go"
)

const defLimit = 10

type service struct {
	repo Repository
	auth magistrala.AuthServiceClient
	sdk  mgsdk.SDK
}

func NewService(repo Repository, auth magistrala.AuthServiceClient, sdk mgsdk.SDK) Service {
	return &service{
		repo: repo,
		auth: auth,
		sdk:  sdk,
	}
}

func (svc *service) SendInvitation(ctx context.Context, token string, invitation Invitation) error {
	if err := CheckRelation(invitation.Relation); err != nil {
		return err
	}

	userID, err := svc.identify(ctx, token)
	if err != nil {
		return err
	}
	invitation.InvitedBy = userID

	if err := svc.checkDomainOrPlatformAdmin(ctx, userID, invitation.Domain); err != nil {
		return err
	}

	joinToken, err := svc.auth.Issue(ctx, &magistrala.IssueReq{UserId: userID, DomainId: &invitation.Domain, Type: uint32(auth.InvitationKey)})
	if err != nil {
		return err
	}
	invitation.Token = joinToken.GetAccessToken()

	if invitation.Resend {
		invitation.UpdatedAt = time.Now()

		return svc.repo.UpdateToken(ctx, invitation)
	}

	invitation.CreatedAt = time.Now()

	return svc.repo.Create(ctx, invitation)
}

func (svc *service) ViewInvitation(ctx context.Context, token, userID, domain string) (invitation Invitation, err error) {
	tokenUserID, err := svc.identify(ctx, token)
	if err != nil {
		return Invitation{}, err
	}
	inv, err := svc.repo.Retrieve(ctx, userID, domain)
	if err != nil {
		return Invitation{}, err
	}

	if tokenUserID == userID {
		return inv, nil
	}

	if inv.InvitedBy == tokenUserID {
		return inv, nil
	}

	if err := svc.checkDomainOrPlatformAdmin(ctx, tokenUserID, domain); err != nil {
		return Invitation{}, err
	}

	return inv, nil
}

func (svc *service) ListInvitations(ctx context.Context, token string, page Page) (invitations InvitationPage, err error) {
	userID, err := svc.identify(ctx, token)
	if err != nil {
		return InvitationPage{}, err
	}

	if err := svc.authorize(ctx, auth.UserType, auth.UsersKind, userID, auth.AdminPermission, auth.PlatformType, auth.MagistralaObject); err == nil {
		return svc.repo.RetrieveAll(ctx, page)
	}

	if page.Domain != "" {
		if err := svc.checkDomainOrPlatformAdmin(ctx, userID, page.Domain); err != nil {
			return InvitationPage{}, err
		}

		return svc.repo.RetrieveAll(ctx, page)
	}

	page.InvitedBy = userID

	return svc.repo.RetrieveAll(ctx, page)
}

func (svc *service) AcceptInvitation(ctx context.Context, token string) ([]string, error) {
	userID, err := svc.identify(ctx, token)
	if err != nil {
		return nil, err
	}

	page, err := svc.repo.RetrieveAll(ctx, Page{UserID: userID, Offset: 0, Limit: defLimit})
	if err != nil {
		return nil, err
	}
	if page.Total > defLimit {
		for i := uint64(defLimit); i < page.Total; i += defLimit {
			p, err := svc.repo.RetrieveAll(ctx, Page{UserID: userID, Offset: i, Limit: defLimit})
			if err != nil {
				return nil, err
			}
			page.Invitations = append(page.Invitations, p.Invitations...)
		}
	}

	domains := []string{}
	for _, inv := range page.Invitations {
		if inv.UserID == userID && inv.ConfirmedAt.IsZero() {
			req := mgsdk.UsersRelationRequest{
				Relation: inv.Relation,
				UserIDs:  []string{userID},
			}
			if sdkerr := svc.sdk.AddUserToDomain(inv.Domain, req, inv.Token); sdkerr != nil {
				return domains, sdkerr
			}

			inv.ConfirmedAt = time.Now()
			if err := svc.repo.UpdateConfirmation(ctx, inv); err != nil {
				return domains, err
			}

			domains = append(domains, inv.Domain)
		}
	}

	return domains, nil
}

func (svc *service) DeleteInvitation(ctx context.Context, token, userID, domainID string) error {
	tokenUserID, err := svc.identify(ctx, token)
	if err != nil {
		return err
	}
	if tokenUserID == userID {
		return svc.repo.Delete(ctx, userID, domainID)
	}

	inv, err := svc.repo.Retrieve(ctx, userID, domainID)
	if err != nil {
		return err
	}
	if inv.InvitedBy == tokenUserID {
		return svc.repo.Delete(ctx, userID, domainID)
	}

	if err := svc.checkDomainOrPlatformAdmin(ctx, tokenUserID, domainID); err != nil {
		return err
	}

	return svc.repo.Delete(ctx, userID, domainID)
}

func (svc *service) identify(ctx context.Context, token string) (string, error) {
	user, err := svc.auth.Identify(ctx, &magistrala.IdentityReq{Token: token})
	if err != nil {
		return "", err
	}

	return user.GetUserId(), nil
}

func (svc *service) authorize(ctx context.Context, subjType, subjKind, subj, perm, objType, obj string) error {
	req := &magistrala.AuthorizeReq{
		SubjectType: subjType,
		SubjectKind: subjKind,
		Subject:     subj,
		Permission:  perm,
		ObjectType:  objType,
		Object:      obj,
	}
	res, err := svc.auth.Authorize(ctx, req)
	if err != nil {
		return errors.Join(svcerr.ErrAuthorization, err)
	}

	if !res.GetAuthorized() {
		return errors.Join(svcerr.ErrAuthorization, err)
	}

	return nil
}

// checkDomainOrPlatformAdmin checks if the given user is a domain or platform administrator.
func (svc *service) checkDomainOrPlatformAdmin(ctx context.Context, userID, domainID string) error {
	if err := svc.authorize(ctx, auth.UserType, auth.UsersKind, userID, auth.AdminPermission, auth.DomainType, domainID); err == nil {
		return nil
	}
	if err := svc.authorize(ctx, auth.UserType, auth.UsersKind, userID, auth.AdminPermission, auth.PlatformType, auth.MagistralaObject); err == nil {
		return nil
	}

	return svcerr.ErrAuthorization
}
