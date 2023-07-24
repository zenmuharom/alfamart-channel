package domain

import "database/sql"

type RcConfig struct {
	Id          int64          `db:"id"`
	Code        sql.NullString `db:"code"`
	GoCode      sql.NullInt64  `db:"go_code"`
	ProductCode sql.NullString `db:"product_code"`
	Httpstatus  sql.NullInt64  `db:"httpstatus"`
	DescEng     sql.NullString `db:"desc_eng"`
	DescId      sql.NullString `db:"desc_id"`
	CreatedAt   sql.NullTime   `db:"created_at"`
	UpdatedAt   sql.NullTime   `db:"updated_at"`
	DeletedAt   sql.NullTime   `db:"deleted_at"`
}
