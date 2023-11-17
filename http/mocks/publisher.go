// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package mocks

import (
	"context"

	"github.com/absmach/magistrala/pkg/messaging"
)

type mockPublisher struct{}

// NewPublisher returns mock message publisher.
func NewPublisher() messaging.Publisher {
	return mockPublisher{}
}

func (pub mockPublisher) Publish(ctx context.Context, topic string, msg *messaging.Message) error {
	return nil
}

func (pub mockPublisher) Close() error {
	return nil
}
