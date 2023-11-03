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
	_ mainflux.Response = (*retrieveDomainRes)(nil)
	_ mainflux.Response = (*assignUsersRes)(nil)
	_ mainflux.Response = (*unassignUsersRes)(nil)
	_ mainflux.Response = (*listDomainsRes)(nil)
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

type retrieveDomainRes struct {
	Data interface{}
}

func (res retrieveDomainRes) Code() int {
	return http.StatusOK
}
func (res retrieveDomainRes) Headers() map[string]string {
	return map[string]string{}
}
func (res retrieveDomainRes) Empty() bool {
	return false
}
func (res retrieveDomainRes) MarshalJSON() ([]byte, error) {
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

type listDomainsRes struct {
	Data interface{}
}

func (res listDomainsRes) Code() int {
	return http.StatusOK
}
func (res listDomainsRes) Headers() map[string]string {
	return map[string]string{}
}
func (res listDomainsRes) Empty() bool {
	return false
}
func (res listDomainsRes) MarshalJSON() ([]byte, error) {
	return json.Marshal(res.Data)
}

type enableDomainRes struct {
}

func (res enableDomainRes) Code() int {
	return http.StatusOK
}
func (res enableDomainRes) Headers() map[string]string {
	return map[string]string{}
}
func (res enableDomainRes) Empty() bool {
	return true
}

type disableDomainRes struct {
}

func (res disableDomainRes) Code() int {
	return http.StatusOK
}
func (res disableDomainRes) Headers() map[string]string {
	return map[string]string{}
}
func (res disableDomainRes) Empty() bool {
	return true
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

type listUserDomainsRes struct {
	Data interface{}
}

func (res listUserDomainsRes) Code() int {
	return http.StatusNoContent
}
func (res listUserDomainsRes) Headers() map[string]string {
	return map[string]string{}
}
func (res listUserDomainsRes) Empty() bool {
	return true
}
func (res listUserDomainsRes) MarshalJSON() ([]byte, error) {
	return json.Marshal(res.Data)
}
