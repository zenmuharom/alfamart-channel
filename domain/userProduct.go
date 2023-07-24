package domain

import "database/sql"

type UserProduct struct {
	Id                int64          `db:"id"`
	Username          sql.NullString `db:"username"`
	Is_Dsposit        sql.NullString `db:"is_deposit"`
	ProductCode       sql.NullString `db:"product_code"`
	Bit18             sql.NullString `db:"bit18"`
	Bit32             sql.NullString `db:"bit32"`
	Bit33             sql.NullString `db:"bit33"`
	Bit62             sql.NullString `db:"bit62"`
	ProductCodeMapped sql.NullString `db:"product_code_mapped"`
	RCSuccess         []uint8        `db:"rc_success"`
	RCErrorNonDebit   []uint8        `db:"rc_error_nondebit"`
	TimeoutBiller     sql.NullInt32  `db:"timeout_biller"`
	SwitchingUrl      sql.NullString `db:"switching_url"`
	AccountNumber     sql.NullString `db:"account_number"`
	CreatedAt         sql.NullString `db:"created_at"`
	ActivatedAt       sql.NullString `db:"activated_at"`
	UpdatedAt         sql.NullString `db:"updated_at"`
	DeletedAt         sql.NullString `db:"deleted_at"`
}
