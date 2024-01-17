// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/absmach/magistrala/pkg/groups"
)

var _ groups.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger *slog.Logger
	svc    groups.Service
}

// LoggingMiddleware adds logging facilities to the groups service.
func LoggingMiddleware(svc groups.Service, logger *slog.Logger) groups.Service {
	return &loggingMiddleware{logger, svc}
}

// CreateGroup logs the create_group request. It logs the group name, id and token and the time it took to complete the request.
// If the request fails, it logs the error.
func (lm *loggingMiddleware) CreateGroup(ctx context.Context, token, kind string, group groups.Group) (g groups.Group, err error) {
	defer func(begin time.Time) {
		message := "Method create_group %s completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("group_name", g.Name),
			slog.String("token", token),
			slog.String("kind", kind),
			slog.String("group_id", g.ID),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.CreateGroup(ctx, token, kind, group)
}

// UpdateGroup logs the update_group request. It logs the group name, id and token and the time it took to complete the request.
// If the request fails, it logs the error.
func (lm *loggingMiddleware) UpdateGroup(ctx context.Context, token string, group groups.Group) (g groups.Group, err error) {
	defer func(begin time.Time) {
		message := "Method update_group completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("group_name", g.Name),
			slog.String("group_id", g.ID),
			slog.String("token", token),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.UpdateGroup(ctx, token, group)
}

// ViewGroup logs the view_group request. It logs the group name, id and token and the time it took to complete the request.
// If the request fails, it logs the error.
func (lm *loggingMiddleware) ViewGroup(ctx context.Context, token, id string) (g groups.Group, err error) {
	defer func(begin time.Time) {
		message := "Method view_group completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("group_name", g.Name),
			slog.String("group_id", g.ID),
			slog.String("token", token),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.ViewGroup(ctx, token, id)
}

// ViewGroupPerms logs the view_group request. It logs the group name, id and token and the time it took to complete the request.
// If the request fails, it logs the error.
func (lm *loggingMiddleware) ViewGroupPerms(ctx context.Context, token, id string) (p []string, err error) {
	defer func(begin time.Time) {
		message := "Method view_group_perms completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("group_id", id),
			slog.String("token", token),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.ViewGroupPerms(ctx, token, id)
}

// ListGroups logs the list_groups request. It logs the token and the time it took to complete the request.
// If the request fails, it logs the error.
func (lm *loggingMiddleware) ListGroups(ctx context.Context, token, memberKind, memberID string, gp groups.Page) (cg groups.Page, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method list_groups %d groups completed", cg.Total)
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("token", token),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.ListGroups(ctx, token, memberKind, memberID, gp)
}

// EnableGroup logs the enable_group request. It logs the group name, id and token and the time it took to complete the request.
// If the request fails, it logs the error.
func (lm *loggingMiddleware) EnableGroup(ctx context.Context, token, id string) (g groups.Group, err error) {
	defer func(begin time.Time) {
		message := "Method enable_group completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("group_name", g.Name),
			slog.String("group_id", g.ID),
			slog.String("token", token),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.EnableGroup(ctx, token, id)
}

// DisableGroup logs the disable_group request. It logs the group name, id and token and the time it took to complete the request.
// If the request fails, it logs the error.
func (lm *loggingMiddleware) DisableGroup(ctx context.Context, token, id string) (g groups.Group, err error) {
	defer func(begin time.Time) {
		message := "Method disable_group completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("group_name", g.Name),
			slog.String("group_id", g.ID),
			slog.String("token", token),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.DisableGroup(ctx, token, id)
}

// ListMembers logs the list_members request. It logs the groupID and token and the time it took to complete the request.
// If the request fails, it logs the error.
func (lm *loggingMiddleware) ListMembers(ctx context.Context, token, groupID, permission, memberKind string) (mp groups.MembersPage, err error) {
	defer func(begin time.Time) {
		message := "Method list_members completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("group_id", groupID),
			slog.String("token", token),
			slog.String("permission", permission),
			slog.String("member_kind", memberKind),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.ListMembers(ctx, token, groupID, permission, memberKind)
}

func (lm *loggingMiddleware) Assign(ctx context.Context, token, groupID, relation, memberKind string, memberIDs ...string) (err error) {
	defer func(begin time.Time) {
		message := "Method assign completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("member_IDS", fmt.Sprintf("%v", memberIDs)),
			slog.String("group_id", groupID),
			slog.String("token", token),
			slog.String("relation", relation),
			slog.String("member_kind", memberKind),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.Assign(ctx, token, groupID, relation, memberKind, memberIDs...)
}

func (lm *loggingMiddleware) Unassign(ctx context.Context, token, groupID, relation, memberKind string, memberIDs ...string) (err error) {
	defer func(begin time.Time) {
		message := "Method unassign completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("member_IDS", fmt.Sprintf("%v", memberIDs)),
			slog.String("group_id", groupID),
			slog.String("token", token),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.Unassign(ctx, token, groupID, relation, memberKind, memberIDs...)
}

func (lm *loggingMiddleware) DeleteGroup(ctx context.Context, token, id string) (err error) {
	defer func(begin time.Time) {
		message := "Method delete_group completed"
		if err != nil {
			lm.logger.Warn(
				fmt.Sprintf("%s with error.", message),
				slog.String("error", err.Error()),
				slog.String("duration", time.Since(begin).String()),
			)
			return
		}
		lm.logger.Info(
			fmt.Sprintf("%s without errors.", message),
			slog.String("group_id", id),
			slog.String("token", token),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.DeleteGroup(ctx, token, id)
}
