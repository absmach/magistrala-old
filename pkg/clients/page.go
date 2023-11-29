// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package clients

// Page contains page metadata that helps navigation.
type Page struct {
	Total      uint64   `json:"total"`
	Offset     uint64   `json:"offset"`
	Limit      uint64   `json:"limit"`
	Name       string   `json:"name,omitempty"`
	Order      string   `json:"order,omitempty"`
	Dir        string   `json:"dir,omitempty"`
	Metadata   Metadata `json:"metadata,omitempty"`
	Owner      string   `json:"owner,omitempty"`
	Tag        string   `json:"tag,omitempty"`
	Permission string   `json:"permission,omitempty"`
	Status     Status   `json:"status,omitempty"`
	IDs        []string `json:"ids,omitempty"`
	Identity   string   `json:"identity,omitempty"`
	Role       *Role    `json:"-"`
	ListPerms  bool     `json:"-"`
}
