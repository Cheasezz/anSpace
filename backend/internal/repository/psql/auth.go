package psql

import (
	"context"
	"fmt"

	"github.com/Cheasezz/anSpace/backend/internal/core"
	"github.com/Cheasezz/anSpace/backend/pkg/postgres"
)

type Auth interface {
	CreateUser(ctx context.Context, signUp core.SignUp) (string, error)
	GetUser(ctx context.Context, signIn core.SignIn) (string, error)
	SetSession(ctx context.Context, session core.Session) error
	GetByRefreshToken(ctx context.Context, refreshToken string) (core.Session, error)
}

type AuthRepo struct {
	db *postgres.Postgres
}

func NewAuthPostgres(db *postgres.Postgres) *AuthRepo {
	return &AuthRepo{db: db}
}

func (r *AuthRepo) CreateUser(ctx context.Context, signUp core.SignUp) (string, error) {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return "", err
	}

	var id string

	query := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) values ($1, $2, $3) RETURNING id", userTable)
	row := tx.QueryRow(ctx, query, signUp.Name, signUp.Username, signUp.Password)
	if err := row.Scan(&id); err != nil {
		if errR := tx.Rollback(ctx); errR != nil {
			return "", errR
		}
		return "", err
	}

	createUserSessionQuery := fmt.Sprintf("INSERT INTO %s (user_id) values ($1)", userSessionTable)
	_, err = tx.Exec(ctx, createUserSessionQuery, id)
	if err != nil {
		if errR := tx.Rollback(ctx); errR != nil {
			return "", errR
		}
		return "", err
	}

	return id, tx.Commit(ctx)
}

func (r *AuthRepo) GetUser(ctx context.Context, signIn core.SignIn) (string, error) {
	var userId string
	query := fmt.Sprintf("SELECT id FROM %s WHERE username=$1 AND password_hash=$2", userTable)
	err := r.db.Scany.Get(ctx, r.db.Pool, &userId, query, signIn.Username, signIn.Password)

	return userId, err
}

func (r *AuthRepo) SetSession(ctx context.Context, session core.Session) error {
	query := fmt.Sprintf("UPDATE %s us SET (refresh_token, expires_at) = ($1, $2) WHERE user_id = $3", userSessionTable)
	_, err := r.db.Pool.Exec(ctx, query, session.RefreshToken, session.ExpiresAt, session.UserId)

	return err
}

func (r *AuthRepo) GetByRefreshToken(ctx context.Context, refreshToken string) (core.Session, error) {
	var session core.Session

	query := fmt.Sprintf("SELECT user_id, expires_at, refresh_token FROM %s WHERE refresh_token=$1", userSessionTable)
	err := r.db.Scany.Get(ctx, r.db.Pool, &session, query, refreshToken)

	return session, err
}
