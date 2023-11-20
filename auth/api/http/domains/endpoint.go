// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package domains

import (
	"context"

	"github.com/absmach/magistrala/auth"
	"github.com/absmach/magistrala/pkg/clients"
	"github.com/go-kit/kit/endpoint"
)

func createDomainEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createDomainReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		d := auth.Domain{
			Name:     req.Name,
			Metadata: req.Metadata,
			Tags:     req.Tags,
			Alias:    req.Alias,
		}
		domain, err := svc.CreateDomain(ctx, req.token, d)
		if err != nil {
			return nil, err
		}

		return createDomainRes{Data: domain}, nil
	}
}

func retrieveDomainEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(retrieveDomainRequest)
		if err := req.validate(); err != nil {
			return nil, err
		}

		domain, err := svc.RetrieveDomain(ctx, req.token, req.domainID)
		if err != nil {
			return nil, err
		}
		return retrieveDomainRes{Data: domain}, nil
	}
}

func updateDomainEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(updateDomainReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		var metadata clients.Metadata
		if req.Metadata != nil {
			metadata = *req.Metadata
		}
		d := auth.DomainReq{
			Name:     req.Name,
			Metadata: &metadata,
			Tags:     req.Tags,
			Alias:    req.Alias,
		}
		domain, err := svc.UpdateDomain(ctx, req.token, req.domainID, d)
		if err != nil {
			return nil, err
		}

		return updateDomainRes{Data: domain}, nil
	}
}

func listDomainsEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listDomainsReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		page := auth.Page{
			Offset:     req.offset,
			Limit:      req.limit,
			Name:       req.name,
			Metadata:   req.metadata,
			Order:      req.order,
			Dir:        req.dir,
			Tag:        req.tag,
			Permission: req.permission,
			Status:     req.status,
		}
		dp, err := svc.ListDomains(ctx, req.token, page)
		if err != nil {
			return nil, err
		}
		return listDomainsRes{Data: dp}, nil
	}
}

func enableDomainEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(enableDomainReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		enable := clients.EnabledStatus
		d := auth.DomainReq{
			Status: &enable,
		}
		if _, err := svc.ChangeDomainStatus(ctx, req.token, req.domainID, d); err != nil {
			return nil, err
		}
		return enableDomainRes{}, nil
	}
}

func disableDomainEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(disableDomainReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		disable := clients.DisabledStatus
		d := auth.DomainReq{
			Status: &disable,
		}
		if _, err := svc.ChangeDomainStatus(ctx, req.token, req.domainID, d); err != nil {
			return nil, err
		}
		return disableDomainRes{}, nil
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

func listUserDomainsEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listUserDomainsReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		page := auth.Page{
			Offset:     req.offset,
			Limit:      req.limit,
			Name:       req.name,
			Metadata:   req.metadata,
			Order:      req.order,
			Dir:        req.dir,
			Tag:        req.tag,
			Permission: req.permission,
			Status:     req.status,
		}
		dp, err := svc.ListUserDomains(ctx, req.token, req.userID, page)
		if err != nil {
			return nil, err
		}
		return listUserDomainsRes{Data: dp}, nil
	}
}
