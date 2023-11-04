// Copyright (c) Magistrala
// SPDX-License-Identifier: Apache-2.0

package jwt

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/absmach/magistrala/auth"
	"github.com/absmach/magistrala/pkg/errors"
	svcerror "github.com/absmach/magistrala/pkg/errors/service"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

var (
	errInvalidIssuer = errors.New("invalid token issuer value")
	// errJWTExpiryKey is used to check if the token is expired.
	errJWTExpiryKey = errors.New(`"exp" not satisfied`)
	// ErrExpiry indicates that the token is expired.
	ErrExpiry = errors.New("token is expired")
)

const (
	issuerName = "magistrala.auth"
	tokenType  = "type"
)

type tokenizer struct {
	secret []byte
}

var _ auth.Tokenizer = (*tokenizer)(nil)

// NewRepository instantiates an implementation of Token repository.
func New(secret []byte) auth.Tokenizer {
	return &tokenizer{
		secret: secret,
	}
}

func (repo *tokenizer) Issue(key auth.Key) (string, error) {
	builder := jwt.NewBuilder()
	builder.
		Issuer(issuerName).
		IssuedAt(key.IssuedAt).
		Subject(key.Subject).
		Claim(tokenType, key.Type).
		Expiration(key.ExpiresAt)
	if key.ID != "" {
		builder.JwtID(key.ID)
	}
	tkn, err := builder.Build()
	if err != nil {
		return "", errors.Wrap(svcerror.ErrAuthentication, err)
	}
	signedTkn, err := jwt.Sign(tkn, jwt.WithKey(jwa.HS512, repo.secret))
	if err != nil {
		return "", err
	}
	return string(signedTkn), nil
}

func (repo *tokenizer) Parse(token string) (auth.Key, error) {
	tkn, err := jwt.Parse(
		[]byte(token),
		jwt.WithValidate(true),
		jwt.WithKey(jwa.HS512, repo.secret),
	)
	if err != nil {
		if errors.Contains(err, errJWTExpiryKey) {
			return auth.Key{}, ErrExpiry
		}

		return auth.Key{}, errors.Wrap(svcerror.ErrAuthentication, err)
	}
	validator := jwt.ValidatorFunc(func(_ context.Context, t jwt.Token) jwt.ValidationError {
		if t.Issuer() != issuerName {
			return jwt.NewValidationError(errInvalidIssuer)
		}
		return nil
	})
	if err := jwt.Validate(tkn, jwt.WithValidator(validator)); err != nil {
		return auth.Key{}, err
	}

	jsn, err := json.Marshal(tkn.PrivateClaims())
	if err != nil {
		return auth.Key{}, err
	}
	var key auth.Key
	if err := json.Unmarshal(jsn, &key); err != nil {
		return auth.Key{}, err
	}

	tType, ok := tkn.Get(tokenType)
	if !ok {
		return auth.Key{}, errors.Wrap(svcerror.ErrAuthentication, err)
	}
	ktype, err := strconv.ParseInt(fmt.Sprintf("%v", tType), 10, 64)
	if err != nil {
		return auth.Key{}, errors.Wrap(svcerror.ErrAuthentication, err)
	}

	key.ID = tkn.JwtID()
	key.Type = auth.KeyType(ktype)
	key.Issuer = tkn.Issuer()
	key.Subject = tkn.Subject()
	key.IssuedAt = tkn.IssuedAt()
	key.ExpiresAt = tkn.Expiration()
	return key, nil
}
