package psql

import (
	"context"
	"fmt"

	"github.com/Cheasezz/anSpace/backend/internal/core"
	"github.com/Cheasezz/anSpace/backend/pkg/postgres"
	"github.com/google/uuid"
)

type Auth interface {
	CreateUser(ctx context.Context, signUp core.AuthCredentials) (uuid.UUID, error)
	GetUserIdByLogPas(ctx context.Context, signIn core.AuthCredentials) (uuid.UUID, error)
	GetUserById(ctx context.Context, userId uuid.UUID) (core.User, error)
	GetUserByEmail(ctx context.Context, email string) (core.User, error)
	SetPassResetCode(ctx context.Context, code core.CodeCredentials) error
	DeletePassResetCode(ctx context.Context, code core.CodeCredentials) error
	SetSession(ctx context.Context, session core.Session) error
	DeleteSession(ctx context.Context, session core.Session) error
	GetUserSessionByRefreshToken(ctx context.Context, refreshToken string) (core.Session, error)
}

type AuthRepo struct {
	db *postgres.Postgres
}

func NewAuthPostgres(db *postgres.Postgres) *AuthRepo {
	return &AuthRepo{db: db}
}

func (r *AuthRepo) CreateUser(ctx context.Context, signUp core.AuthCredentials) (uuid.UUID, error) {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return uuid.UUID{}, err
	}

	var id uuid.UUID

	query := fmt.Sprintf("INSERT INTO %s (email, password_hash) values ($1, $2) RETURNING id", userTable)
	row := tx.QueryRow(ctx, query, signUp.Email, signUp.Password)
	if err := row.Scan(&id); err != nil {
		if errR := tx.Rollback(ctx); errR != nil {
			return uuid.UUID{}, errR
		}
		return uuid.UUID{}, err
	}

	return id, tx.Commit(ctx)
}

func (r *AuthRepo) GetUserIdByLogPas(ctx context.Context, signIn core.AuthCredentials) (uuid.UUID, error) {
	var userId uuid.UUID
	query := fmt.Sprintf("SELECT id FROM %s WHERE email=$1 AND password_hash=$2", userTable)
	err := r.db.Scany.Get(ctx, r.db.Pool, &userId, query, signIn.Email, signIn.Password)

	return userId, err
}

func (r *AuthRepo) GetUserById(ctx context.Context, userId uuid.UUID) (core.User, error) {
	var user core.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", userTable)
	err := r.db.Scany.Get(ctx, r.db.Pool, &user, query, userId)

	return user, err
}

func (r *AuthRepo) GetUserByEmail(ctx context.Context, email string) (core.User, error) {
	var user core.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE email=$1", userTable)
	err := r.db.Scany.Get(ctx, r.db.Pool, &user, query, email)

	return user, err
}

func (r *AuthRepo) SetPassResetCode(ctx context.Context, code core.CodeCredentials) error {
	query := fmt.Sprintf("INSERT INTO %s (user_email, code, code_type, expires_at) values ($1, $2, $3, $4)", codesTable)
	_, err := r.db.Pool.Exec(ctx, query, code.Email, code.Code, code.CodeType, code.ExpiresAt)

	return err
}

func (r *AuthRepo) DeletePassResetCode(ctx context.Context, code core.CodeCredentials) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE code=$1 AND user_email=$2", codesTable)
	_, err := r.db.Pool.Exec(ctx, query, code.Code, code.Email)

	return err
}

func (r *AuthRepo) SetSession(ctx context.Context, session core.Session) error {
	query := fmt.Sprintf("INSERT INTO  %s (user_id,refresh_token) values ($1, $2)", userSessionTable)
	_, err := r.db.Pool.Exec(ctx, query, session.UserId, session.RefreshToken)

	return err
}
func (r *AuthRepo) DeleteSession(ctx context.Context, session core.Session) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE refresh_token = $1 AND user_id = $2", userSessionTable)
	_, err := r.db.Pool.Exec(ctx, query, session.RefreshToken, session.UserId)

	return err
}

func (r *AuthRepo) GetUserSessionByRefreshToken(ctx context.Context, refreshToken string) (core.Session, error) {
	var session core.Session

	query := fmt.Sprintf("SELECT user_id, refresh_token FROM %s WHERE refresh_token=$1", userSessionTable)
	err := r.db.Scany.Get(ctx, r.db.Pool, &session, query, refreshToken)

	return session, err
}
