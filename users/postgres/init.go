// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package postgres

import (
	_ "github.com/jackc/pgx/v5/stdlib" // required for SQL access
	migrate "github.com/rubenv/sql-migrate"
)

// Migration of Users service.
func Migration() *migrate.MemoryMigrationSource {
	return &migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{
			{
				Id: "clients_01",
				// VARCHAR(36) for colums with IDs as UUIDS have a maximum of 36 characters
				// STATUS 0 to imply enabled and 1 to imply disabled
				// Role 0 to imply user role and 1 to imply admin role
				Up: []string{
					`CREATE TABLE IF NOT EXISTS clients (
						id          VARCHAR(36) PRIMARY KEY,
						name        VARCHAR(254),
						owner_id    VARCHAR(36),
						identity    VARCHAR(254) NOT NULL UNIQUE,
						secret      TEXT NOT NULL,
						tags        TEXT[],
						metadata    JSONB,
						created_at  TIMESTAMP,
						updated_at  TIMESTAMP,
						updated_by  VARCHAR(254),
						status      SMALLINT NOT NULL DEFAULT 0 CHECK (status >= 0),
						role        SMALLINT DEFAULT 0 CHECK (status >= 0)
					)`,
				},
				Down: []string{
					`DROP TABLE IF EXISTS clients`,
				},
			},
			{
				Id: "clients_02",
				Up: []string{
					`ALTER TABLE clients ALTER COLUMN name SET NOT NULL`,
					`ALTER TABLE clients ADD CONSTRAINT clients_name_unique UNIQUE (name)`,
				},
				Down: []string{
					`ALTER TABLE clients ALTER COLUMN name DROP NOT NULL`,
					`ALTER TABLE clients DROP CONSTRAINT clients_name_unique`,
				},
			},
		},
	}
}
