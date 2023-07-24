package domain

import "database/sql"

type Fee struct {
	Id                  int64           `db:"id"`
	Username            sql.NullString  `db:"username"`
	ProductCode         sql.NullString  `db:"product_code"`
	FeeMin              sql.NullFloat64 `db:"fee_min"`
	FeeMax              sql.NullFloat64 `db:"fee_max"`
	FeeCaFix            sql.NullInt64   `db:"fee_ca_fix"`
	FeeCaPercentage     sql.NullFloat64 `db:"fee_ca_percentage"`
	FeeBillerFix        sql.NullInt64   `db:"fee_biller_fix"`
	FeeBillerPercentage sql.NullFloat64 `db:"fee_biller_percentage"`
	FeeFinnetFix        sql.NullInt64   `db:"fee_finnet_fix"`
	FeeFinnetPercentage sql.NullFloat64 `db:"fee_finnet_percentage"`
	IsSurcharge         sql.NullBool    `db:"is_surcharge"`
	CreatedAt           sql.NullTime    `db:"created_at"`
	UpdatedAt           sql.NullTime    `db:"updated_at"`
	DeletedAt           sql.NullTime    `db:"deleted_at"`
}
