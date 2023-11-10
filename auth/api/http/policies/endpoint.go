// Copyright (c) Magistrala
// SPDX-License-Identifier: Apache-2.0

package policies

import (
	"context"

	"github.com/absmach/magistrala/auth"
	"github.com/go-kit/kit/endpoint"
)

func createPolicyEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(policiesReq)
		if err := req.validate(); err != nil {
			return createPolicyRes{}, err
		}
		if err := svc.AddPolicies(ctx, []auth.PolicyReq{}); err != nil {
			return createPolicyRes{}, err
		}

		return createPolicyRes{created: true}, nil
	}
}

func deletePoliciesEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(policiesReq)
		if err := req.validate(); err != nil {
			return deletePoliciesRes{}, err
		}
		if err := svc.DeletePolicies(ctx, []auth.PolicyReq{}); err != nil {
			return deletePoliciesRes{}, err
		}

		return deletePoliciesRes{deleted: true}, nil
	}
}
