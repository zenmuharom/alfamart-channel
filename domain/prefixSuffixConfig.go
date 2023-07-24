package domain

import "database/sql"

type PrefixSuffixConfig struct {
	Id             sql.NullString `db:"id"`
	PresufType     sql.NullString `db:"presuf_type"`
	FixNumberStart sql.NullInt64  `db:"fix_number_start"`
	FixNumberEnd   sql.NullInt64  `db:"fix_number_end"`
	FixCharacter   sql.NullString `db:"fix_character"`
	Prefix         sql.NullString `db:"prefix"`
	SuffixLength   sql.NullInt64  `db:"suffix_length"`
	Length         sql.NullInt64  `db:"length"`
	ProductCode    sql.NullString `db:"product_code"`
	CreatedAt      sql.NullTime   `db:"created_at"`
	ActivatedAt    sql.NullTime   `db:"activated_at"`
	UpdatedAt      sql.NullTime   `db:"updated_at"`
	DeletedAt      sql.NullTime   `db:"deleted_at"`
}
