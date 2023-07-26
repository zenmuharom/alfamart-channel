package domain

import "database/sql"

type Route struct {
	Id          int64          `db:"id"`
	Method      sql.NullString `db:"method"`
	Path        sql.NullString `db:"path"`
	CreatedAt   sql.NullTime   `db:"created_at"`
	UpdatedAt   sql.NullTime   `db:"updated_at"`
	ActivatedAt sql.NullTime   `db:"activated_at"`
	DeletedAt   sql.NullTime   `db:"deleted_at"`
}
