// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/absmach/magistrala/internal/postgres"
	mgclients "github.com/absmach/magistrala/pkg/clients"
	pgclients "github.com/absmach/magistrala/pkg/clients/postgres"
	"github.com/absmach/magistrala/pkg/errors"
	repoerr "github.com/absmach/magistrala/pkg/errors/repository"
)

var _ mgclients.Repository = (*clientRepo)(nil)

type clientRepo struct {
	pgclients.ClientRepository
}

type Repository interface {
	mgclients.Repository

	// Save persists the client account. A non-nil error is returned to indicate
	// operation failure.
	Save(ctx context.Context, client mgclients.Client) (mgclients.Client, error)

	CheckSuperAdmin(ctx context.Context, adminID string) error

	// RetrieveNames returns a list of client names that match the given query.
	RetrieveNames(ctx context.Context, pm mgclients.Page) (mgclients.ClientsPage, error)
}

// NewRepository instantiates a PostgreSQL
// implementation of Clients repository.
func NewRepository(db postgres.Database) Repository {
	return &clientRepo{
		ClientRepository: pgclients.ClientRepository{DB: db},
	}
}

func (repo clientRepo) Save(ctx context.Context, c mgclients.Client) (mgclients.Client, error) {
	q := `INSERT INTO clients (id, name, tags, owner_id, identity, secret, metadata, created_at, status, role)
        VALUES (:id, :name, :tags, :owner_id, :identity, :secret, :metadata, :created_at, :status, :role)
        RETURNING id, name, tags, identity, metadata, COALESCE(owner_id, '') AS owner_id, status, created_at`
	dbc, err := pgclients.ToDBClient(c)
	if err != nil {
		return mgclients.Client{}, errors.Wrap(repoerr.ErrCreateEntity, err)
	}

	row, err := repo.ClientRepository.DB.NamedQueryContext(ctx, q, dbc)
	if err != nil {
		return mgclients.Client{}, postgres.HandleError(repoerr.ErrCreateEntity, err)
	}

	defer row.Close()
	row.Next()
	dbc = pgclients.DBClient{}
	if err := row.StructScan(&dbc); err != nil {
		return mgclients.Client{}, errors.Wrap(repoerr.ErrFailedOpDB, err)
	}

	client, err := pgclients.ToClient(dbc)
	if err != nil {
		return mgclients.Client{}, errors.Wrap(repoerr.ErrFailedOpDB, err)
	}

	return client, nil
}

func (repo clientRepo) CheckSuperAdmin(ctx context.Context, adminID string) error {
	q := "SELECT 1 FROM clients WHERE id = $1 AND role = $2"
	rows, err := repo.ClientRepository.DB.QueryContext(ctx, q, adminID, mgclients.AdminRole)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.ErrAuthorization
		}
		return errors.Wrap(errors.ErrAuthorization, err)
	}
	defer rows.Close()
	if !rows.Next() {
		return errors.ErrAuthorization
	}
	if err := rows.Err(); err != nil {
		return errors.Wrap(errors.ErrAuthorization, err)
	}
	return nil
}

func (repo clientRepo) RetrieveNames(ctx context.Context, pm mgclients.Page) (mgclients.ClientsPage, error) {
	query := constructQuery(pm)

	q := fmt.Sprintf("SELECT name FROM clients %s LIMIT :limit OFFSET :offset", query)

	dbPage, err := pgclients.ToDBClientsPage(pm)
	if err != nil {
		return mgclients.ClientsPage{}, errors.Wrap(postgres.ErrFailedToRetrieveAll, err)
	}

	rows, err := repo.DB.NamedQueryContext(ctx, q, dbPage)
	if err != nil {
		return mgclients.ClientsPage{}, errors.Wrap(postgres.ErrFailedToRetrieveAll, err)
	}
	defer rows.Close()

	var items []mgclients.Client
	for rows.Next() {
		dbc := pgclients.DBClient{}
		if err := rows.StructScan(&dbc); err != nil {
			return mgclients.ClientsPage{}, errors.Wrap(repoerr.ErrViewEntity, err)
		}

		c, err := pgclients.ToClient(dbc)
		if err != nil {
			return mgclients.ClientsPage{}, err
		}

		items = append(items, c)
	}
	cq := fmt.Sprintf(`SELECT COUNT(*) FROM clients c %s;`, query)

	total, err := postgres.Total(ctx, repo.DB, cq, dbPage)
	if err != nil {
		return mgclients.ClientsPage{}, errors.Wrap(repoerr.ErrViewEntity, err)
	}

	page := mgclients.ClientsPage{
		Clients: items,
		Page: mgclients.Page{
			Total:  total,
			Offset: pm.Offset,
			Limit:  pm.Limit,
		},
	}

	return page, nil
}

func constructQuery(pm mgclients.Page) string {
	var query []string
	var emq string
	if pm.Name != "" {
		query = append(query, fmt.Sprintf("name ILIKE '%%%s%%'", pm.Name))
	}
	if pm.Identity != "" {
		query = append(query, fmt.Sprintf("identity ILIKE '%%%s%%'", pm.Identity))
	}

	if len(query) > 0 {
		emq = fmt.Sprintf("WHERE %s", strings.Join(query, " AND "))
	}

	if pm.Order != "" && (pm.Order == "name" || pm.Order == "email" || pm.Order == "created_at" || pm.Order == "updated_at") {
		emq = fmt.Sprintf("%s ORDER BY %s", emq, pm.Order)
	}

	if pm.Dir != "" && (pm.Dir == "asc" || pm.Dir == "desc") {
		emq = fmt.Sprintf("%s %s", emq, pm.Dir)
	}

	return emq
}
