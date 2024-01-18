// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package grpc_test

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/absmach/magistrala"
	"github.com/absmach/magistrala/auth"
	grpcapi "github.com/absmach/magistrala/auth/api/grpc"
	"github.com/absmach/magistrala/internal/apiutil"
	"github.com/absmach/magistrala/pkg/errors"
	svcerr "github.com/absmach/magistrala/pkg/errors/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

const (
	port            = 8081
	secret          = "secret"
	email           = "test@example.com"
	id              = "testID"
	thingsType      = "things"
	usersType       = "users"
	description     = "Description"
	groupName       = "mgx"
	adminpermission = "admin"

	authoritiesObj  = "authorities"
	memberRelation  = "member"
	loginDuration   = 30 * time.Minute
	refreshDuration = 24 * time.Hour
	invalidDuration = 7 * 24 * time.Hour
	validToken      = "valid"
	inValidToken    = "invalid"
	validPolicy     = "valid"
	validID         = "d4ebb847-5d0e-4e46-bdd9-b6aceaaa3a22"
	domainID        = "d4ebb847-5d0e-4e46-bdd9-b6aceaaa3a22"
)

func startGRPCServer(svc auth.Service, port int) {
	listener, _ := net.Listen("tcp", fmt.Sprintf(":%d", port))
	server := grpc.NewServer()
	magistrala.RegisterAuthServiceServer(server, grpcapi.NewServer(svc))
	go func() {
		if err := server.Serve(listener); err != nil {
			panic(fmt.Sprintf("failed to serve: %s", err))
		}
	}()
}

