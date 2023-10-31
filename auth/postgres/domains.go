// Copyright (c) Magistrala
// SPDX-License-Identifier: Apache-2.0

package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/absmach/magistrala/auth"
	"github.com/absmach/magistrala/internal/postgres"
	"github.com/absmach/magistrala/pkg/clients"
	"github.com/absmach/magistrala/pkg/errors"
	"github.com/jackc/pgtype"
	"github.com/jmoiron/sqlx"
)

var _ auth.DomainsRepository = (*domainRepo)(nil)

var (
	errRollbackTx = errors.New("failed to rollback transaction")
)

type domainRepo struct {
	db postgres.Database
}

// NewDomainRepository instantiates a PostgreSQL
// implementation of Domain repository.
func NewDomainRepository(db postgres.Database) auth.DomainsRepository {
	return &domainRepo{
		db: db,
	}
}

// Save
func (repo domainRepo) Save(ctx context.Context, d auth.Domain) (ad auth.Domain, err error) {

	q := `INSERT INTO domains (id, name, email, tags, metadata, created_at, updated_at, updated_by, created_by, status)
	VALUES (:id, :name, :email, :tags, :metadata, :created_at, :updated_at, :updated_by, :created_by, :status)
	RETURNING id, name, email, tags, metadata, created_at, updated_at, updated_by, created_by, status;`

	dbd, err := toDBDomains(d)
	if err != nil {
		return auth.Domain{}, errors.Wrap(errors.ErrCreateEntity, err)
	}

	row, err := repo.db.NamedQueryContext(ctx, q, dbd)
	if err != nil {
		return auth.Domain{}, postgres.HandleError(err, errors.ErrCreateEntity)
	}

	defer row.Close()
	row.Next()
	dbd = dbDomain{}
	if err := row.StructScan(&dbd); err != nil {
		return auth.Domain{}, err
	}

	domain, err := toDomain(dbd)
	if err != nil {
		return auth.Domain{}, err
	}

	return domain, nil
}

// RetrieveByID retrieves Domain by its unique ID.
func (repo domainRepo) RetrieveByID(ctx context.Context, id string) (auth.Domain, error) {
	q := `SELECT id, name, email, tags, metadata, created_at, updated_at, updated_by, created_by, status
        FROM domains WHERE id = :id`

	dbd := dbDomain{
		ID: id,
	}

	row, err := repo.db.NamedQueryContext(ctx, q, dbd)
	if err != nil {
		if err == sql.ErrNoRows {
			return auth.Domain{}, errors.Wrap(errors.ErrNotFound, err)
		}
		return auth.Domain{}, errors.Wrap(errors.ErrViewEntity, err)
	}

	defer row.Close()
	row.Next()
	dbd = dbDomain{}
	if err := row.StructScan(&dbd); err != nil {
		return auth.Domain{}, errors.Wrap(errors.ErrNotFound, err)
	}

	return toDomain(dbd)
}

// RetrieveAllByIDs retrieves for given Domain IDs .
func (repo domainRepo) RetrieveAllByIDs(ctx context.Context, pm auth.Page) (auth.DomainsPage, error) {
	var q string
	if len(pm.IDs) <= 0 {
		return auth.DomainsPage{}, nil
	}
	query, err := buildPageQuery(pm)
	if err != nil {
		return auth.DomainsPage{}, err
	}
	if query == "" {
		return auth.DomainsPage{}, nil
	}

	q = `SELECT id, name, email, tags, metadata, created_at, updated_at, updated_by, created_by, status
	FROM domains`
	q = fmt.Sprintf("%s %s ORDER BY :order :dir  LIMIT :limit OFFSET :offset;", q, query)

	dbPage, err := toDBClientsPage(pm)
	if err != nil {
		return auth.DomainsPage{}, errors.Wrap(postgres.ErrFailedToRetrieveAll, err)
	}

	rows, err := repo.db.NamedQueryContext(ctx, q, dbPage)
	if err != nil {
		return auth.DomainsPage{}, errors.Wrap(postgres.ErrFailedToRetrieveAll, err)
	}
	defer rows.Close()

	domains, err := repo.processRows(rows)
	if err != nil {
		return auth.DomainsPage{}, errors.Wrap(postgres.ErrFailedToRetrieveAll, err)
	}

	cq := "SELECT COUNT(*) FROM domains"
	if query != "" {
		cq = fmt.Sprintf(" %s %s", cq, query)
	}

	total, err := postgres.Total(ctx, repo.db, cq, dbPage)
	if err != nil {
		return auth.DomainsPage{}, errors.Wrap(postgres.ErrFailedToRetrieveAll, err)
	}

	pm.Total = total
	return auth.DomainsPage{
		Page:    pm,
		Domains: domains,
	}, nil
}

