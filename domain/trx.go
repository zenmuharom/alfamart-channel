package domain

import "database/sql"

type Trx struct {
	Pid            string          `db:"pid", json:"pid"`
	AmendmentDate  sql.NullTime    `db:"amendment_date", json:"amendment_date"`
	SourceMerchant sql.NullString  `db:"source_merchant", json:"source_merchant"`
	TargetProduct  sql.NullString  `db:"target_product", json:"target_product"`
	Amount         sql.NullFloat64 `db:"amount", json:"amount"`
	Status         sql.NullString  `db:"status", json:"status"`
	Rc             sql.NullString  `db:"rc", json:"rc"`
	RcDesc         sql.NullString  `db:"rc_desc", json:"rc_desc"`
	ElapsedTime    sql.NullInt64   `db:"elapsed_time", json:"elapsed_time"`
}

type TrxDetail struct {
}
