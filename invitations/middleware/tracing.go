// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"context"

	"github.com/absmach/magistrala/invitations"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var _ invitations.Service = (*tracing)(nil)

type tracing struct {
	tracer trace.Tracer
	svc    invitations.Service
}

func Tracing(svc invitations.Service, tracer trace.Tracer) invitations.Service {
	return &tracing{tracer, svc}
}

func (tm *tracing) SendInvitation(ctx context.Context, token string, invitation invitations.Invitation) (err error) {
	ctx, span := tm.tracer.Start(ctx, "send_invitation", trace.WithAttributes(
		attribute.String("domain", invitation.Domain),
		attribute.String("user_id", invitation.UserID),
	))
	defer span.End()

	return tm.svc.SendInvitation(ctx, token, invitation)
}

func (tm *tracing) ViewInvitation(ctx context.Context, token, userID, domain string) (invitation invitations.Invitation, err error) {
	ctx, span := tm.tracer.Start(ctx, "view_invitation", trace.WithAttributes(
		attribute.String("user_id", userID),
		attribute.String("domain", domain),
	))
	defer span.End()

	return tm.svc.ViewInvitation(ctx, token, userID, domain)
}

func (tm *tracing) ListInvitations(ctx context.Context, token string, page invitations.Page) (invs invitations.InvitationPage, err error) {
	ctx, span := tm.tracer.Start(ctx, "list_invitations", trace.WithAttributes(
		attribute.Int("limit", int(page.Limit)),
		attribute.Int("offset", int(page.Offset)),
		attribute.String("user_id", page.UserID),
		attribute.String("domain", page.Domain),
		attribute.String("invited_by", page.InvitedBy),
	))
	defer span.End()

	return tm.svc.ListInvitations(ctx, token, page)
}

func (tm *tracing) AcceptInvitation(ctx context.Context, token string) (domains []string, err error) {
	ctx, span := tm.tracer.Start(ctx, "accept_invitation", trace.WithAttributes(
		attribute.StringSlice("domains", domains),
	))
	defer span.End()

	return tm.svc.AcceptInvitation(ctx, token)
}

func (tm *tracing) DeleteInvitation(ctx context.Context, token, userID, domain string) (err error) {
	ctx, span := tm.tracer.Start(ctx, "delete_invitation", trace.WithAttributes(
		attribute.String("user_id", userID),
		attribute.String("domain", domain),
	))
	defer span.End()

	return tm.svc.DeleteInvitation(ctx, token, userID, domain)
}
