// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/absmach/magistrala"
	mgclients "github.com/absmach/magistrala/pkg/clients"
	"github.com/absmach/magistrala/users"
)

var _ users.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger *slog.Logger
	svc    users.Service
}

// LoggingMiddleware adds logging facilities to the clients service.
func LoggingMiddleware(svc users.Service, logger *slog.Logger) users.Service {
	return &loggingMiddleware{logger, svc}
}

// RegisterClient logs the register_client request. It logs the client id and token and the time it took to complete the request.
// If the request fails, it logs the error.
func (lm *loggingMiddleware) RegisterClient(ctx context.Context, token string, client mgclients.Client) (c mgclients.Client, err error) {
	defer func(begin time.Time) {
		message := "Method register_client completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error: %s.", message, err),
				slog.String("method", "register_client"),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("method", "register_client"),
			slog.String("client_id", c.ID),
			slog.String("token", token),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.RegisterClient(ctx, token, client)
}

// IssueToken logs the issue_token request. It logs the client identity and token type and the time it took to complete the request.
// If the request fails, it logs the error.
func (lm *loggingMiddleware) IssueToken(ctx context.Context, identity, secret, domainID string) (t *magistrala.Token, err error) {
	defer func(begin time.Time) {
		message := "Method issue_token"
		if t != nil {
			message = fmt.Sprintf("%s of type %s", message, t.AccessType)
		}
		message = fmt.Sprintf("%s completed", message)
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error", message),
				slog.String("method", "issue_token"),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("method", "issue_token"),
			slog.String("identity", identity),
			slog.String("domain_id", domainID),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.IssueToken(ctx, identity, secret, domainID)
}

// RefreshToken logs the refresh_token request. It logs the refreshtoken, token type and the time it took to complete the request.
// If the request fails, it logs the error.
func (lm *loggingMiddleware) RefreshToken(ctx context.Context, refreshToken, domainID string) (t *magistrala.Token, err error) {
	defer func(begin time.Time) {
		message := "Method refresh_token"
		if t != nil {
			message = fmt.Sprintf("%s of type %s", message, t.AccessType)
		}
		message = fmt.Sprintf("%s completed", message)
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error", message),
				slog.String("method", "refresh_token"),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("method", "refresh_token"),
			slog.String("refresh_token", refreshToken),
			slog.String("domain_id", domainID),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.RefreshToken(ctx, refreshToken, domainID)
}

// ViewClient logs the view_client request. It logs the client id and token and the time it took to complete the request.
// If the request fails, it logs the error.
func (lm *loggingMiddleware) ViewClient(ctx context.Context, token, id string) (c mgclients.Client, err error) {
	defer func(begin time.Time) {
		message := "Method view_client completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error", message),
				slog.String("method", "view_client"),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("method", "view_client"),
			slog.String("client_id", id),
			slog.String("token", token),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.ViewClient(ctx, token, id)
}

// ViewProfile logs the view_profile request. It logs the client id and token and the time it took to complete the request.
// If the request fails, it logs the error.
func (lm *loggingMiddleware) ViewProfile(ctx context.Context, token string) (c mgclients.Client, err error) {
	defer func(begin time.Time) {
		message := "Method view_profile completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error", message),
				slog.String("method", "view_profile"),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("method", "view_profile"),
			slog.String("token", token),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.ViewProfile(ctx, token)
}

// ListClients logs the list_clients request. It logs the token and page metadata and the time it took to complete the request.
// If the request fails, it logs the error.
func (lm *loggingMiddleware) ListClients(ctx context.Context, token string, pm mgclients.Page) (cp mgclients.ClientsPage, err error) {
	defer func(begin time.Time) {
		message := "Method list_clients completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error", message),
				slog.String("method", "list_clients"),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("method", "list_clients"),
			slog.String("token", token),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.ListClients(ctx, token, pm)
}

// UpdateClient logs the update_client request. It logs the client id and token and the time it took to complete the request.
// If the request fails, it logs the error.
func (lm *loggingMiddleware) UpdateClient(ctx context.Context, token string, client mgclients.Client) (c mgclients.Client, err error) {
	defer func(begin time.Time) {
		message := "Method update_client completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error", message),
				slog.String("method", "update_client"),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("method", "update_client"),
			slog.String("client_id", c.ID),
			slog.String("token", token),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.UpdateClient(ctx, token, client)
}

// UpdateClientTags logs the update_client_tags request. It logs the client id and token and the time it took to complete the request.
// If the request fails, it logs the error.
func (lm *loggingMiddleware) UpdateClientTags(ctx context.Context, token string, client mgclients.Client) (c mgclients.Client, err error) {
	defer func(begin time.Time) {
		message := "Method update_client_tags completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("method", "update_client_tags"),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("method", "update_client_tags"),
			slog.String("client_id", c.ID),
			slog.String("token", token),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.UpdateClientTags(ctx, token, client)
}

// UpdateClientIdentity logs the update_client_identity request. It logs the client id and token and the time it took to complete the request.
// If the request fails, it logs the error.
func (lm *loggingMiddleware) UpdateClientIdentity(ctx context.Context, token, id, identity string) (c mgclients.Client, err error) {
	defer func(begin time.Time) {
		message := "Method update_client_identity completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("method", "update_client_identity"),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("method", "update_client_identity"),
			slog.String("client_id", c.ID),
			slog.String("token", token),
			slog.String("identity", identity),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.UpdateClientIdentity(ctx, token, id, identity)
}

// UpdateClientSecret logs the update_client_secret request. It logs the client id and token and the time it took to complete the request.
// If the request fails, it logs the error.
func (lm *loggingMiddleware) UpdateClientSecret(ctx context.Context, token, oldSecret, newSecret string) (c mgclients.Client, err error) {
	defer func(begin time.Time) {
		message := "Method update_client_secret completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("method", "update_client_secret"),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("method", "update_client_secret"),
			slog.String("client_id", c.ID),
			slog.String("token", token),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.UpdateClientSecret(ctx, token, oldSecret, newSecret)
}

// GenerateResetToken logs the generate_reset_token request. It logs the email and host and the time it took to complete the request.
// If the request fails, it logs the error.
func (lm *loggingMiddleware) GenerateResetToken(ctx context.Context, email, host string) (err error) {
	defer func(begin time.Time) {
		message := "Method generate_reset_token completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("method", "generate_reset_token"),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("method", "generate_reset_token"),
			slog.String("email", email),
			slog.String("host", host),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.GenerateResetToken(ctx, email, host)
}

// ResetSecret logs the reset_secret request. It logs the token and the time it took to complete the request.
// If the request fails, it logs the error.
func (lm *loggingMiddleware) ResetSecret(ctx context.Context, token, secret string) (err error) {
	defer func(begin time.Time) {
		message := "Method reset_secret completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("method", "reset_secret"),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("method", "reset_secret"),
			slog.String("token", token),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.ResetSecret(ctx, token, secret)
}

// SendPasswordReset logs the send_password_reset request. It logs the token and the time it took to complete the request.
// If the request fails, it logs the error.
func (lm *loggingMiddleware) SendPasswordReset(ctx context.Context, host, email, user, token string) (err error) {
	defer func(begin time.Time) {
		message := "Method send_password_reset completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("method", "send_password_reset"),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("method", "send_password_reset"),
			slog.String("host", host),
			slog.String("email", email),
			slog.String("user", user),
			slog.String("token", token),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.SendPasswordReset(ctx, host, email, user, token)
}

// UpdateClientRole logs the update_client_role request. It logs the client id and token and the time it took to complete the request.
// If the request fails, it logs the error.
func (lm *loggingMiddleware) UpdateClientRole(ctx context.Context, token string, client mgclients.Client) (c mgclients.Client, err error) {
	defer func(begin time.Time) {
		message := "Method update_client_role completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("method", "update_client_role"),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("method", "update_client_role"),
			slog.String("client_id", c.ID),
			slog.String("role", client.Role.String()),
			slog.String("token", token),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.UpdateClientRole(ctx, token, client)
}

// EnableClient logs the enable_client request. It logs the client id and token and the time it took to complete the request.
// If the request fails, it logs the error.
func (lm *loggingMiddleware) EnableClient(ctx context.Context, token, id string) (c mgclients.Client, err error) {
	defer func(begin time.Time) {
		message := "Method enable_client completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("method", "enable_client"),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("method", "enable_client"),
			slog.String("client_id", c.ID),
			slog.String("token", token),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.EnableClient(ctx, token, id)
}

// DisableClient logs the disable_client request. It logs the client id and token and the time it took to complete the request.
// If the request fails, it logs the error.
func (lm *loggingMiddleware) DisableClient(ctx context.Context, token, id string) (c mgclients.Client, err error) {
	defer func(begin time.Time) {
		message := "Method disable_client completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("method", "disable_client"),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("method", "disable_client"),
			slog.String("client_id", c.ID),
			slog.String("token", token),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.DisableClient(ctx, token, id)
}

// ListMembers logs the list_members request. It logs the group id, token and the time it took to complete the request.
// If the request fails, it logs the error.
func (lm *loggingMiddleware) ListMembers(ctx context.Context, token, objectKind, objectID string, cp mgclients.Page) (mp mgclients.MembersPage, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method list_members %d members completed", mp.Total)
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("method", "list_members"),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("method", "list_members"),
			slog.String("object_kind", objectKind),
			slog.String("object_id", objectID),
			slog.String("token", token),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.ListMembers(ctx, token, objectKind, objectID, cp)
}

// Identify logs the identify request. It logs the token and the time it took to complete the request.
func (lm *loggingMiddleware) Identify(ctx context.Context, token string) (id string, err error) {
	defer func(begin time.Time) {
		message := "Method identify completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("method", "identify"),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),			
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("method", "identify"),
			slog.String("client_id", id),
			slog.String("token", token),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.Identify(ctx, token)
}
