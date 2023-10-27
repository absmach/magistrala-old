// Copyright (c) Magistrala
// SPDX-License-Identifier: Apache-2.0

package domains

import (
	"encoding/json"
	"net/http"

	"github.com/mainflux/mainflux"
)

var (
	_ mainflux.Response = (*createDomainRes)(nil)
	_ mainflux.Response = (*viewDomainRes)(nil)
	_ mainflux.Response = (*assignUsersRes)(nil)
	_ mainflux.Response = (*unassignUsersRes)(nil)
)

type createDomainRes struct {
	Data interface{}
}

func (res createDomainRes) Code() int {
	return http.StatusOK
}
func (res createDomainRes) Headers() map[string]string {
	return map[string]string{}
}
func (res createDomainRes) Empty() bool {
	return false
}
func (res createDomainRes) MarshalJSON() ([]byte, error) {
	return json.Marshal(res.Data)
}

type viewDomainRes struct {
	Data interface{}
}

func (res viewDomainRes) Code() int {
	return http.StatusOK
}
func (res viewDomainRes) Headers() map[string]string {
	return map[string]string{}
}
func (res viewDomainRes) Empty() bool {
	return false
}
func (res viewDomainRes) MarshalJSON() ([]byte, error) {
	return json.Marshal(res.Data)
}

type updateDomainRes struct {
	Data interface{}
}

func (res updateDomainRes) Code() int {
	return http.StatusOK
}
func (res updateDomainRes) Headers() map[string]string {
	return map[string]string{}
}
func (res updateDomainRes) Empty() bool {
	return false
}
func (res updateDomainRes) MarshalJSON() ([]byte, error) {
	return json.Marshal(res.Data)
}

type assignUsersRes struct{}

func (res assignUsersRes) Code() int {
	return http.StatusCreated
}
func (res assignUsersRes) Headers() map[string]string {
	return map[string]string{}
}
func (res assignUsersRes) Empty() bool {
	return true
}

type unassignUsersRes struct{}

func (res unassignUsersRes) Code() int {
	return http.StatusNoContent
}
func (res unassignUsersRes) Headers() map[string]string {
	return map[string]string{}
}
func (res unassignUsersRes) Empty() bool {
	return true
}