// Update updates the client name and metadata.
func (repo domainRepo) Update(ctx context.Context, dr auth.DomainReq) (auth.Domain, error) {
	var query []string
	var upq string
	var d auth.Domain
	if dr.Name != nil && *dr.Name != "" {
		query = append(query, "name = :name, ")
		d.Name = *dr.Name
	}
	if dr.Email != nil && *dr.Email != "" {
		query = append(query, "email = :email, ")
		d.Email = *dr.Email
	}
	if dr.Metadata != nil {
		query = append(query, "metadata = :metadata, ")
		d.Metadata = *dr.Metadata
	}
	if dr.Tags != nil {
		query = append(query, "tags = :tags, ")
		d.Tags = *dr.Tags
	}
	if dr.Status != nil {
		query = append(query, "status = :status, ")
		d.Status = *dr.Status
	}
	if len(query) > 0 {
		upq = strings.Join(query, " ")
	}
	q := fmt.Sprintf(`UPDATE clients SET %s  updated_at = :updated_at, updated_by = :updated_by
        WHERE id = :id AND status = :status
        RETURNING id, name, email, tags, metadata, created_at, updated_at, updated_by, created_by, status;`,
		upq)

	dbd, err := toDBDomains(d)
	if err != nil {
		return auth.Domain{}, errors.Wrap(errors.ErrCreateEntity, err)
	}
	row, err := repo.db.NamedQueryContext(ctx, q, dbd)
	if err != nil {
		return auth.Domain{}, postgres.HandleError(err, errors.ErrCreateEntity)
	}

	defer row.Close()
	row.Next()
	dbd = dbDomain{}
	if err := row.StructScan(&dbd); err != nil {
		return auth.Domain{}, err
	}

	domain, err := toDomain(dbd)
	if err != nil {
		return auth.Domain{}, err
	}

	return domain, nil
}

// Delete
func (repo domainRepo) Delete(ctx context.Context, id string) error {
	return nil
}

// SavePolicyCopy
func (repo domainRepo) SavePolicyCopy(ctx context.Context, pc auth.PolicyCopy) error {
	q := `INSERT INTO policies_copy (subject_type, subject_id, subject_relation, relation, object_type, object_id)
	VALUES (:subject_type, :subject_id, :subject_relation, :relation, :object_type, :object_id)
	RETURNING subject_type, subject_id, subject_relation, relation, object_type, object_id;`

	dbpc := toDBPolicyCopy(pc)
	row, err := repo.db.NamedQueryContext(ctx, q, dbpc)
	if err != nil {
		return postgres.HandleError(err, errors.ErrCreateEntity)
	}
	defer row.Close()

	return nil
}

// DeletePolicyCopy
func (repo domainRepo) DeletePolicyCopy(ctx context.Context, pc auth.PolicyCopy) (err error) {
	q := `
		DELETE FROM
			policies_copy
		WHERE
			subject_type = :subject_type
			AND subject_id = :subject_id
			AND subject_relation = :subject_relation
			AND relation = :relation
			AND object_type = :object_type
			AND object_id = :object_id
		;`

	dbpc := toDBPolicyCopy(pc)
	row, err := repo.db.NamedQueryContext(ctx, q, dbpc)
	if err != nil {
		return postgres.HandleError(err, errors.ErrRemoveEntity)
	}
	defer row.Close()

	return nil
}

func (repo domainRepo) processRows(rows *sqlx.Rows) ([]auth.Domain, error) {
	var items []auth.Domain
	for rows.Next() {
		dbd := dbDomain{}
		if err := rows.StructScan(&dbd); err != nil {
			return items, err
		}
		d, err := toDomain(dbd)
		if err != nil {
			return items, err
		}
		items = append(items, d)
	}
	return items, nil
}

type dbDomain struct {
	ID        string           `db:"id"`
	Name      string           `db:"name"`
	Email     string           `db:"Email"`
	Metadata  []byte           `db:"metadata,omitempty"`
	Tags      pgtype.TextArray `db:"tags,omitempty"`
	Alias     *string          `db:"alias,omitempty"`
	Status    clients.Status   `db:"status"`
	CreatedBy string           `db:"created_by"`
	CreatedAt time.Time        `db:"created_at"`
	UpdatedBy *string          `db:"updated_by,omitempty"`
	UpdatedAt sql.NullTime     `db:"updated_at,omitempty"`
}

