// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package grpc_test

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	svc, krepo = newService()
	startGRPCServer(svc, port)

	code := m.Run()

	os.Exit(code)
}
