package service

import (
	"github.com/absmach/magistrala/pkg/errors"
)

// Wrapper for Service errors
var (
	// ErrAuthentication indicates failure occurred while authenticating the entity.
	ErrAuthentication = errors.New("failed to perform authentication over the entity")

	// ErrAuthorization indicates failure occurred while authorizing the entity.
	ErrAuthorization = errors.New("failed to perform authorization over the entity")

	// ErrMalformedEntity indicates a malformed entity specification.
	ErrMalformedEntity = errors.New("malformed entity specification")

	// ErrNotFound indicates a non-existent entity request.
	ErrNotFound = errors.New("entity not found")

	// ErrConflict indicates that entity already exists.
	ErrConflict = errors.New("entity already exists")

	// ErrCreateEntity indicates error in creating entity or entities.
	ErrCreateEntity = errors.New("failed to create entity in the db")

	// ErrRemoveEntity indicates error in removing entity.
	ErrRemoveEntity = errors.New("failed to remove entity")
)
