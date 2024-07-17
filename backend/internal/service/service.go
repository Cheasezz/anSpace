package service

import (
	repositories "github.com/Cheasezz/anSpace/backend/internal/repository"
	"github.com/Cheasezz/anSpace/backend/pkg/auth"
	"github.com/Cheasezz/anSpace/backend/pkg/hasher"
)

type Services struct {
	Auth
}

type Deps struct {
	Repos        *repositories.Repositories
	Hasher       hasher.PasswordHasher
	TokenManager auth.TokenManager
}

func NewServices(d Deps) *Services {
	return &Services{
		Auth: newAuthService(d.Repos.Psql.Auth, d.Hasher, d.TokenManager),
	}
}
