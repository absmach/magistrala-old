// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package mocks

import "fmt"

// Since mocks will store data in map, and they need to resemble the real
// identifiers as much as possible, a key will be created as combination of
// owner and their id. This will allow searching by prefix or suffix.
func key(owner string, id string) string {
	return fmt.Sprintf("%s-%s", owner, id)
}
