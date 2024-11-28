package core

import (
	"github.com/google/uuid"
)

type User struct {
	Id            uuid.UUID `json:"-" db:"id"`
	Email         string    `json:"email" db:"email"`
	Username      string    `json:"username" db:"username"`
	PasswordHash  string    `json:"passwordHash" db:"password_hash"`
}
