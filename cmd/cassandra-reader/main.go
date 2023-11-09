// Copyright (c) Magistrala
// SPDX-License-Identifier: Apache-2.0

// Package main contains cassandra-reader main function to start the cassandra-reader service.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/absmach/magistrala"
	"github.com/absmach/magistrala/internal"
	cassandraclient "github.com/absmach/magistrala/internal/clients/cassandra"
	"github.com/absmach/magistrala/internal/server"
	httpserver "github.com/absmach/magistrala/internal/server/http"
	mglog "github.com/absmach/magistrala/logger"
	"github.com/absmach/magistrala/pkg/auth"
	"github.com/absmach/magistrala/pkg/uuid"
	"github.com/absmach/magistrala/readers"
	"github.com/absmach/magistrala/readers/api"
	"github.com/absmach/magistrala/readers/cassandra"
	"github.com/caarlos0/env/v10"
	"github.com/gocql/gocql"
	chclient "github.com/mainflux/callhome/pkg/client"
	"golang.org/x/sync/errgroup"
)

const (
	svcName        = "cassandra-reader"
	envPrefixDB    = "MG_CASSANDRA_"
	envPrefixHTTP  = "MG_CASSANDRA_READER_HTTP_"
	envPrefixAuth  = "MG_AUTH_GRPC_"
	envPrefixAuthz = "MG_THINGS_AUTH_GRPC_"
	defSvcHTTPPort = "9003"
)

type config struct {
	LogLevel      string `env:"MG_CASSANDRA_READER_LOG_LEVEL"     envDefault:"info"`
	SendTelemetry bool   `env:"MG_SEND_TELEMETRY"                 envDefault:"true"`
	InstanceID    string `env:"MG_CASSANDRA_READER_INSTANCE_ID"   envDefault:""`
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	g, ctx := errgroup.WithContext(ctx)

	// Create cassandra reader service configurations
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("failed to load %s service configuration : %s", svcName, err)
	}

	logger, err := mglog.New(os.Stdout, cfg.LogLevel)
	if err != nil {
		log.Fatalf("failed to init logger: %s", err)
	}

	var exitCode int
	defer mglog.ExitWithError(&exitCode)

	if cfg.InstanceID == "" {
		if cfg.InstanceID, err = uuid.New().ID(); err != nil {
			logger.Error(fmt.Sprintf("failed to generate instanceID: %s", err))
			exitCode = 1
			return
		}
	}

	ac, acHandler, err := auth.Setup(envPrefixAuth)
	if err != nil {
		logger.Error(err.Error())
		exitCode = 1
		return
	}
	defer acHandler.Close()

	logger.Info("Successfully connected to auth grpc server " + acHandler.Secure())

	tc, tcHandler, err := auth.SetupAuthz(envPrefixAuthz)
	if err != nil {
		logger.Error(err.Error())
		exitCode = 1
		return
	}
	defer tcHandler.Close()

	logger.Info("Successfully connected to things grpc server " + tcHandler.Secure())

	// Create new cassandra client
	csdSession, err := cassandraclient.Setup(envPrefixDB)
	if err != nil {
		logger.Error(err.Error())
		exitCode = 1
		return
	}
	defer csdSession.Close()

	// Create new service
	repo := newService(csdSession, logger)

	// Create new http server
	httpServerConfig := server.Config{Port: defSvcHTTPPort}
	if err := env.ParseWithOptions(&httpServerConfig, env.Options{Prefix: envPrefixHTTP}); err != nil {
		logger.Error(fmt.Sprintf("failed to load %s HTTP server configuration : %s", svcName, err))
		exitCode = 1
		return
	}
	hs := httpserver.New(ctx, cancel, svcName, httpServerConfig, api.MakeHandler(repo, ac, tc, svcName, cfg.InstanceID), logger)

	if cfg.SendTelemetry {
		chc := chclient.New(svcName, magistrala.Version, logger, cancel)
		go chc.CallHome(ctx)
	}

	// Start servers
	g.Go(func() error {
		return hs.Start()
	})

	g.Go(func() error {
		return server.StopSignalHandler(ctx, cancel, logger, svcName, hs)
	})

	if err := g.Wait(); err != nil {
		logger.Error(fmt.Sprintf("Cassandra reader service terminated: %s", err))
	}
}

func newService(csdSession *gocql.Session, logger mglog.Logger) readers.MessageRepository {
	repo := cassandra.New(csdSession)
	repo = api.LoggingMiddleware(repo, logger)
	counter, latency := internal.MakeMetrics("cassandra", "message_reader")
	repo = api.MetricsMiddleware(repo, counter, latency)
	return repo
}
