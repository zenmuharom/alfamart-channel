package domain

import "database/sql"

type MiddlewareRequest struct {
	Id                 int            `db:"id"`
	Process            sql.NullString `db:"process"`
	Field              sql.NullString `db:"field"`
	ParentId           sql.NullInt64  `db:"parent_id"`
	FieldParent        sql.NullString `db:"field_parent"`
	Type               sql.NullString `db:"type"`
	ConditionFieldId   sql.NullInt64  `db:"condition_field_id"`
	ConditionFieldName sql.NullString `db:"condition_field_name"`
	ConditionOperator  sql.NullString `db:"condition_operator"`
	ConditionValue     sql.NullString `db:"condition_value"`
	Middleware         sql.NullString `db:"middleware"`
	ProductCode        sql.NullString `db:"product_code"`
	CreatedAt          sql.NullTime   `db:"created_at"`
	UpdatedAt          sql.NullTime   `db:"updated_at"`
	ActivatedAt        sql.NullTime   `db:"activated_at"`
	DeletedAt          sql.NullTime   `db:"deleted_at"`
}
