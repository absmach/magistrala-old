// Copyright (c) Magistrala
// SPDX-License-Identifier: Apache-2.0

package auth

import (
	"github.com/absmach/magistrala"
	authgrpc "github.com/absmach/magistrala/auth/api/grpc"
	"github.com/absmach/magistrala/pkg/errors"
	thingsauth "github.com/absmach/magistrala/things/api/grpc"
	"github.com/caarlos0/env/v10"
)

var errGrpcConfig = errors.New("failed to load grpc configuration")

// Setup loads Auth gRPC configuration and creates new Auth gRPC client.
//
// Example:
//
//	client, handler, err := auth.Setup("MG_AUTH_GRPC_")
func Setup(envPrefix string) (magistrala.AuthServiceClient, Handler, error) {
	cfg := Config{}
	if err := env.ParseWithOptions(&cfg, env.Options{Prefix: envPrefix}); err != nil {
		return nil, nil, errors.Wrap(errGrpcConfig, err)
	}

	client, err := newClient(cfg)
	if err != nil {
		return nil, nil, err
	}

	return authgrpc.NewClient(client.Connection(), cfg.Timeout), client, nil
}

// Setup loads Authz gRPC configuration and creates new Authz gRPC client.
//
// Example:
//
//	client, handler, err := auth.SetupAuthz("MG_THINGS_AUTH_GRPC_")
func SetupAuthz(envPrefix string) (magistrala.AuthzServiceClient, Handler, error) {
	cfg := Config{}
	if err := env.ParseWithOptions(&cfg, env.Options{Prefix: envPrefix}); err != nil {
		return nil, nil, errors.Wrap(errGrpcConfig, err)
	}

	client, err := newClient(cfg)
	if err != nil {
		return nil, nil, err
	}

	return thingsauth.NewClient(client.Connection(), cfg.Timeout), client, nil
}
