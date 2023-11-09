package service

import (
	"github.com/absmach/magistrala/pkg/errors"
)

// Wrapper for Service errors
var (
	// ErrAuthentication indicates failure occurred while authenticating the entity.
	ErrAuthentication = errors.New("authentication error")

	// ErrAuthorization indicates failure occurred while authorizing the entity.
	ErrAuthorization = errors.New("permission denied")
)
