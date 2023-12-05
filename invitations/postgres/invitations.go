// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/absmach/magistrala/internal/postgres"
	"github.com/absmach/magistrala/invitations"
	"github.com/absmach/magistrala/pkg/errors"
)

type repository struct {
	db postgres.Database
}

func NewRepository(db postgres.Database) invitations.Repository {
	return &repository{db: db}
}

func (repo *repository) Create(ctx context.Context, invitation invitations.Invitation) (err error) {
	q := `INSERT INTO invitations (invited_by, user_id, domain, token, relation, created_at, updated_at, confirmed_at)
		VALUES (:invited_by, :user_id, :domain, :token, :relation, :created_at, :updated_at, :confirmed_at)`

	if _, err = repo.db.NamedExecContext(ctx, q, invitation); err != nil {
		return postgres.HandleError(errors.ErrCreateEntity, err)
	}

	return nil
}

func (repo *repository) Retrieve(ctx context.Context, userID, domainID string) (invitations.Invitation, error) {
	q := `SELECT invited_by, user_id, domain, relation, created_at, updated_at, confirmed_at FROM invitations WHERE user_id = :user_id AND domain = :domain`

	inv := invitations.Invitation{
		UserID: userID,
		Domain: domainID,
	}

	rows, err := repo.db.NamedQueryContext(ctx, q, inv)
	if err != nil {
		if err == sql.ErrNoRows {
			return invitations.Invitation{}, errors.ErrNotFound
		}

		return invitations.Invitation{}, postgres.HandleError(errors.ErrViewEntity, err)
	}
	defer rows.Close()

	var item invitations.Invitation
	if rows.Next() {
		if err = rows.StructScan(&item); err != nil {
			return invitations.Invitation{}, postgres.HandleError(errors.ErrViewEntity, err)
		}

		return item, nil
	}

	return invitations.Invitation{}, errors.ErrNotFound
}

func (repo *repository) RetrieveAll(ctx context.Context, page invitations.Page) (invitations.InvitationPage, error) {
	query := pageQuery(page)
	q := fmt.Sprintf(`SELECT invited_by, user_id, domain, relation, created_at, updated_at, confirmed_at FROM invitations %s LIMIT :limit OFFSET :offset`, query)

	rows, err := repo.db.NamedQueryContext(ctx, q, page)
	if err != nil {
		return invitations.InvitationPage{}, postgres.HandleError(errors.ErrViewEntity, err)
	}
	defer rows.Close()

	var items []invitations.Invitation
	for rows.Next() {
		var item invitations.Invitation
		if err = rows.StructScan(&item); err != nil {
			return invitations.InvitationPage{}, postgres.HandleError(errors.ErrViewEntity, err)
		}
		items = append(items, item)
	}

	tq := fmt.Sprintf(`SELECT COUNT(*) FROM invitations %s`, query)

	total, err := postgres.Total(ctx, repo.db, tq, page)
	if err != nil {
		return invitations.InvitationPage{}, postgres.HandleError(errors.ErrViewEntity, err)
	}

	invPage := invitations.InvitationPage{
		Total:       total,
		Offset:      page.Offset,
		Limit:       page.Limit,
		Invitations: items,
	}

	return invPage, nil
}

func (repo *repository) UpdateToken(ctx context.Context, invitation invitations.Invitation) (err error) {
	q := `UPDATE invitations SET token = :token, updated_at = :updated_at WHERE user_id = :user_id AND domain = :domain`

	if _, err = repo.db.NamedExecContext(ctx, q, invitation); err != nil {
		return postgres.HandleError(errors.ErrUpdateEntity, err)
	}

	return nil
}

func (repo *repository) UpdateConfirmation(ctx context.Context, invitation invitations.Invitation) (err error) {
	q := `UPDATE invitations SET confirmed_at = :confirmed_at WHERE user_id = :user_id AND domain = :domain`

	if _, err = repo.db.NamedExecContext(ctx, q, invitation); err != nil {
		return postgres.HandleError(errors.ErrUpdateEntity, err)
	}

	return nil
}

func (repo *repository) Delete(ctx context.Context, userID, domain string) (err error) {
	q := `DELETE FROM invitations WHERE user_id = $1 AND domain = $2`

	if _, err = repo.db.ExecContext(ctx, q, userID, domain); err != nil {
		return postgres.HandleError(errors.ErrRemoveEntity, err)
	}

	return nil
}

func pageQuery(pm invitations.Page) string {
	var query []string
	var emq string
	if pm.Domain != "" {
		query = append(query, "domain = :domain")
	}
	if pm.UserID != "" {
		query = append(query, "user_id = :user_id")
	}
	if pm.InvitedBy != "" {
		query = append(query, "invited_by = :invited_by")
	}
	if pm.Relation != "" {
		query = append(query, "relation = :relation")
	}

	if len(query) > 0 {
		emq = fmt.Sprintf("WHERE %s", strings.Join(query, " AND "))
	}

	return emq
}
