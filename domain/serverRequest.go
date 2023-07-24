package domain

import "database/sql"

type ServerRequest struct {
	Id           int64          `db:"id"`
	Endpoint     sql.NullString `db:"endpoint"`
	Order        sql.NullInt64  `db:"order"`
	Field        sql.NullString `db:"field"`
	Type         sql.NullString `db:"type"`
	Length       sql.NullInt64  `db:"length"`
	Required     sql.NullBool   `db:"required"`
	Forbidded    sql.NullBool   `db:"forbidded"`
	ParentId     sql.NullInt64  `db:"parent_id"`
	FieldAs      sql.NullString `db:"field_as"`
	FieldParent  sql.NullString `db:"field_parent"`
	ProductCode  sql.NullString `db:"productCode"`
	Created_At   sql.NullString `db:"created_at"`
	Activated_At sql.NullString `db:"activated_at"`
	Updated_At   sql.NullString `db:"updated_at"`
	Deleted_At   sql.NullString `db:"deleted_at"`
}
