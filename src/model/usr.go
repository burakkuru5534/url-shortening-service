package model

type Usr struct {
	ID                 int64  `db:"id"`
	Email              string `db:"email"`
	Password           string `db:"password"`
	ShorteningUrlLimit int64  `db:"shortening_url_limit"`
	AccountType        string `db:"account_type"`
}
