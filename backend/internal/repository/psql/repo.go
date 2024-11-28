package psql

import (
	"github.com/Cheasezz/anSpace/backend/pkg/postgres"
)

const (
	userTable        = "users"
	userSessionTable = "users_sessions"
	codesTable       = "codes"
)

type Repository struct {
	Auth
}

func NewPsqlRepository(db *postgres.Postgres) *Repository {
	return &Repository{
		Auth: NewAuthPostgres(db),
	}
}
