package domain

import "database/sql"

type Product struct {
	Id          int64          `db:"id"`
	Code        sql.NullString `db:"code"`
	Name        sql.NullString `db:"name"`
	Bit103      sql.NullString `db:"bit103"`
	Bit104      sql.NullString `db:"bit104"`
	MitraId     sql.NullInt64  `db:"mitra_id"`
	CreatedAt   sql.NullTime   `db:"created_at"`
	ActivatedAt sql.NullTime   `db:"activated_at"`
	UpdatedAt   sql.NullTime   `db:"updated_at"`
	DeletedAt   sql.NullTime   `db:"deleted_at"`
}