func toDBDomains(d auth.Domain) (dbDomain, error) {
	data := []byte("{}")
	if len(d.Metadata) > 0 {
		b, err := json.Marshal(d.Metadata)
		if err != nil {
			return dbDomain{}, errors.Wrap(errors.ErrMalformedEntity, err)
		}
		data = b
	}
	var tags pgtype.TextArray
	if err := tags.Set(d.Tags); err != nil {
		return dbDomain{}, err
	}
	var alias *string
	if d.Alias != "" {
		alias = &d.Alias
	}
	var updatedBy *string
	if d.UpdatedBy != "" {
		updatedBy = &d.UpdatedBy
	}
	var updatedAt sql.NullTime
	if d.UpdatedAt != (time.Time{}) {
		updatedAt = sql.NullTime{Time: d.UpdatedAt, Valid: true}
	}

	return dbDomain{
		ID:        d.ID,
		Name:      d.Name,
		Email:     d.Email,
		Metadata:  data,
		Tags:      tags,
		Alias:     alias,
		Status:    d.Status,
		CreatedBy: d.CreatedBy,
		CreatedAt: d.CreatedAt,
		UpdatedBy: updatedBy,
		UpdatedAt: updatedAt,
	}, nil
}

func toDomain(d dbDomain) (auth.Domain, error) {
	var metadata clients.Metadata
	if d.Metadata != nil {
		if err := json.Unmarshal([]byte(d.Metadata), &metadata); err != nil {
			return auth.Domain{}, errors.Wrap(errors.ErrMalformedEntity, err)
		}
	}
	var tags []string
	for _, e := range d.Tags.Elements {
		tags = append(tags, e.String)
	}
	var alias string
	if d.Alias != nil {
		alias = *d.Alias
	}
	var updatedBy string
	if d.UpdatedBy != nil {
		updatedBy = *d.UpdatedBy
	}
	var updatedAt time.Time
	if d.UpdatedAt.Valid {
		updatedAt = d.UpdatedAt.Time
	}

	return auth.Domain{
		ID:        d.ID,
		Name:      d.Name,
		Email:     d.Email,
		Metadata:  metadata,
		Tags:      tags,
		Alias:     alias,
		Status:    d.Status,
		CreatedBy: d.CreatedBy,
		CreatedAt: d.CreatedAt,
		UpdatedBy: updatedBy,
		UpdatedAt: updatedAt,
	}, nil
}

type dbDomainsPage struct {
	Total    uint64         `db:"total"`
	Limit    uint64         `db:"limit"`
	Offset   uint64         `db:"offset"`
	Order    string         `db:"order"`
	Dir      string         `db:"dir"`
	Name     string         `db:"name"`
	Email    string         `db:"email"`
	ID       string         `db:"id"`
	IDs      []string       `db:"ids"`
	Metadata []byte         `db:"metadata"`
	Tag      string         `db:"tag"`
	Status   clients.Status `db:"status"`
}

func toDBClientsPage(pm auth.Page) (dbDomainsPage, error) {
	_, data, err := postgres.CreateMetadataQuery("", pm.Metadata)
	if err != nil {
		return dbDomainsPage{}, errors.Wrap(errors.ErrViewEntity, err)
	}
	return dbDomainsPage{
		Total:    pm.Total,
		Limit:    pm.Limit,
		Offset:   pm.Offset,
		Order:    pm.Order,
		Dir:      pm.Dir,
		Name:     pm.Name,
		Email:    pm.Email,
		ID:       pm.ID,
		IDs:      pm.IDs,
		Metadata: data,
		Tag:      pm.Tag,
		Status:   pm.Status,
	}, nil
}

func buildPageQuery(pm auth.Page) (string, error) {
	var query []string
	var emq string

	if pm.ID != "" {
		query = append(query, "id = :id")
	}

	if len(pm.IDs) != 0 {
		query = append(query, fmt.Sprintf("id IN ('%s')", strings.Join(pm.IDs, "','")))
	}

	if pm.Status != clients.AllStatus {
		query = append(query, "d.status = :status")
	}

	if pm.Email != "" {
		query = append(query, "email = :email")
	}

	if pm.Name != "" {
		query = append(query, "name = :name")
	}

	if pm.Tag != "" {
		query = append(query, ":tag = ANY(d.tags)")
	}

	mq, _, err := postgres.CreateMetadataQuery("", pm.Metadata)
	if err != nil {
		return "", errors.Wrap(errors.ErrViewEntity, err)
	}
	if mq != "" {
		query = append(query, mq)
	}

	if len(query) > 0 {
		emq = fmt.Sprintf("WHERE %s", strings.Join(query, " AND "))
	}

	return emq, nil
}

type dbPolicyCopy struct {
	SubjectType     string `db:"subject_type,omitempty"`
	SubjectID       string `db:"subject_id,omitempty"`
	SubjectRelation string `db:"subject_relation,omitempty"`
	Relation        string `db:"relation,omitempty"`
	ObjectType      string `db:"object_type,omitempty"`
	ObjectID        string `db:"object_id,omitempty"`
}

func toDBPolicyCopy(pc auth.PolicyCopy) dbPolicyCopy {
	return dbPolicyCopy{
		SubjectType:     pc.SubjectType,
		SubjectID:       pc.SubjectID,
		SubjectRelation: pc.SubjectRelation,
		Relation:        pc.Relation,
		ObjectType:      pc.ObjectType,
		ObjectID:        pc.ObjectID,
	}
}
