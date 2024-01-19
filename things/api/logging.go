// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"
	"log/slog"
	"time"

	"github.com/absmach/magistrala"
	mgclients "github.com/absmach/magistrala/pkg/clients"
	"github.com/absmach/magistrala/things"
)

var _ things.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger *slog.Logger
	svc    things.Service
}

func LoggingMiddleware(svc things.Service, logger *slog.Logger) things.Service {
	return &loggingMiddleware{logger, svc}
}

func (lm *loggingMiddleware) CreateThings(ctx context.Context, token string, clients ...mgclients.Client) (cs []mgclients.Client, err error) {
	defer func(begin time.Time) {
		args := []interface{}{
			slog.String("duration", time.Since(begin).String()), slog.Any("no_of_things", len(clients)),
		}
		if err != nil {
			args = append(args, slog.String("error", err.Error()))
			lm.logger.Warn("Create thing failed to complete successfully", args...)
			return
		}
		lm.logger.Info("Create thing completed successfully", args...)
	}(time.Now())
	return lm.svc.CreateThings(ctx, token, clients...)
}

func (lm *loggingMiddleware) ViewClient(ctx context.Context, token, id string) (c mgclients.Client, err error) {
	defer func(begin time.Time) {
		args := []interface{}{
			slog.String("duration", time.Since(begin).String()),
			slog.String("id", id),
		}
		if err != nil {
			args = append(args, slog.String("error", err.Error()))
			lm.logger.Warn("View client failed to complete successfully", args...)
			return
		}
		lm.logger.Info("View client completed successfully", args...)
	}(time.Now())
	return lm.svc.ViewClient(ctx, token, id)
}

func (lm *loggingMiddleware) ViewClientPerms(ctx context.Context, token, id string) (p []string, err error) {
	defer func(begin time.Time) {
		args := []interface{}{
			slog.String("duration", time.Since(begin).String()),
			slog.String("id", id),
		}
		if err != nil {
			args = append(args, slog.String("error", err.Error()))
			lm.logger.Warn("View client permissions failed to complete successfully", args...)
			return
		}
		lm.logger.Info("View client permissions completed successfully", args...)
	}(time.Now())
	return lm.svc.ViewClientPerms(ctx, token, id)
}

func (lm *loggingMiddleware) ListClients(ctx context.Context, token, reqUserID string, pm mgclients.Page) (cp mgclients.ClientsPage, err error) {
	defer func(begin time.Time) {
		args := []interface{}{
			slog.String("duration", time.Since(begin).String()),
			slog.String("user_id", reqUserID),
			slog.Group(
				"page",
				slog.Any("limit", pm.Limit),
				slog.Any("offset", pm.Offset),
			),
		}
		if err != nil {
			args = append(args, slog.String("error", err.Error()))
			lm.logger.Warn("List clients failed to complete successfully", args...)
			return
		}
		lm.logger.Info("List clients completed successfully", args...)
	}(time.Now())
	return lm.svc.ListClients(ctx, token, reqUserID, pm)
}

func (lm *loggingMiddleware) UpdateClient(ctx context.Context, token string, client mgclients.Client) (c mgclients.Client, err error) {
	defer func(begin time.Time) {
		args := []interface{}{
			slog.String("duration", time.Since(begin).String()),
			slog.Group(
				"client",
				slog.String("id", client.ID),
				slog.String("name", client.Name),
				slog.Any("metadata", client.Metadata),
			),
		}
		if err != nil {
			args = append(args, slog.String("error", err.Error()))
			lm.logger.Warn("Update client failed to complete successfully", args...)
			return
		}
		lm.logger.Info("Update client completed successfully", args...)
	}(time.Now())
	return lm.svc.UpdateClient(ctx, token, client)
}

func (lm *loggingMiddleware) UpdateClientTags(ctx context.Context, token string, client mgclients.Client) (c mgclients.Client, err error) {
	defer func(begin time.Time) {
		args := []interface{}{
			slog.String("duration", time.Since(begin).String()),
			slog.Group(
				"client",
				slog.String("id", client.ID),
				slog.Any("tags", client.Tags),
			),
		}
		if err != nil {
			args := append(args, slog.String("error", err.Error()))
			lm.logger.Warn("Update client tags failed to complete successfully", args...)
			return
		}
		lm.logger.Info("Update client tags completed successfully", args...)
	}(time.Now())
	return lm.svc.UpdateClientTags(ctx, token, client)
}

func (lm *loggingMiddleware) UpdateClientSecret(ctx context.Context, token, oldSecret, newSecret string) (c mgclients.Client, err error) {
	defer func(begin time.Time) {
		args := []interface{}{
			slog.String("duration", time.Since(begin).String()),
			slog.Group(
				"things",
				slog.String("id", c.ID),
			),
		}
		if err != nil {
			args = append(args, slog.String("error", err.Error()))
			lm.logger.Warn("Update client secret failed to complete successfully", args...)
			return
		}
		lm.logger.Info("Update client secret completed successfully", args...)
	}(time.Now())
	return lm.svc.UpdateClientSecret(ctx, token, oldSecret, newSecret)
}