func TestIssue(t *testing.T) {
	authAddr := fmt.Sprintf("localhost:%d", port)
	conn, _ := grpc.Dial(authAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	client := grpcapi.NewClient(conn, time.Second)

	cases := []struct {
		desc          string
		id            string
		email         string
		kind          auth.KeyType
		issueResponse auth.Token
		err           error
		code          codes.Code
	}{
		{
			desc:  "issue for user with valid token",
			id:    id,
			email: email,
			kind:  auth.AccessKey,
			issueResponse: auth.Token{
				AccessToken:  validToken,
				RefreshToken: validToken,
			},
			err:  nil,
			code: codes.OK,
		},
		{
			desc:  "issue recovery key",
			id:    id,
			email: email,
			kind:  auth.RecoveryKey,
			issueResponse: auth.Token{
				AccessToken:  validToken,
				RefreshToken: validToken,
			},
			err:  nil,
			code: codes.OK,
		},
		{
			desc:          "issue API key unauthenticated",
			id:            id,
			email:         email,
			kind:          auth.APIKey,
			issueResponse: auth.Token{},
			err:           errors.ErrAuthentication,
			code:          codes.Unauthenticated,
		},
		{
			desc:          "issue for invalid key type",
			id:            id,
			email:         email,
			kind:          32,
			issueResponse: auth.Token{},
			err:           errors.ErrMalformedEntity,
			code:          codes.InvalidArgument,
		},
		{
			desc:          "issue for user that does notexist",
			id:            "",
			email:         "",
			kind:          auth.APIKey,
			issueResponse: auth.Token{},
			err:           errors.ErrAuthentication,
			code:          codes.Unauthenticated,
		},
	}

	for _, tc := range cases {
		repoCall := svc.On("Issue", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tc.issueResponse, tc.err)
		_, err := client.Issue(context.Background(), &magistrala.IssueReq{UserId: tc.id, Type: uint32(tc.kind)})
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
		repoCall.Unset()
	}
}

func TestRefresh(t *testing.T) {
	authAddr := fmt.Sprintf("localhost:%d", port)
	conn, _ := grpc.Dial(authAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	client := grpcapi.NewClient(conn, time.Second)

	cases := []struct {
		desc          string
		token         string
		issueResponse auth.Token
		err           error
		code          codes.Code
	}{
		{
			desc:  "refresh token with valid token",
			token: validToken,
			issueResponse: auth.Token{
				AccessToken:  validToken,
				RefreshToken: validToken,
			},
			err:  nil,
			code: codes.OK,
		},
		{
			desc:          "refresh token with invalid token",
			token:         inValidToken,
			issueResponse: auth.Token{},
			err:           errors.ErrAuthentication,
			code:          codes.Unauthenticated,
		},
		{
			desc:          "refresh token with empty token",
			token:         "",
			issueResponse: auth.Token{},
			err:           apiutil.ErrMissingSecret,
			code:          codes.InvalidArgument,
		},
	}

	for _, tc := range cases {
		repoCall := svc.On("Issue", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tc.issueResponse, tc.err)
		_, err := client.Refresh(context.Background(), &magistrala.RefreshReq{RefreshToken: tc.token})
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
		repoCall.Unset()
	}
}

func TestIdentify(t *testing.T) {
	authAddr := fmt.Sprintf("localhost:%d", port)
	conn, _ := grpc.Dial(authAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	client := grpcapi.NewClient(conn, time.Second)

	cases := []struct {
		desc   string
		token  string
		idt    *magistrala.IdentityRes
		svcErr error
		err    error
		code   codes.Code
	}{
		{
			desc:  "identify user with valid user token",
			token: validToken,
			idt:   &magistrala.IdentityRes{Id: id, UserId: email, DomainId: domainID},
			err:   nil,
			code:  codes.OK,
		},
		{
			desc:   "identify user with invalid user token",
			token:  "invalid",
			idt:    &magistrala.IdentityRes{},
			svcErr: svcerr.ErrAuthentication,
			err:    svcerr.ErrAuthentication,
			code:   codes.Unauthenticated,
		},
		{
			desc:  "identify user with empty token",
			token: "",
			idt:   &magistrala.IdentityRes{},
			err:   apiutil.ErrBearerToken,
			code:  codes.Unauthenticated,
		},
	}

	for _, tc := range cases {
		repoCall := svc.On("Identify", mock.Anything, mock.Anything, mock.Anything).Return(auth.Key{Subject: id, User: email, Domain: domainID}, tc.svcErr)
		idt, err := client.Identify(context.Background(), &magistrala.IdentityReq{Token: tc.token})
		if idt != nil {
			assert.Equal(t, tc.idt, idt, fmt.Sprintf("%s: expected %v got %v", tc.desc, tc.idt, idt))
		}
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
		repoCall.Unset()
	}
}

func TestAuthorize(t *testing.T) {
	authAddr := fmt.Sprintf("localhost:%d", port)
	conn, _ := grpc.Dial(authAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	client := grpcapi.NewClient(conn, time.Second)

	cases := []struct {
		desc        string
		token       string
		subject     string
		subjecttype string
		object      string
		objecttype  string
		relation    string
		permission  string
		ar          *magistrala.AuthorizeRes
		err         error
		code        codes.Code
	}{
		{
			desc:        "authorize user with authorized token",
			token:       validToken,
			subject:     id,
			subjecttype: usersType,
			object:      authoritiesObj,
			objecttype:  usersType,
			relation:    memberRelation,
			permission:  adminpermission,
			ar:          &magistrala.AuthorizeRes{Authorized: true},
			err:         nil,
		},
		{
			desc:        "authorize user with unauthorized token",
			token:       inValidToken,
			subject:     id,
			subjecttype: usersType,
			object:      authoritiesObj,
			objecttype:  usersType,
			relation:    memberRelation,
			permission:  adminpermission,
			ar:          &magistrala.AuthorizeRes{Authorized: false},
			err:         svcerr.ErrAuthorization,
		},
		{
			desc:        "authorize user with empty subject",
			token:       validToken,
			subject:     "",
			subjecttype: usersType,
			object:      authoritiesObj,
			objecttype:  usersType,
			relation:    memberRelation,
			permission:  adminpermission,
			ar:          &magistrala.AuthorizeRes{Authorized: false},
			err:         apiutil.ErrMissingPolicySub,
		},
		{
			desc:        "authorize user with empty subject type",
			token:       validToken,
			subject:     id,
			subjecttype: "",
			object:      authoritiesObj,
			objecttype:  usersType,
			relation:    memberRelation,
			permission:  adminpermission,
			ar:          &magistrala.AuthorizeRes{Authorized: false},
			err:         apiutil.ErrMissingPolicySub,
		},
		{
			desc:        "authorize user with empty object",
			token:       validToken,
			subject:     id,
			subjecttype: usersType,
			object:      "",
			objecttype:  usersType,
			relation:    memberRelation,
			permission:  adminpermission,
			ar:          &magistrala.AuthorizeRes{Authorized: false},
			err:         apiutil.ErrMissingPolicyObj,
		},
		{
			desc:        "authorize user with empty object type",
			token:       validToken,
			subject:     id,
			subjecttype: usersType,
			object:      authoritiesObj,
			objecttype:  "",
			relation:    memberRelation,
			permission:  adminpermission,
			ar:          &magistrala.AuthorizeRes{Authorized: false},
			err:         apiutil.ErrMissingPolicyObj,
		},
		{
			desc:        "authorize user with empty permission",
			token:       validToken,
			subject:     id,
			subjecttype: usersType,
			object:      authoritiesObj,
			objecttype:  usersType,
			relation:    memberRelation,
			permission:  "",
			ar:          &magistrala.AuthorizeRes{Authorized: false},
			err:         apiutil.ErrMalformedPolicyPer,
		},
	}
	for _, tc := range cases {
		repocall := svc.On("Authorize", mock.Anything, mock.Anything).Return(tc.err)
		ar, err := client.Authorize(context.Background(), &magistrala.AuthorizeReq{Subject: tc.subject, SubjectType: tc.subjecttype, Object: tc.object, ObjectType: tc.objecttype, Relation: tc.relation, Permission: tc.permission})
		if ar != nil {
			assert.Equal(t, tc.ar, ar, fmt.Sprintf("%s: expected %v got %v", tc.desc, tc.ar, ar))
		}
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
		repocall.Unset()
	}
}

func TestAddPolicy(t *testing.T) {
	authAddr := fmt.Sprintf("localhost:%d", port)
	conn, _ := grpc.Dial(authAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	client := grpcapi.NewClient(conn, time.Second)

	groupAdminObj := "groupadmin"

	cases := []struct {
		desc        string
		token       string
		subject     string
		subjecttype string
		object      string
		objecttype  string
		relation    string
		permission  string
		ar          *magistrala.AddPolicyRes
		err         error
		code        codes.Code
	}{
		{
			desc:        "add groupadmin policy to user",
			token:       validToken,
			subject:     id,
			subjecttype: usersType,
			object:      groupAdminObj,
			objecttype:  usersType,
			relation:    memberRelation,
			permission:  adminpermission,
			ar:          &magistrala.AddPolicyRes{Authorized: true},
			err:         nil,
		},
		{
			desc:        "add groupadmin policy to user with invalid token",
			token:       inValidToken,
			subject:     id,
			subjecttype: usersType,
			object:      groupAdminObj,
			objecttype:  usersType,
			relation:    memberRelation,
			permission:  adminpermission,
			ar:          &magistrala.AddPolicyRes{Authorized: false},
			err:         svcerr.ErrAuthorization,
		},
	}
	for _, tc := range cases {
		repoCall := svc.On("AddPolicy", mock.Anything, mock.Anything).Return(tc.err)
		apr, err := client.AddPolicy(context.Background(), &magistrala.AddPolicyReq{Subject: tc.subject, SubjectType: tc.subjecttype, Object: tc.object, ObjectType: tc.objecttype, Relation: tc.relation, Permission: tc.permission})
		if apr != nil {
			assert.Equal(t, tc.ar, apr, fmt.Sprintf("%s: expected %v got %v", tc.desc, tc.ar, apr))
		}
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
		repoCall.Unset()
	}
}

func TestAddPolicies(t *testing.T) {
	authAddr := fmt.Sprintf("localhost:%d", port)
	conn, _ := grpc.Dial(authAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	client := grpcapi.NewClient(conn, time.Second)

	groupAdminObj := "groupadmin"

	cases := []struct {
		desc  string
		token string
		pr    *magistrala.AddPoliciesReq
		ar    *magistrala.AddPoliciesRes
		err   error
		code  codes.Code
	}{
		{
			desc:  "add groupadmin policy to user",
			token: validToken,
			pr: &magistrala.AddPoliciesReq{
				AddPoliciesReq: []*magistrala.AddPolicyReq{
					{
						Subject:     id,
						SubjectType: usersType,
						Object:      groupAdminObj,
						ObjectType:  usersType,
						Relation:    memberRelation,
						Permission:  adminpermission,
					},
				},
			},
			ar:  &magistrala.AddPoliciesRes{Authorized: true},
			err: nil,
		},
		{
			desc:  "add groupadmin policy to user with invalid token",
			token: inValidToken,
			pr: &magistrala.AddPoliciesReq{
				AddPoliciesReq: []*magistrala.AddPolicyReq{
					{
						Subject:     id,
						SubjectType: usersType,
						Object:      groupAdminObj,
						ObjectType:  usersType,
						Relation:    memberRelation,
						Permission:  adminpermission,
					},
				},
			},
			ar:  &magistrala.AddPoliciesRes{Authorized: false},
			err: svcerr.ErrAuthorization,
		},
	}
	for _, tc := range cases {
		repoCall := svc.On("AddPolicies", mock.Anything, mock.Anything).Return(tc.err)
		apr, err := client.AddPolicies(context.Background(), tc.pr)
		if apr != nil {
			assert.Equal(t, tc.ar, apr, fmt.Sprintf("%s: expected %v got %v", tc.desc, tc.ar, apr))
		}
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
		repoCall.Unset()
	}
}

func TestDeletePolicy(t *testing.T) {
	authAddr := fmt.Sprintf("localhost:%d", port)
	conn, _ := grpc.Dial(authAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	client := grpcapi.NewClient(conn, time.Second)

	readRelation := "read"
	thingID := "thing"

	cases := []struct {
		desc        string
		token       string
		subject     string
		subjecttype string
		object      string
		objecttype  string
		relation    string
		permission  string
		dpr         *magistrala.DeletePolicyRes
		err         error
	}{
		{
			desc:        "delete valid policy",
			token:       validToken,
			subject:     id,
			subjecttype: usersType,
			object:      thingID,
			objecttype:  thingsType,
			relation:    readRelation,
			permission:  readRelation,
			dpr:         &magistrala.DeletePolicyRes{Deleted: true},
			err:         nil,
		},
		{
			desc:        "delete invalid policy with invalid token",
			token:       inValidToken,
			subject:     id,
			subjecttype: usersType,
			object:      thingID,
			objecttype:  thingsType,
			relation:    readRelation,
			permission:  readRelation,
			dpr:         &magistrala.DeletePolicyRes{Deleted: false},
			err:         svcerr.ErrAuthorization,
		},
	}
	for _, tc := range cases {
		repoCall := svc.On("DeletePolicy", mock.Anything, mock.Anything).Return(tc.err)
		dpr, err := client.DeletePolicy(context.Background(), &magistrala.DeletePolicyReq{Subject: tc.subject, SubjectType: tc.subjecttype, Object: tc.object, ObjectType: tc.objecttype, Relation: tc.relation})
		assert.Equal(t, tc.dpr.GetDeleted(), dpr.GetDeleted(), fmt.Sprintf("%s: expected %v got %v", tc.desc, tc.dpr.GetDeleted(), dpr.GetDeleted()))
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
		repoCall.Unset()
	}
}

func TestDeletePolicies(t *testing.T) {
	authAddr := fmt.Sprintf("localhost:%d", port)
	conn, _ := grpc.Dial(authAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	client := grpcapi.NewClient(conn, time.Second)

	readRelation := "read"
	thingID := "thing"

	cases := []struct {
		desc  string
		token string
		pr    *magistrala.DeletePoliciesReq
		ar    *magistrala.DeletePoliciesRes
		err   error
		code  codes.Code
	}{
		{
			desc:  "delete policies with valid token",
			token: validToken,
			pr: &magistrala.DeletePoliciesReq{
				DeletePoliciesReq: []*magistrala.DeletePolicyReq{
					{
						Subject:     id,
						SubjectType: usersType,
						Object:      thingID,
						ObjectType:  thingsType,
						Relation:    readRelation,
						Permission:  readRelation,
					},
				},
			},
			ar:  &magistrala.DeletePoliciesRes{Deleted: true},
			err: nil,
		},
		{
			desc:  "delete policies with invalid token",
			token: inValidToken,
			pr: &magistrala.DeletePoliciesReq{
				DeletePoliciesReq: []*magistrala.DeletePolicyReq{
					{
						Subject:     id,
						SubjectType: usersType,
						Object:      thingID,
						ObjectType:  thingsType,
						Relation:    readRelation,
						Permission:  readRelation,
					},
				},
			},
			ar:  &magistrala.DeletePoliciesRes{Deleted: false},
			err: svcerr.ErrAuthorization,
		},
	}
	for _, tc := range cases {
		repoCall := svc.On("DeletePolicies", mock.Anything, mock.Anything).Return(tc.err)
		apr, err := client.DeletePolicies(context.Background(), tc.pr)
		if apr != nil {
			assert.Equal(t, tc.ar, apr, fmt.Sprintf("%s: expected %v got %v", tc.desc, tc.ar, apr))
		}
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
		repoCall.Unset()
	}
}

func TestListObjects(t *testing.T) {
	authAddr := fmt.Sprintf("localhost:%d", port)
	conn, _ := grpc.Dial(authAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	client := grpcapi.NewClient(conn, time.Second)

	cases := []struct {
		desc  string
		token string
		pr    *magistrala.ListObjectsReq
		ar    *magistrala.ListObjectsRes
		err   error
		code  codes.Code
	}{
		{
			desc:  "list objects with valid token",
			token: validToken,
			pr: &magistrala.ListObjectsReq{
				Domain:     domainID,
				ObjectType: thingsType,
				Relation:   memberRelation,
				Permission: adminpermission,
			},
			ar: &magistrala.ListObjectsRes{
				Policies: []string{validPolicy},
			},
			err: nil,
		},
		{
			desc:  "list objects with invalid token",
			token: inValidToken,
			pr: &magistrala.ListObjectsReq{
				Domain:     domainID,
				ObjectType: thingsType,
				Relation:   memberRelation,
				Permission: adminpermission,
			},
			ar:  &magistrala.ListObjectsRes{},
			err: svcerr.ErrAuthorization,
		},
	}
	for _, tc := range cases {
		repoCall := svc.On("ListObjects", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(auth.PolicyPage{Policies: tc.ar.Policies}, tc.err)
		apr, err := client.ListObjects(context.Background(), tc.pr)
		if apr != nil {
			assert.Equal(t, tc.ar, apr, fmt.Sprintf("%s: expected %v got %v", tc.desc, tc.ar, apr))
		}
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
		repoCall.Unset()
	}
}

func TestListAllObjects(t *testing.T) {
	authAddr := fmt.Sprintf("localhost:%d", port)
	conn, _ := grpc.Dial(authAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	client := grpcapi.NewClient(conn, time.Second)
	cases := []struct {
		desc  string
		token string
		pr    *magistrala.ListObjectsReq
		ar    *magistrala.ListObjectsRes
		err   error
		code  codes.Code
	}{
		{
			desc:  "list all objects with valid token",
			token: validToken,
			pr: &magistrala.ListObjectsReq{
				Domain:     domainID,
				ObjectType: thingsType,
				Relation:   memberRelation,
				Permission: adminpermission,
			},
			ar: &magistrala.ListObjectsRes{
				Policies: []string{validPolicy},
			},
			err: nil,
		},
		{
			desc:  "list all objects with invalid token",
			token: inValidToken,
			pr: &magistrala.ListObjectsReq{
				Domain:     domainID,
				ObjectType: thingsType,
				Relation:   memberRelation,
				Permission: adminpermission,
			},
			ar:  &magistrala.ListObjectsRes{},
			err: svcerr.ErrAuthorization,
		},
	}
	for _, tc := range cases {
		repoCall := svc.On("ListAllObjects", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(auth.PolicyPage{Policies: tc.ar.Policies}, tc.err)
		apr, err := client.ListAllObjects(context.Background(), tc.pr)
		if apr != nil {
			assert.Equal(t, tc.ar, apr, fmt.Sprintf("%s: expected %v got %v", tc.desc, tc.ar, apr))
		}
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
		repoCall.Unset()
	}
}

func TestCountObects(t *testing.T) {
	authAddr := fmt.Sprintf("localhost:%d", port)
	conn, _ := grpc.Dial(authAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	client := grpcapi.NewClient(conn, time.Second)
	cases := []struct {
		desc  string
		token string
		pr    *magistrala.CountObjectsReq
		ar    *magistrala.CountObjectsRes
		err   error
		code  codes.Code
	}{
		{
			desc:  "count objects with valid token",
			token: validToken,
			pr: &magistrala.CountObjectsReq{
				Domain:     domainID,
				ObjectType: thingsType,
				Relation:   memberRelation,
				Permission: adminpermission,
			},
			ar: &magistrala.CountObjectsRes{
				Count: 1,
			},
			err: nil,
		},
		{
			desc:  "count objects with invalid token",
			token: inValidToken,
			pr: &magistrala.CountObjectsReq{
				Domain:     domainID,
				ObjectType: thingsType,
				Relation:   memberRelation,
				Permission: adminpermission,
			},
			ar:  &magistrala.CountObjectsRes{},
			err: svcerr.ErrAuthorization,
		},
	}
	for _, tc := range cases {
		repoCall := svc.On("CountObjects", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(int(tc.ar.Count), tc.err)
		apr, err := client.CountObjects(context.Background(), tc.pr)
		if apr != nil {
			assert.Equal(t, tc.ar, apr, fmt.Sprintf("%s: expected %v got %v", tc.desc, tc.ar, apr))
		}
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
		repoCall.Unset()
	}
}

func TestListSubjects(t *testing.T) {
	authAddr := fmt.Sprintf("localhost:%d", port)
	conn, _ := grpc.Dial(authAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	client := grpcapi.NewClient(conn, time.Second)
	cases := []struct {
		desc  string
		token string
		pr    *magistrala.ListSubjectsReq
		ar    *magistrala.ListSubjectsRes
		err   error
		code  codes.Code
	}{
		{
			desc:  "list subjects with valid token",
			token: validToken,
			pr: &magistrala.ListSubjectsReq{
				Domain:      domainID,
				SubjectType: usersType,
				Relation:    memberRelation,
				Permission:  adminpermission,
			},
			ar: &magistrala.ListSubjectsRes{
				Policies: []string{validPolicy},
			},
			err: nil,
		},
		{
			desc:  "list subjects with invalid token",
			token: inValidToken,
			pr: &magistrala.ListSubjectsReq{
				Domain:      domainID,
				SubjectType: usersType,
				Relation:    memberRelation,
				Permission:  adminpermission,
			},
			ar:  &magistrala.ListSubjectsRes{},
			err: svcerr.ErrAuthorization,
		},
	}
	for _, tc := range cases {
		repoCall := svc.On("ListSubjects", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(auth.PolicyPage{Policies: tc.ar.Policies}, tc.err)
		apr, err := client.ListSubjects(context.Background(), tc.pr)
		if apr != nil {
			assert.Equal(t, tc.ar, apr, fmt.Sprintf("%s: expected %v got %v", tc.desc, tc.ar, apr))
		}
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
		repoCall.Unset()
	}
}

func TestListAllSubjects(t *testing.T) {
	authAddr := fmt.Sprintf("localhost:%d", port)
	conn, _ := grpc.Dial(authAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	client := grpcapi.NewClient(conn, time.Second)
	cases := []struct {
		desc  string
		token string
		pr    *magistrala.ListSubjectsReq
		ar    *magistrala.ListSubjectsRes
		err   error
		code  codes.Code
	}{
		{
			desc:  "list all subjects with valid token",
			token: validToken,
			pr: &magistrala.ListSubjectsReq{
				Domain:      domainID,
				SubjectType: auth.UserType,
				Relation:    memberRelation,
				Permission:  adminpermission,
			},
			ar: &magistrala.ListSubjectsRes{
				Policies: []string{validPolicy},
			},
			err: nil,
		},
		{
			desc:  "list all subjects with invalid token",
			token: inValidToken,
			pr: &magistrala.ListSubjectsReq{
				Domain:      domainID,
				SubjectType: usersType,
				Relation:    memberRelation,
				Permission:  adminpermission,
			},
			ar:  &magistrala.ListSubjectsRes{},
			err: svcerr.ErrAuthorization,
		},
	}
	for _, tc := range cases {
		repoCall := svc.On("ListAllSubjects", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(auth.PolicyPage{Policies: tc.ar.Policies}, tc.err)
		apr, err := client.ListAllSubjects(context.Background(), tc.pr)
		if apr != nil {
			assert.Equal(t, tc.ar, apr, fmt.Sprintf("%s: expected %v got %v", tc.desc, tc.ar, apr))
		}
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
		repoCall.Unset()
	}
}

func TestCountSubjects(t *testing.T) {
	authAddr := fmt.Sprintf("localhost:%d", port)
	conn, _ := grpc.Dial(authAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	client := grpcapi.NewClient(conn, time.Second)
	cases := []struct {
		desc  string
		token string
		pr    *magistrala.CountSubjectsReq
		ar    *magistrala.CountSubjectsRes
		err   error
		code  codes.Code
	}{
		{
			desc:  "count subjects with valid token",
			token: validToken,
			pr: &magistrala.CountSubjectsReq{
				Domain:      domainID,
				SubjectType: usersType,
				Relation:    memberRelation,
				Permission:  adminpermission,
			},
			ar: &magistrala.CountSubjectsRes{
				Count: 1,
			},
			code: codes.OK,
			err:  nil,
		},
		{
			desc:  "count subjects with invalid token",
			token: inValidToken,
			pr: &magistrala.CountSubjectsReq{
				Domain:      domainID,
				SubjectType: usersType,
				Relation:    memberRelation,
				Permission:  adminpermission,
			},
			ar:   &magistrala.CountSubjectsRes{},
			err:  svcerr.ErrAuthentication,
			code: codes.Unauthenticated,
		},
	}
	for _, tc := range cases {
		repoCall := svc.On("CountSubjects", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(int(tc.ar.Count), tc.err)
		apr, err := client.CountSubjects(context.Background(), tc.pr)
		if apr != nil {
			assert.Equal(t, tc.ar, apr, fmt.Sprintf("%s: expected %v got %v", tc.desc, tc.ar, apr))
		}
		e, ok := status.FromError(err)
		assert.True(t, ok, "gRPC status can't be extracted from the error")
		assert.Equal(t, tc.code, e.Code(), fmt.Sprintf("%s: expected %s got %s", tc.desc, tc.code, e.Code()))
		repoCall.Unset()
	}
}

func TestListPermissions(t *testing.T) {
	authAddr := fmt.Sprintf("localhost:%d", port)
	conn, _ := grpc.Dial(authAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	client := grpcapi.NewClient(conn, time.Second)
	cases := []struct {
		desc  string
		token string
		pr    *magistrala.ListPermissionsReq
		ar    *magistrala.ListPermissionsRes
		err   error
		code  codes.Code
	}{
		{
			desc:  "list permissions of thing type with valid token",
			token: validToken,
			pr: &magistrala.ListPermissionsReq{
				Domain:            domainID,
				SubjectType:       auth.UserType,
				Subject:           id,
				ObjectType:        auth.ThingType,
				Object:            validID,
				FilterPermissions: []string{"view"},
			},
			ar: &magistrala.ListPermissionsRes{
				SubjectType: auth.UserType,
				Subject:     id,
				ObjectType:  auth.ThingType,
				Object:      validID,
				Permissions: []string{"view"},
			},
			err: nil,
		},
		{
			desc:  "list permissions of thing type with valid token",
			token: validToken,
			pr: &magistrala.ListPermissionsReq{
				Domain:            domainID,
				SubjectType:       auth.UserType,
				Subject:           id,
				ObjectType:        auth.GroupType,
				Object:            validID,
				FilterPermissions: []string{"view"},
			},
			ar: &magistrala.ListPermissionsRes{
				SubjectType: auth.UserType,
				Subject:     id,
				ObjectType:  auth.GroupType,
				Object:      validID,
				Permissions: []string{"view"},
			},
			err: nil,
		},
		{
			desc:  "list permissions of platform type with valid token",
			token: validToken,
			pr: &magistrala.ListPermissionsReq{
				Domain:            domainID,
				SubjectType:       auth.UserType,
				Subject:           id,
				ObjectType:        auth.PlatformType,
				Object:            validID,
				FilterPermissions: []string{"view"},
			},
			ar: &magistrala.ListPermissionsRes{
				SubjectType: auth.UserType,
				Subject:     id,
				ObjectType:  auth.PlatformType,
				Object:      validID,
				Permissions: []string{"view"},
			},
			err: nil,
		},
		{
			desc:  "list permissions of thing type with valid token",
			token: validToken,
			pr: &magistrala.ListPermissionsReq{
				Domain:            domainID,
				SubjectType:       auth.UserType,
				Subject:           id,
				ObjectType:        auth.DomainType,
				Object:            validID,
				FilterPermissions: []string{"view"},
			},
			ar: &magistrala.ListPermissionsRes{
				SubjectType: auth.UserType,
				Subject:     id,
				ObjectType:  auth.DomainType,
				Object:      validID,
				Permissions: []string{"view"},
			},
			err: nil,
		},
		{
			desc:  "list permissions of thing type with invalid token",
			token: inValidToken,
			pr: &magistrala.ListPermissionsReq{
				Domain:            domainID,
				SubjectType:       auth.UserType,
				Subject:           id,
				ObjectType:        auth.ThingType,
				Object:            validID,
				FilterPermissions: []string{"view"},
			},
			ar:  &magistrala.ListPermissionsRes{},
			err: svcerr.ErrAuthentication,
		},
		{
			desc:  "list permissions with invalid object type",
			token: validToken,
			pr: &magistrala.ListPermissionsReq{
				Domain:            domainID,
				SubjectType:       auth.UserType,
				Subject:           id,
				ObjectType:        "invalid",
				Object:            validID,
				FilterPermissions: []string{"view"},
			},
			ar:  &magistrala.ListPermissionsRes{},
			err: apiutil.ErrMalformedPolicy,
		},
	}
	for _, tc := range cases {
		repoCall := svc.On("ListPermissions", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(auth.Permissions{"view"}, tc.err)
		apr, err := client.ListPermissions(context.Background(), tc.pr)
		if apr != nil {
			assert.Equal(t, tc.ar, apr, fmt.Sprintf("%s: expected %v got %v", tc.desc, tc.ar, apr))
		}
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
		repoCall.Unset()
	}
}
