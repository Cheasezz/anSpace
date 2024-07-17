package repositories

import (
	"github.com/Cheasezz/anSpace/backend/internal/repository/psql"
	"github.com/Cheasezz/anSpace/backend/pkg/postgres"
)

type Repositories struct {
	Psql *psql.Repository
}

func NewRepositories(pg *postgres.Postgres) *Repositories {
	return &Repositories{
		Psql: psql.NewPsqlRepository(pg),
	}
}
