// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

//go:build !test

package api

import (
	"context"
	"time"

	"github.com/absmach/magistrala/auth"
	"github.com/go-kit/kit/metrics"
)

var _ auth.Service = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	counter metrics.Counter
	latency metrics.Histogram
	svc     auth.Service
}

// MetricsMiddleware instruments core service by tracking request count and latency.
func MetricsMiddleware(svc auth.Service, counter metrics.Counter, latency metrics.Histogram) auth.Service {
	return &metricsMiddleware{
		counter: counter,
		latency: latency,
		svc:     svc,
	}
}

func (ms *metricsMiddleware) ListObjects(ctx context.Context, pr auth.PolicyReq, nextPageToken string, limit int32) (p auth.PolicyPage, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "list_objects").Add(1)
		ms.latency.With("method", "list_objects").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ListObjects(ctx, pr, nextPageToken, limit)
}

func (ms *metricsMiddleware) ListAllObjects(ctx context.Context, pr auth.PolicyReq) (p auth.PolicyPage, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "list_all_objects").Add(1)
		ms.latency.With("method", "list_all_objects").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ListAllObjects(ctx, pr)
}

func (ms *metricsMiddleware) CountObjects(ctx context.Context, pr auth.PolicyReq) (count int, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "count_objects").Add(1)
		ms.latency.With("method", "count_objects").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return ms.svc.CountObjects(ctx, pr)
}

func (ms *metricsMiddleware) ListSubjects(ctx context.Context, pr auth.PolicyReq, nextPageToken string, limit int32) (p auth.PolicyPage, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "list_subjects").Add(1)
		ms.latency.With("method", "list_subjects").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ListSubjects(ctx, pr, nextPageToken, limit)
}

func (ms *metricsMiddleware) ListAllSubjects(ctx context.Context, pr auth.PolicyReq) (p auth.PolicyPage, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "list_all_subjects").Add(1)
		ms.latency.With("method", "list_all_subjects").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ListAllSubjects(ctx, pr)
}

func (ms *metricsMiddleware) CountSubjects(ctx context.Context, pr auth.PolicyReq) (count int, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "count_subjects").Add(1)
		ms.latency.With("method", "count_subjects").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return ms.svc.CountSubjects(ctx, pr)
}

func (ms *metricsMiddleware) ListPermissions(ctx context.Context, pr auth.PolicyReq, filterPermissions []string) (p auth.Permissions, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "list_permissions").Add(1)
		ms.latency.With("method", "list_permissions").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ListPermissions(ctx, pr, filterPermissions)
}

func (ms *metricsMiddleware) Issue(ctx context.Context, token string, key auth.Key) (auth.Token, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "issue_key").Add(1)
		ms.latency.With("method", "issue_key").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.Issue(ctx, token, key)
}

func (ms *metricsMiddleware) Revoke(ctx context.Context, token, id string) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "revoke_key").Add(1)
		ms.latency.With("method", "revoke_key").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.Revoke(ctx, token, id)
}

func (ms *metricsMiddleware) RetrieveKey(ctx context.Context, token, id string) (auth.Key, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "retrieve_key").Add(1)
		ms.latency.With("method", "retrieve_key").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.RetrieveKey(ctx, token, id)
}

func (ms *metricsMiddleware) Identify(ctx context.Context, token string) (auth.Key, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "identify").Add(1)
		ms.latency.With("method", "identify").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.Identify(ctx, token)
}

func (ms *metricsMiddleware) Authorize(ctx context.Context, pr auth.PolicyReq) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "authorize").Add(1)
		ms.latency.With("method", "authorize").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return ms.svc.Authorize(ctx, pr)
}

func (ms *metricsMiddleware) AddPolicy(ctx context.Context, pr auth.PolicyReq) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "add_policy").Add(1)
		ms.latency.With("method", "add_policy").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return ms.svc.AddPolicy(ctx, pr)
}

func (ms *metricsMiddleware) AddPolicies(ctx context.Context, prs []auth.PolicyReq) (err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "create_policy_bulk").Add(1)
		ms.latency.With("method", "create_policy_bulk").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.AddPolicies(ctx, prs)
}

func (ms *metricsMiddleware) DeletePolicy(ctx context.Context, pr auth.PolicyReq) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "delete_policy").Add(1)
		ms.latency.With("method", "delete_policy").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return ms.svc.DeletePolicy(ctx, pr)
}

func (ms *metricsMiddleware) DeletePolicies(ctx context.Context, prs []auth.PolicyReq) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "delete_policies").Add(1)
		ms.latency.With("method", "delete_policies").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return ms.svc.DeletePolicies(ctx, prs)
}

func (ms *metricsMiddleware) CreateDomain(ctx context.Context, token string, d auth.Domain) (auth.Domain, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "create_domain").Add(1)
		ms.latency.With("method", "create_domain").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return ms.svc.CreateDomain(ctx, token, d)
}

func (ms *metricsMiddleware) RetrieveDomain(ctx context.Context, token, id string) (auth.Domain, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "retrieve_domain").Add(1)
		ms.latency.With("method", "retrieve_domain").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return ms.svc.RetrieveDomain(ctx, token, id)
}

func (ms *metricsMiddleware) RetrieveDomainPermissions(ctx context.Context, token string, id string) (auth.Permissions, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "retrieve_domain_permissions").Add(1)
		ms.latency.With("method", "retrieve_domain_permissions").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return ms.svc.RetrieveDomainPermissions(ctx, token, id)
}

func (ms *metricsMiddleware) UpdateDomain(ctx context.Context, token, id string, d auth.DomainReq) (auth.Domain, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "update_domain").Add(1)
		ms.latency.With("method", "update_domain").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return ms.svc.UpdateDomain(ctx, token, id, d)
}

func (ms *metricsMiddleware) ChangeDomainStatus(ctx context.Context, token, id string, d auth.DomainReq) (auth.Domain, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "change_domain_status").Add(1)
		ms.latency.With("method", "change_domain_status").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return ms.svc.ChangeDomainStatus(ctx, token, id, d)
}

func (ms *metricsMiddleware) ListDomains(ctx context.Context, token string, page auth.Page) (auth.DomainsPage, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "list_domains").Add(1)
		ms.latency.With("method", "list_domains").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return ms.svc.ListDomains(ctx, token, page)
}

func (ms *metricsMiddleware) AssignUsers(ctx context.Context, token, id string, userIds []string, relation string) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "assign_users").Add(1)
		ms.latency.With("method", "assign_users").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return ms.svc.AssignUsers(ctx, token, id, userIds, relation)
}

func (ms *metricsMiddleware) UnassignUsers(ctx context.Context, token, id string, userIds []string, relation string) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "unassign_users").Add(1)
		ms.latency.With("method", "unassign_users").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return ms.svc.UnassignUsers(ctx, token, id, userIds, relation)
}

func (ms *metricsMiddleware) ListUserDomains(ctx context.Context, token, userID string, page auth.Page) (auth.DomainsPage, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "list_user_domains").Add(1)
		ms.latency.With("method", "list_user_domains").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return ms.svc.ListUserDomains(ctx, token, userID, page)
}
