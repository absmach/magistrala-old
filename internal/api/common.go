// Copyright (c) Magistrala
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/absmach/magistrala"
	"github.com/absmach/magistrala/internal/apiutil"
	"github.com/absmach/magistrala/internal/postgres"
	mgclients "github.com/absmach/magistrala/pkg/clients"
	"github.com/absmach/magistrala/pkg/errors"
	repoerror "github.com/absmach/magistrala/pkg/errors/repository"
	svcerror "github.com/absmach/magistrala/pkg/errors/service"

	"github.com/gofrs/uuid"
)

const (
	MemberKindKey    = "member_kind"
	PermissionKey    = "permission"
	RelationKey      = "relation"
	StatusKey        = "status"
	OffsetKey        = "offset"
	OrderKey         = "order"
	LimitKey         = "limit"
	MetadataKey      = "metadata"
	ParentKey        = "parent_id"
	OwnerKey         = "owner_id"
	ClientKey        = "client"
	IdentityKey      = "identity"
	GroupKey         = "group"
	ActionKey        = "action"
	TagKey           = "tag"
	NameKey          = "name"
	TotalKey         = "total"
	SubjectKey       = "subject"
	ObjectKey        = "object"
	LevelKey         = "level"
	TreeKey          = "tree"
	DirKey           = "dir"
	VisibilityKey    = "visibility"
	SharedByKey      = "shared_by"
	TokenKey         = "token"
	DefPermission    = "view"
	DefTotal         = uint64(100)
	DefOffset        = 0
	DefOrder         = "updated_by"
	DefDir           = "asc"
	DefLimit         = 10
	DefLevel         = 0
	DefStatus        = "enabled"
	DefClientStatus  = mgclients.Enabled
	DefGroupStatus   = mgclients.Enabled
	SharedVisibility = "shared"
	MyVisibility     = "mine"
	AllVisibility    = "all"
	// ContentType represents JSON content type.
	ContentType = "application/json"

	// MaxNameSize limits name size to prevent making them too complex.
	MaxLimitSize = 100
	MaxNameSize  = 1024
	NameOrder    = "name"
	IDOrder      = "id"
	AscDir       = "asc"
	DescDir      = "desc"
)

// ValidateUUID validates UUID format.
func ValidateUUID(extID string) (err error) {
	id, err := uuid.FromString(extID)
	if id.String() != extID || err != nil {
		return apiutil.ErrInvalidIDFormat
	}

	return nil
}

// EncodeResponse encodes successful response.
func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	if ar, ok := response.(magistrala.Response); ok {
		for k, v := range ar.Headers() {
			w.Header().Set(k, v)
		}
		w.Header().Set("Content-Type", ContentType)
		w.WriteHeader(ar.Code())

		if ar.Empty() {
			return nil
		}
	}

	return json.NewEncoder(w).Encode(response)
}

// EncodeError encodes an error response.
func EncodeError(_ context.Context, err error, w http.ResponseWriter) {
	var wrapper error
	if errors.Contains(err, apiutil.ErrValidation) {
		wrapper, err = errors.Unwrap(err)
	}

	w.Header().Set("Content-Type", ContentType)
	switch {
	case errors.Contains(err, apiutil.ErrInvalidSecret),
		errors.Contains(err, repoerror.ErrMalformedEntity),
		errors.Contains(err, apiutil.ErrMissingID),
		errors.Contains(err, apiutil.ErrEmptyList),
		errors.Contains(err, apiutil.ErrMissingMemberType),
		errors.Contains(err, apiutil.ErrMissingMemberKind),
		errors.Contains(err, apiutil.ErrNameSize):
		w.WriteHeader(http.StatusBadRequest)
	case errors.Contains(err, svcerror.ErrAuthentication):
		w.WriteHeader(http.StatusUnauthorized)
	case errors.Contains(err, repoerror.ErrNotFound):
		w.WriteHeader(http.StatusNotFound)
	case errors.Contains(err, repoerror.ErrConflict):
		w.WriteHeader(http.StatusConflict)
	case errors.Contains(err, svcerror.ErrAuthorization):
		w.WriteHeader(http.StatusForbidden)
	case errors.Contains(err, postgres.ErrMemberAlreadyAssigned):
		w.WriteHeader(http.StatusConflict)
	case errors.Contains(err, apiutil.ErrUnsupportedContentType):
		w.WriteHeader(http.StatusUnsupportedMediaType)
	case errors.Contains(err, repoerror.ErrCreateEntity),
		errors.Contains(err, repoerror.ErrUpdateEntity),
		errors.Contains(err, repoerror.ErrViewEntity),
		errors.Contains(err, repoerror.ErrRemoveEntity):
		w.WriteHeader(http.StatusInternalServerError)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	if wrapper != nil {
		err = errors.Wrap(wrapper, err)
	}

	if errorVal, ok := err.(errors.Error); ok {
		if err := json.NewEncoder(w).Encode(errorVal); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
