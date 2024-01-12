// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"

	"github.com/absmach/magistrala/internal/apiutil"
	"github.com/absmach/magistrala/pkg/errors"
	"github.com/absmach/magistrala/pkg/groups"
	"github.com/go-kit/kit/endpoint"
)

func CreateGroupEndpoint(svc groups.Service, kind string) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createGroupReq)
		if err := req.validate(); err != nil {
			return createGroupRes{created: false}, errors.Wrap(apiutil.ErrValidation, err)
		}

		group, err := svc.CreateGroup(ctx, req.token, kind, req.Group)
		if err != nil {
			return createGroupRes{created: false}, err
		}

		return createGroupRes{created: true, Group: group}, nil
	}
}

func ViewGroupEndpoint(svc groups.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(groupReq)
		if err := req.validate(); err != nil {
			return viewGroupRes{}, errors.Wrap(apiutil.ErrValidation, err)
		}

		group, err := svc.ViewGroup(ctx, req.token, req.id)
		if err != nil {
			return viewGroupRes{}, err
		}

		return viewGroupRes{Group: group}, nil
	}
}

func ViewGroupPermsEndpoint(svc groups.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(groupPermsReq)
		if err := req.validate(); err != nil {
			return viewGroupPermsRes{}, errors.Wrap(apiutil.ErrValidation, err)
		}

		p, err := svc.ViewGroupPerms(ctx, req.token, req.id)
		if err != nil {
			return viewGroupPermsRes{}, err
		}

		return viewGroupPermsRes{Permissions: p}, nil
	}
}

func UpdateGroupEndpoint(svc groups.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(updateGroupReq)
		if err := req.validate(); err != nil {
			return updateGroupRes{}, errors.Wrap(apiutil.ErrValidation, err)
		}

		group := groups.Group{
			ID:          req.id,
			Name:        req.Name,
			Description: req.Description,
			Metadata:    req.Metadata,
		}

		group, err := svc.UpdateGroup(ctx, req.token, group)
		if err != nil {
			return updateGroupRes{}, err
		}

		return updateGroupRes{Group: group}, nil
	}
}

func EnableGroupEndpoint(svc groups.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(changeGroupStatusReq)
		if err := req.validate(); err != nil {
			return changeStatusRes{}, errors.Wrap(apiutil.ErrValidation, err)
		}
		group, err := svc.EnableGroup(ctx, req.token, req.id)
		if err != nil {
			return changeStatusRes{}, err
		}
		return changeStatusRes{Group: group}, nil
	}
}

func DisableGroupEndpoint(svc groups.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(changeGroupStatusReq)
		if err := req.validate(); err != nil {
			return changeStatusRes{}, errors.Wrap(apiutil.ErrValidation, err)
		}
		group, err := svc.DisableGroup(ctx, req.token, req.id)
		if err != nil {
			return changeStatusRes{}, err
		}
		return changeStatusRes{Group: group}, nil
	}
}

func ListGroupsEndpoint(svc groups.Service, memberKind string) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listGroupsReq)
		if memberKind != "" {
			req.memberKind = memberKind
		}
		if err := req.validate(); err != nil {
			return groupPageRes{}, errors.Wrap(apiutil.ErrValidation, err)
		}
		page, err := svc.ListGroups(ctx, req.token, req.memberKind, req.memberID, req.Page)
		if err != nil {
			return groupPageRes{}, err
		}

		if req.tree {
			return buildGroupsResponseTree(page), nil
		}

		return buildGroupsResponse(page), nil
	}
}

func ListMembersEndpoint(svc groups.Service, memberKind string) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listMembersReq)
		if memberKind != "" {
			req.memberKind = memberKind
		}
		if err := req.validate(); err != nil {
			return listMembersRes{}, errors.Wrap(apiutil.ErrValidation, err)
		}

		page, err := svc.ListMembers(ctx, req.token, req.groupID, req.permission, req.memberKind)
		if err != nil {
			return listMembersRes{}, err
		}

		return listMembersRes{
			pageRes: pageRes{
				Limit:  page.Limit,
				Offset: page.Offset,
				Total:  page.Total,
			},
			Members: page.Members,
		}, nil
	}
}

func AssignMembersEndpoint(svc groups.Service, relation, memberKind string) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(assignReq)
		if relation != "" {
			req.Relation = relation
		}
		if memberKind != "" {
			req.MemberKind = memberKind
		}
		if err := req.validate(); err != nil {
			return assignRes{}, errors.Wrap(apiutil.ErrValidation, err)
		}
		if err := svc.Assign(ctx, req.token, req.groupID, req.Relation, req.MemberKind, req.Members...); err != nil {
			return assignRes{}, err
		}
		return assignRes{assigned: true}, nil
	}
}

func UnassignMembersEndpoint(svc groups.Service, relation, memberKind string) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(unassignReq)
		if relation != "" {
			req.Relation = relation
		}
		if memberKind != "" {
			req.MemberKind = memberKind
		}
		if err := req.validate(); err != nil {
			return unassignRes{}, errors.Wrap(apiutil.ErrValidation, err)
		}

		if err := svc.Unassign(ctx, req.token, req.groupID, req.Relation, req.MemberKind, req.Members...); err != nil {
			return unassignRes{}, err
		}
		return unassignRes{unassigned: true}, nil
	}
}

func DeleteGroupEndpoint(svc groups.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(groupReq)
		if err := req.validate(); err != nil {
			return deleteGroupRes{}, errors.Wrap(apiutil.ErrValidation, err)
		}
		if err := svc.DeleteGroup(ctx, req.token, req.id); err != nil {
			return deleteGroupRes{}, err
		}
		return deleteGroupRes{deleted: true}, nil
	}
}

func buildGroupsResponseTree(page groups.Page) groupPageRes {
	groupsMap := map[string]*groups.Group{}
	// Parents' map keeps its array of children.
	parentsMap := map[string][]*groups.Group{}
	for i := range page.Groups {
		if _, ok := groupsMap[page.Groups[i].ID]; !ok {
			groupsMap[page.Groups[i].ID] = &page.Groups[i]
			parentsMap[page.Groups[i].ID] = make([]*groups.Group, 0)
		}
	}

	for _, group := range groupsMap {
		if children, ok := parentsMap[group.Parent]; ok {
			children = append(children, group)
			parentsMap[group.Parent] = children
		}
	}

	res := groupPageRes{
		pageRes: pageRes{
			Limit:  page.Limit,
			Offset: page.Offset,
			Total:  page.Total,
			Level:  page.Level,
		},
		Groups: []viewGroupRes{},
	}

	for _, group := range groupsMap {
		if children, ok := parentsMap[group.ID]; ok {
			group.Children = children
		}
	}

	for _, group := range groupsMap {
		view := toViewGroupRes(*group)
		if children, ok := parentsMap[group.Parent]; len(children) == 0 || !ok {
			res.Groups = append(res.Groups, view)
		}
	}

	return res
}

func toViewGroupRes(group groups.Group) viewGroupRes {
	view := viewGroupRes{
		Group: group,
	}
	return view
}

func buildGroupsResponse(gp groups.Page) groupPageRes {
	res := groupPageRes{
		pageRes: pageRes{
			Total: gp.Total,
			Level: gp.Level,
		},
		Groups: []viewGroupRes{},
	}

	for _, group := range gp.Groups {
		view := viewGroupRes{
			Group: group,
		}
		res.Groups = append(res.Groups, view)
	}

	return res
}
