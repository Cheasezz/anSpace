package core

import "time"

type CodeCredentials struct {
	Email     string    `db:"email"`
	Code      string    `db:"code"`
	CodeType  string    `db:"code_type"`
	ExpiresAt time.Time `db:"expires_at"`
}
