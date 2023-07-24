package domain

import "database/sql"

type ServerResponse struct {
	Id                 int            `db:"id"`
	Endpoint           sql.NullString `db:"endpoint"`
	Field              sql.NullString `db:"field"`
	Type               sql.NullString `db:"type"`
	ParentId           sql.NullInt64  `db:"parent_id"`
	ProductCode        sql.NullString `db:"product_code"`
	FieldParent        sql.NullString `db:"field_parent"`
	FieldParentType    sql.NullString `db:"field_parent_type"`
	ConditionFieldId   sql.NullInt64  `db:"condition_field_id"`
	ConditionFieldName sql.NullString `db:"condition_field_name"`
	ConditionOperator  sql.NullString `db:"condition_operator"`
	ConditionValue     sql.NullString `db:"condition_value"`
	FieldAs            sql.NullString `db:"field_as"`
	CreatedAt          sql.NullString `db:"created_at"`
	ActivatedAt        sql.NullString `db:"activated_at"`
	UpdatedAt          sql.NullString `db:"updated_at"`
	DeletedAt          sql.NullString `db:"deleted_at"`
}