func (lm *loggingMiddleware) EnableClient(ctx context.Context, token, id string) (c mgclients.Client, err error) {
	defer func(begin time.Time) {
		args := []interface{}{
			slog.String("duration", time.Since(begin).String()),
			slog.String("id", id),
		}
		if err != nil {
			args = append(args, slog.String("error", err.Error()))
			lm.logger.Warn("Enable client failed to complete successfully", args...)
			return
		}
		lm.logger.Info("Enable client completed successfully", args...)
	}(time.Now())
	return lm.svc.EnableClient(ctx, token, id)
}

func (lm *loggingMiddleware) DisableClient(ctx context.Context, token, id string) (c mgclients.Client, err error) {
	defer func(begin time.Time) {
		args := []interface{}{
			slog.String("duration", time.Since(begin).String()),
			slog.String("id", id),
		}
		if err != nil {
			args = append(args, slog.String("error", err.Error()))
			lm.logger.Warn("Disable client failed to complete successfully", args...)
			return
		}
		lm.logger.Info("Disable client completed successfully", args...)
	}(time.Now())
	return lm.svc.DisableClient(ctx, token, id)
}

func (lm *loggingMiddleware) ListClientsByGroup(ctx context.Context, token, channelID string, cp mgclients.Page) (mp mgclients.MembersPage, err error) {
	defer func(begin time.Time) {
		args := []interface{}{
			slog.String("duration", time.Since(begin).String()),
			slog.String("channel_id", channelID),
			slog.Group(
				"page",
				slog.Any("offset", cp.Offset),
				slog.Any("limit", cp.Limit),
			),
		}
		if err != nil {
			args = append(args, slog.String("error", err.Error()))
			lm.logger.Warn("List clients by group failed to complete successfully", args...)
			return
		}
		lm.logger.Info("List clients by group completed successfully", args...)
	}(time.Now())
	return lm.svc.ListClientsByGroup(ctx, token, channelID, cp)
}

func (lm *loggingMiddleware) Identify(ctx context.Context, key string) (id string, err error) {
	defer func(begin time.Time) {
		args := []interface{}{
			slog.String("duration", time.Since(begin).String()),
			slog.Group(
				"thing",
				slog.String("key", key),
				slog.String("id", id),
			),
		}
		if err != nil {
			args = append(args, slog.String("error", err.Error()))
			lm.logger.Warn("Identify failed to complete successfully", args...)
			return
		}
		lm.logger.Info("Identify completed successfully", args...)
	}(time.Now())
	return lm.svc.Identify(ctx, key)
}

func (lm *loggingMiddleware) Authorize(ctx context.Context, req *magistrala.AuthorizeReq) (id string, err error) {
	defer func(begin time.Time) {
		args := []interface{}{
			slog.String("duration", time.Since(begin).String()),
			slog.String("thing_key", req.Subject),
			slog.String("channel_id", req.Object),
		}
		if err != nil {
			args = append(args, slog.String("error", err.Error()))
			lm.logger.Warn("Authorize failed to complete successfully", args...)
			return
		}
		lm.logger.Info("Authorize completed successfully", args...)
	}(time.Now())
	return lm.svc.Authorize(ctx, req)
}

func (lm *loggingMiddleware) Share(ctx context.Context, token, id, relation string, userids ...string) (err error) {
	defer func(begin time.Time) {
		args := []interface{}{
			slog.String("duration", time.Since(begin).String()),
			slog.String("thing_id", id),
			slog.Any("user_ids", userids),
			slog.String("relation", relation),
		}
		if err != nil {
			args = append(args, slog.String("error", err.Error()))
			lm.logger.Warn("Share failed to complete successfully", args...)
			return
		}
		lm.logger.Info("Share completed successfully", args...)
	}(time.Now())
	return lm.svc.Share(ctx, token, id, relation, userids...)
}

func (lm *loggingMiddleware) Unshare(ctx context.Context, token, id, relation string, userids ...string) (err error) {
	defer func(begin time.Time) {
		args := []interface{}{
			slog.String("duration", time.Since(begin).String()),
			slog.String("thing_id", id),
			slog.Any("user_ids", userids),
			slog.String("relation", relation),
		}
		if err != nil {
			args = append(args, slog.String("error", err.Error()))
			lm.logger.Warn("Unshare failed to complete successfully", args...)
			return
		}
		lm.logger.Info("Unshare completed successfully", args...)
	}(time.Now())
	return lm.svc.Unshare(ctx, token, id, relation, userids...)
}

func (lm *loggingMiddleware) DeleteClient(ctx context.Context, token, id string) (err error) {
	defer func(begin time.Time) {
		args := []interface{}{
			slog.String("duration", time.Since(begin).String()),
			slog.String("thing_id", id),
		}
		if err != nil {
			args = append(args, slog.String("error", err.Error()))
			lm.logger.Warn("Delete client failed to complete successfully", args...)
			return
		}
		lm.logger.Info("Delete client completed successfully", args...)
	}(time.Now())
	return lm.svc.DeleteClient(ctx, token, id)
}
