package domain

import "database/sql"

type MiddlewareResponse struct {
	Id           int64          `db:"id"`
	Process      sql.NullString `db:"process"`
	Field        sql.NullString `db:"field"`
	Type         sql.NullString `db:"type"`
	Middleware   sql.NullString `db:"middleware"`
	FieldAs      sql.NullString `db:"field_as"`
	ParentId     sql.NullInt64  `db:"parent_id"`
	FieldParent  sql.NullString `db:"field_parent"`
	Created_At   sql.NullTime   `db:"created_at"`
	Updated_At   sql.NullTime   `db:"updated_at"`
	Activated_At sql.NullTime   `db:"activated_at"`
	Deleted_At   sql.NullTime   `db:"deleted_at"`
}
