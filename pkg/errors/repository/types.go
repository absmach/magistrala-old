// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package repository

import "github.com/absmach/magistrala/pkg/errors"

// Wrapper for Repository errors.
var (
	// ErrMalformedEntity indicates a malformed entity specification.
	ErrMalformedEntity = errors.New("malformed entity specification")

	// ErrNotFound indicates a non-existent entity request.
	ErrNotFound = errors.New("entity not found")

	// ErrConflict indicates that entity already exists.
	ErrConflict = errors.New("entity already exists")

	// ErrCreateEntity indicates error in creating entity or entities.
	ErrCreateEntity = errors.New("failed to create entity in the db")

	// ErrViewEntity indicates error in viewing entity or entities.
	ErrViewEntity = errors.New("view entity failed")

	// ErrUpdateEntity indicates error in updating entity or entities.
	ErrUpdateEntity = errors.New("update entity failed")

	// ErrRemoveEntity indicates error in removing entity.
	ErrRemoveEntity = errors.New("failed to remove entity")

	// ErrScanMetadata indicates problem with metadata in db.
	ErrScanMetadata = errors.New("failed to scan metadata in db")

	// ErrWrongSecret indicates a wrong secret was provided.
	ErrWrongSecret = errors.New("wrong secret")

	// ErrLogin indicates wrong login credentials.
	ErrLogin = errors.New("invalid user id or secret")

	//ErrUniqueID indicates an error in generating a unique ID.
	ErrUniqueID = errors.New("failed to generate unique identifier")

	//ErrFailedOpDB indicates a failure in a database operation.
	ErrFailedOpDB = errors.New("operation on db element failed")
)
