// Copyright (c) Magistrala
// SPDX-License-Identifier: Apache-2.0

package domains

import (
	"context"

	"github.com/absmach/magistrala/auth"
	"github.com/go-kit/kit/endpoint"
)

func createDomainEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createDomainReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		d := auth.Domain{
			Name:        req.Name,
			Description: req.Description,
			Metadata:    req.Metadata,
			Tags:        req.Tags,
			Alias:       req.Alias,
		}
		domain, err := svc.CreateDomain(ctx, req.token, d)
		if err != nil {
			return nil, err
		}

		return createDomainRes{Data: domain}, nil
	}
}

func viewDomainEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(viewDomainRequest)
		if err := req.validate(); err != nil {
			return nil, err
		}

		domain, err := svc.ViewDomain(ctx, req.token, req.domainID)
		if err != nil {
			return nil, err
		}
		return viewDomainRes{Data: domain}, nil

	}
}

func updateDomainEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(updateDomainReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		d := auth.Domain{
			Name:        req.Name,
			Description: req.Description,
			Metadata:    req.Metadata,
			Tags:        req.Tags,
			Alias:       req.Alias,
		}
		domain, err := svc.UpdateDomain(ctx, req.token, req.domainID, d)
		if err != nil {
			return nil, err
		}

		return updateDomainRes{Data: domain}, nil
	}
}

func assignDomainUsersEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(assignUsersReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		if err := svc.AssignUsers(ctx, req.token, req.domainID, req.UserIDs, req.Relation); err != nil {
			return nil, err
		}
		return assignUsersRes{}, nil
	}
}

func unassignDomainUsersEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(unassignUsersReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		if err := svc.UnassignUsers(ctx, req.token, req.domainID, req.UserIDs, req.Relation); err != nil {
			return nil, err
		}
		return unassignUsersRes{}, nil
	}
}
