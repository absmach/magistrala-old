// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package auth_test

import (
	"fmt"
	"testing"

	"github.com/absmach/magistrala/auth"
	"github.com/absmach/magistrala/internal/apiutil"
	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	assert.Equal(t, auth.DisabledStatus.String(), "disabled")
	assert.Equal(t, auth.EnabledStatus.String(), "enabled")
	assert.Equal(t, auth.FreezeStatus.String(), "freezed")
	assert.Equal(t, auth.AllStatus.String(), "all")
	assert.Equal(t, auth.Status(5).String(), "unknown")
}

func TestToStatus(t *testing.T) {
	cases := []struct {
		status   string
		expected auth.Status
		err      error
	}{
		{"disabled", auth.DisabledStatus, nil},
		{"enabled", auth.EnabledStatus, nil},
		{"freezed", auth.FreezeStatus, nil},
		{"all", auth.AllStatus, nil},
		{"unknown", auth.EnabledStatus, apiutil.ErrInvalidStatus},
	}
	for _, c := range cases {
		status, err := auth.ToStatus(c.status)
		assert.Equal(t, status, c.expected, fmt.Sprintf("To %s failed with status: %v", c.status, status))
		assert.Equal(t, c.err, err, fmt.Sprintf("To %s failed with error: %v", c.status, err))
	}
}

func TestMarshalJSON(t *testing.T) {
	cases := []struct {
		status   auth.Status
		expected string
	}{
		{auth.DisabledStatus, `"disabled"`},
		{auth.EnabledStatus, `"enabled"`},
		{auth.FreezeStatus, `"freezed"`},
		{auth.AllStatus, `"all"`},
		{auth.Status(5), `"unknown"`},
	}
	for _, c := range cases {
		b, err := c.status.MarshalJSON()
		assert.Equal(t, c.expected, string(b), fmt.Sprintf("MarshalJSON failed with status: %v", c.status))
		assert.Nil(t, err, fmt.Sprintf("MarshalJSON failed with error: %v", err))
	}
}

func TestUnmarshalJSON(t *testing.T) {
	cases := []struct {
		status   string
		expected auth.Status
		err      error
	}{
		{`"disabled"`, auth.DisabledStatus, nil},
		{`"enabled"`, auth.EnabledStatus, nil},
		{`"freezed"`, auth.FreezeStatus, nil},
		{`"all"`, auth.AllStatus, nil},
		{`"unknown"`, auth.EnabledStatus, apiutil.ErrInvalidStatus},
	}
	for _, c := range cases {
		var status auth.Status
		err := status.UnmarshalJSON([]byte(c.status))
		assert.Equal(t, c.expected, status, fmt.Sprintf("UnmarshalJSON failed with status: %v", status))
		assert.Equal(t, c.err, err, fmt.Sprintf("UnmarshalJSON failed with error: %v", err))
	}
}
