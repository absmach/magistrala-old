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
		message := fmt.Sprintf("Method create_things %d things completed", len(cs))
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("mathod", "create_things"),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("mathod", "create_things"),
			slog.String("token", token),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.CreateThings(ctx, token, clients...)
}

func (lm *loggingMiddleware) ViewClient(ctx context.Context, token, id string) (c mgclients.Client, err error) {
	defer func(begin time.Time) {
		message := "Method view_thing completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("method", "view_thing"),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("method", "view_thing"),
			slog.String("thing_id", id),
			slog.String("token", token),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.ViewClient(ctx, token, id)
}

func (lm *loggingMiddleware) ViewClientPerms(ctx context.Context, token, id string) (p []string, err error) {
	defer func(begin time.Time) {
		message := "Method view_thing_permissions completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("method", "view_thing_permissions"),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("method", "view_thing_permissions"),
			slog.String("thing_id", id),
			slog.String("token", token),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.ViewClientPerms(ctx, token, id)
}

func (lm *loggingMiddleware) ListClients(ctx context.Context, token, reqUserID string, pm mgclients.Page) (cp mgclients.ClientsPage, err error) {
	defer func(begin time.Time) {
		message := "Method list_things completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("method", "list_things"),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("method", "list_things"),
			slog.String("token", token),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.ListClients(ctx, token, reqUserID, pm)
}

func (lm *loggingMiddleware) UpdateClient(ctx context.Context, token string, client mgclients.Client) (c mgclients.Client, err error) {
	defer func(begin time.Time) {
		message := "Method update_thing_name_and_metadata completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("method", "update_thing_name_and_metadata"),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("method", "update_thing_name_and_metadata"),
			slog.String("thing_id", c.ID),
			slog.String("token", token),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.UpdateClient(ctx, token, client)
}

func (lm *loggingMiddleware) UpdateClientTags(ctx context.Context, token string, client mgclients.Client) (c mgclients.Client, err error) {
	defer func(begin time.Time) {
		message := "Method update_thing_tags completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("method", "update_thing_tags"),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("method", "update_thing_tags"),
			slog.String("thing_id", c.ID),
			slog.String("token", token),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.UpdateClientTags(ctx, token, client)
}

func (lm *loggingMiddleware) UpdateClientSecret(ctx context.Context, token, oldSecret, newSecret string) (c mgclients.Client, err error) {
	defer func(begin time.Time) {
		message := "Method update_thing_secret completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error", message),
				slog.String("method", "update_thing_secret"),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("method", "update_thing_secret"),
			slog.String("thing_id", c.ID),
			slog.String("token", token),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.UpdateClientSecret(ctx, token, oldSecret, newSecret)
}

func (lm *loggingMiddleware) EnableClient(ctx context.Context, token, id string) (c mgclients.Client, err error) {
	defer func(begin time.Time) {
		message := "Method enable_thing completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with errors.", message),
				slog.String("method", "enable_thing"),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("method", "enable_thing"),
			slog.String("thing_id", id),
			slog.String("token", token),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.EnableClient(ctx, token, id)
}

func (lm *loggingMiddleware) DisableClient(ctx context.Context, token, id string) (c mgclients.Client, err error) {
	defer func(begin time.Time) {
		message := "Method disable_thing completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("method", "disable_thing"),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("method", "disable_thing"),
			slog.String("thing_id", id),
			slog.String("token", token),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.DisableClient(ctx, token, id)
}

func (lm *loggingMiddleware) ListClientsByGroup(ctx context.Context, token, channelID string, cp mgclients.Page) (mp mgclients.MembersPage, err error) {
	defer func(begin time.Time) {
		message := "Method list_things_by_channel completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("method", "list_things_by_channel"),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("method", "list_things_by_channel"),
			slog.String("channel_id", channelID),
			slog.String("token", token),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.ListClientsByGroup(ctx, token, channelID, cp)
}

func (lm *loggingMiddleware) Identify(ctx context.Context, key string) (id string, err error) {
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
			slog.String("thing_id", id),
			slog.String("thing_key", key),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.Identify(ctx, key)
}

func (lm *loggingMiddleware) Authorize(ctx context.Context, req *magistrala.AuthorizeReq) (id string, err error) {
	defer func(begin time.Time) {
		message := "Method authorise completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("method", "authorise"),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("method", "authorize"),
			slog.String("thing_key", req.Subject),
			slog.String("channel_id", req.Object),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.Authorize(ctx, req)
}

func (lm *loggingMiddleware) Share(ctx context.Context, token, id, relation string, userids ...string) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method share for users %v completex", userids)
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("method", "share"),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("method", "share"),
			slog.String("thing_id", id),
			slog.String("relation", relation),
			slog.String("token", token),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.Share(ctx, token, id, relation, userids...)
}

func (lm *loggingMiddleware) Unshare(ctx context.Context, token, id, relation string, userids ...string) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method unshare for users %v completed", userids)
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error: %s.", message, err),
				slog.String("method", "unshare"),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("method", "unshare"),
			slog.String("thing_id", id),
			slog.String("relation", relation),
			slog.String("token", token),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.Unshare(ctx, token, id, relation, userids...)
}

func (lm *loggingMiddleware) DeleteClient(ctx context.Context, token, id string) (err error) {
	defer func(begin time.Time) {
		message := "Method delete_client completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error: %s.", message, err),
				slog.String("method", "delete_client"),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("method", "delete_client"),
			slog.String("thing_id", id),
			slog.String("token", token),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.DeleteClient(ctx, token, id)
}
