package domain

import "database/sql"

type MiddlewareRequestValue struct {
	Id                       int64          `db:"id"`
	FieldId                  sql.NullInt64  `db:"field_id"`
	FieldName                sql.NullString `db:"field_name"`
	FieldType                sql.NullString `db:"field_type"`
	ServerRequestId          sql.NullInt64  `db:"server_request_id"`
	ServerRequestField       sql.NullString `db:"server_request_field"`
	ServerRequestParentField sql.NullString `db:"server_request_parent_field"`
	ServerRequestValue       sql.NullString `db:"server_request_value"`
	Args                     sql.NullString `db:"args"`
	ParentId                 sql.NullInt64  `db:"parent_id"`
	ConditionFieldId         sql.NullInt64  `db:"condition_field_id"`
	ConditionFieldName       sql.NullString `db:"condition_field_name"`
	ConditionOperator        sql.NullString `db:"condition_operator"`
	ConditionValue           sql.NullString `db:"condition_value"`
	CaCode                   sql.NullString `db:"ca_code"`
	ProductCode              sql.NullString `db:"product_code"`
	CreatedAt                sql.NullTime   `db:"created_at"`
	UpdatedAt                sql.NullTime   `db:"updated_at"`
	ActivatedAt              sql.NullTime   `db:"activated_at"`
	DeletedAt                sql.NullTime   `db:"deleted_at"`
}
