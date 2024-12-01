package core

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	UserId       uuid.UUID `db:"user_id"`
	RefreshToken string    `db:"refresh_token"`
	ExpiresAt    time.Time `db:"expires_at"`
}

type AuthCredentials struct {
	Email    string `json:"email" binding:"required" db:"email" example:"example@gmail.com"`
	Password string `json:"password" binding:"required" db:"password_hash" example:"qwerty123456"`
}
type Email struct {
	Email string `json:"email" binding:"required" db:"email" example:"example@gmail.com"`
}
type Password struct {
	Password string `json:"password" binding:"required" db:"password_hash" example:"qwerty123456"`
}
