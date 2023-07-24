package domain

import "database/sql"

type Response struct {
	Id          sql.NullString `db:"id"`
	Field       sql.NullString `db:"field"`
	Type        sql.NullString `db:"type"`
	ParentId    sql.NullString `db:"parent_id"`
	ProductCode sql.NullString `db:"product_code"`
	CreatedAt   sql.NullString `db:"created_at"`
	ActivatedAt sql.NullString `db:"activated_at"`
	UpdatedAt   sql.NullString `db:"updated_at"`
	DeletedAt   sql.NullString `db:"deleted_at"`
}
