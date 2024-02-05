package domain

import "database/sql"

type ServerResponseValue struct {
	Id                      int64          `db:"id"`
	FieldId                 sql.NullInt64  `db:"field_id"`
	FieldName               sql.NullString `db:"field_name"`
	FieldType               sql.NullString `db:"field_type"`
	MiddlewareResponseId    sql.NullInt64  `db:"middleware_response_id"`
	MiddlewareResponseField sql.NullString `db:"middleware_response_field"`
	MiddlewareResponseValue sql.NullString `db:"middleware_response_value"` //buat diassign dari middleware response
	Args                    sql.NullString `db:"args"`
	ParentId                sql.NullInt64  `db:"parent_id"`
	ConditionFieldId        sql.NullInt64  `db:"condition_field_id"`
	ConditionFieldName      sql.NullString `db:"condition_field_name"`
	ConditionOperator       sql.NullString `db:"condition_operator"`
	ConditionValue          sql.NullString `db:"condition_value"`
	ProductCode             sql.NullString `db:"product_code"`
	FieldAs                 sql.NullString `db:"field_as"`
	CreatedAt               sql.NullTime   `db:"created_at"`
	UpdatedAt               sql.NullTime   `db:"updated_at"`
	ActivatedAt             sql.NullTime   `db:"activated_at"`
	DeletedAt               sql.NullTime   `db:"deleted_at"`
}
