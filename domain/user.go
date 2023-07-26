package domain

import "database/sql"

type User struct {
	Username      sql.NullString `db:"username"`
	Password      sql.NullString `db:"password"`
	Secret        sql.NullString `db:"secret"`
	IsDeposit     sql.NullBool   `db:"is_deposit"`
	NevaUsername  sql.NullString `db:"neva_username"`
	NevaPassword  sql.NullString `db:"neva_password"`
	ChannelCode   sql.NullString `db:"channel_code"`
	AccountNumber sql.NullString `db:"account_number"`
	Mitraco       sql.NullString `db:"mitraco"`
	CreatedAt     sql.NullTime   `db:"created_at"`
	ActivatedAt   sql.NullTime   `db:"activated_at"`
	UpdatedAt     sql.NullTime   `db:"updated_at"`
	DeletedAt     sql.NullTime   `db:"deleted_at"`
}
