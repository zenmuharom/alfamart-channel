package domain

import "database/sql"

type FieldValue struct {
	Id                int64          `db:"id"`
	FieldId           int64          `db:"field_id"`
	FieldName         sql.NullString `db:"field_name"`
	FieldType         string         `db:"field_type"`
	FieldAs           sql.NullString `db:"field_as"`
	FieldParentId     sql.NullInt64  `db:"field_parent"`
	FieldRefId        sql.NullInt64  `db:"field_ref_id"`
	GlobalField       sql.NullString `db:"global_field"`
	Args              sql.NullString `db:"args"`
	ConditionFieldId  sql.NullInt64  `db:"condition_field_id"`
	ConditionOperator sql.NullString `db:"condition_operator"`
	ConditionValue    sql.NullString `db:"condition_value"`
	CreatedAt         sql.NullTime   `db:"created_at"`
	UpdatedAt         sql.NullTime   `db:"updated_at"`
	DeletedAt         sql.NullTime   `db:"deleted_at"`
}
