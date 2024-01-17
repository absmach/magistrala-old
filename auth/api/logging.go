// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

//go:build !test

package api

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/absmach/magistrala/auth"
)

var _ auth.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger *slog.Logger
	svc    auth.Service
}

// LoggingMiddleware adds logging facilities to the core service.
func LoggingMiddleware(svc auth.Service, logger *slog.Logger) auth.Service {
	return &loggingMiddleware{logger, svc}
}

func (lm *loggingMiddleware) ListObjects(ctx context.Context, pr auth.PolicyReq, nextPageToken string, limit int32) (p auth.PolicyPage, err error) {
	defer func(begin time.Time) {
		message := "Method list_objects completed"
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
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.ListObjects(ctx, pr, nextPageToken, limit)
}

func (lm *loggingMiddleware) ListAllObjects(ctx context.Context, pr auth.PolicyReq) (p auth.PolicyPage, err error) {
	defer func(begin time.Time) {
		message := "Method list_all_objects completed"
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
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.ListAllObjects(ctx, pr)
}

func (lm *loggingMiddleware) CountObjects(ctx context.Context, pr auth.PolicyReq) (count int, err error) {
	defer func(begin time.Time) {
		message := "Method count_objects completed"
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
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.CountObjects(ctx, pr)
}

func (lm *loggingMiddleware) ListSubjects(ctx context.Context, pr auth.PolicyReq, nextPageToken string, limit int32) (p auth.PolicyPage, err error) {
	defer func(begin time.Time) {
		message := "Method list_subjects completed"
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
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.ListSubjects(ctx, pr, nextPageToken, limit)
}

func (lm *loggingMiddleware) ListAllSubjects(ctx context.Context, pr auth.PolicyReq) (p auth.PolicyPage, err error) {
	defer func(begin time.Time) {
		message := "Method list_all_subjects completed"
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
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.ListAllSubjects(ctx, pr)
}

func (lm *loggingMiddleware) CountSubjects(ctx context.Context, pr auth.PolicyReq) (count int, err error) {
	defer func(begin time.Time) {
		message := "Method count_subjects completed"
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
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.CountSubjects(ctx, pr)
}

func (lm *loggingMiddleware) ListPermissions(ctx context.Context, pr auth.PolicyReq, filterPermissions []string) (p auth.Permissions, err error) {
	defer func(begin time.Time) {
		message := "Method list_permissions completed"
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
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.ListPermissions(ctx, pr, filterPermissions)
}

func (lm *loggingMiddleware) Issue(ctx context.Context, token string, key auth.Key) (tkn auth.Token, err error) {
	defer func(begin time.Time) {
		d := ""
		if key.Type != auth.AccessKey && !key.ExpiresAt.IsZero() {
			d = fmt.Sprintf("with expiration date %v", key.ExpiresAt)
		}
		message := fmt.Sprintf("Method issue for key %s completed", d)
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
			slog.String("key_type", key.Type.String()),
			slog.String("key_id", key.ID),
			slog.String("token", token),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.Issue(ctx, token, key)
}

func (lm *loggingMiddleware) Revoke(ctx context.Context, token, id string) (err error) {
	defer func(begin time.Time) {
		message := "Method revoke completed"
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
			slog.String("key_id", id),
			slog.String("token", token),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.Revoke(ctx, token, id)
}

func (lm *loggingMiddleware) RetrieveKey(ctx context.Context, token, id string) (key auth.Key, err error) {
	defer func(begin time.Time) {
		message := "Method retrieve_key completed"
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
			slog.String("key_id", id),
			slog.String("token", token),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.RetrieveKey(ctx, token, id)
}

func (lm *loggingMiddleware) Identify(ctx context.Context, token string) (id auth.Key, err error) {
	defer func(begin time.Time) {
		message := "Method identify completed"
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

	return lm.svc.Identify(ctx, token)
}

func (lm *loggingMiddleware) Authorize(ctx context.Context, pr auth.PolicyReq) (err error) {
	defer func(begin time.Time) {
		message := "Method authorize completed"
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
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.Authorize(ctx, pr)
}

func (lm *loggingMiddleware) AddPolicy(ctx context.Context, pr auth.PolicyReq) (err error) {
	defer func(begin time.Time) {
		message := "Method add_policy completed"
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
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.AddPolicy(ctx, pr)
}

func (lm *loggingMiddleware) AddPolicies(ctx context.Context, prs []auth.PolicyReq) (err error) {
	defer func(begin time.Time) {
		message := "Method add_policies completed"
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
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())

	return lm.svc.AddPolicies(ctx, prs)
}

func (lm *loggingMiddleware) DeletePolicy(ctx context.Context, pr auth.PolicyReq) (err error) {
	defer func(begin time.Time) {
		message := "Method delete_policy completed"
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
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.DeletePolicy(ctx, pr)
}

func (lm *loggingMiddleware) DeletePolicies(ctx context.Context, prs []auth.PolicyReq) (err error) {
	defer func(begin time.Time) {
		message := "Method delete_policies completed"
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
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.DeletePolicies(ctx, prs)
}

func (lm *loggingMiddleware) CreateDomain(ctx context.Context, token string, d auth.Domain) (do auth.Domain, err error) {
	defer func(begin time.Time) {
		message := "Method create_domain completed"
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
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.CreateDomain(ctx, token, d)
}

func (lm *loggingMiddleware) RetrieveDomain(ctx context.Context, token, id string) (do auth.Domain, err error) {
	defer func(begin time.Time) {
		message := "Method retrieve_domain completed"
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
			slog.String("domain_id", id),
			slog.String("token", token),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.RetrieveDomain(ctx, token, id)
}

func (lm *loggingMiddleware) RetrieveDomainPermissions(ctx context.Context, token, id string) (permissions auth.Permissions, err error) {
	defer func(begin time.Time) {
		message := "Method retrieve_domain_permissions completed"
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
			slog.String("domain_id", id),
			slog.String("token", token),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.RetrieveDomainPermissions(ctx, token, id)
}

func (lm *loggingMiddleware) UpdateDomain(ctx context.Context, token, id string, d auth.DomainReq) (do auth.Domain, err error) {
	defer func(begin time.Time) {
		message := "Method update_domain completed"
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
			slog.String("domain_id", id),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.UpdateDomain(ctx, token, id, d)
}

func (lm *loggingMiddleware) ChangeDomainStatus(ctx context.Context, token, id string, d auth.DomainReq) (do auth.Domain, err error) {
	defer func(begin time.Time) {
		message := "Method change_domain_status completed"
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
			slog.String("domain_id", "id"),
			slog.String("token", token),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.ChangeDomainStatus(ctx, token, id, d)
}

func (lm *loggingMiddleware) ListDomains(ctx context.Context, token string, page auth.Page) (do auth.DomainsPage, err error) {
	defer func(begin time.Time) {
		message := "Method list_domains completed"
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
	return lm.svc.ListDomains(ctx, token, page)
}

func (lm *loggingMiddleware) AssignUsers(ctx context.Context, token, id string, userIds []string, relation string) (err error) {
	defer func(begin time.Time) {
		message := "Method assign_users completed"
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
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.AssignUsers(ctx, token, id, userIds, relation)
}

func (lm *loggingMiddleware) UnassignUsers(ctx context.Context, token, id string, userIds []string, relation string) (err error) {
	defer func(begin time.Time) {
		message := "Method unassign_users completed"
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
			slog.String("user_ids", fmt.Sprintf("%v", userIds)),
			slog.String("relation", relation),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.UnassignUsers(ctx, token, id, userIds, relation)
}

func (lm *loggingMiddleware) ListUserDomains(ctx context.Context, token, userID string, page auth.Page) (do auth.DomainsPage, err error) {
	defer func(begin time.Time) {
		message := "Method list_user_domains completed"
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
			slog.String("user_id", userID),
			slog.String("duration", time.Since(begin).String()),
		)
	}(time.Now())
	return lm.svc.ListUserDomains(ctx, token, userID, page)
}
