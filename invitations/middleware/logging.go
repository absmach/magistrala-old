// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/absmach/magistrala/invitations"
	mglog "github.com/absmach/magistrala/logger"
)

var _ invitations.Service = (*logging)(nil)

type logging struct {
	logger mglog.Logger
	svc    invitations.Service
}

func Logging(logger mglog.Logger, svc invitations.Service) invitations.Service {
	return &logging{logger, svc}
}

func (lm *logging) SendInvitation(ctx context.Context, token string, invitation invitations.Invitation) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method send_invitation to user_id %s from domain %s took %s to complete", invitation.UserID, invitation.Domain, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())
	return lm.svc.SendInvitation(ctx, token, invitation)
}

func (lm *logging) ViewInvitation(ctx context.Context, token string, userID, domain string) (invitation invitations.Invitation, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method view_invitation took %s to complete", time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())
	return lm.svc.ViewInvitation(ctx, token, userID, domain)
}

func (lm *logging) ListInvitations(ctx context.Context, token string, page invitations.Page) (invitations invitations.InvitationPage, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method list_invitations took %s to complete", time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())
	return lm.svc.ListInvitations(ctx, token, page)
}

func (lm *logging) AcceptInvitation(ctx context.Context, token string) (domains []string, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method accept_invitation took %s to complete", time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())
	return lm.svc.AcceptInvitation(ctx, token)
}

func (lm *logging) DeleteInvitation(ctx context.Context, token string, userID, domain string) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method delete_invitation took %s to complete", time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())
	return lm.svc.DeleteInvitation(ctx, token, userID, domain)
}
