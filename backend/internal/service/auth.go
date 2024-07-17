package service

import (
	"context"
	"time"

	"github.com/Cheasezz/anSpace/backend/internal/core"
	"github.com/Cheasezz/anSpace/backend/internal/repository/psql"
	"github.com/Cheasezz/anSpace/backend/pkg/auth"
	"github.com/Cheasezz/anSpace/backend/pkg/hasher"
)

type Auth interface {
	SignUp(ctx context.Context, signUp core.SignUp) (string, error)
	SignIn(ctx context.Context, signIn core.SignIn) (auth.Tokens, error)
	LogOut(ctx context.Context, refreshToken string) (auth.Tokens, error)
	RefreshAccessToken(ctx context.Context, refreshToken string) (auth.Tokens, error)
}

var (
	daysForUpdRtToken = 5
)

type AuthService struct {
	repo         psql.Auth
	hasher       hasher.PasswordHasher
	tokenManager auth.TokenManager
}

func newAuthService(r psql.Auth, h hasher.PasswordHasher, tm auth.TokenManager) *AuthService {
	return &AuthService{
		repo:         r,
		hasher:       h,
		tokenManager: tm,
	}
}

// Hash password and write new user into db.
// With method repo.CreateUser.
// Return userId and error.
func (s *AuthService) SignUp(ctx context.Context, signUp core.SignUp) (string, error) {
	pass, err := s.hasher.Hash(signUp.Password)
	if err != nil {
		return "", err
	}
	signUp.Password = pass
	userId, err := s.repo.CreateUser(ctx, signUp)
	return userId, err
}

// Hash password and search ind db userId with method repo.GetUser.
// Pass userId into createSession method.
// Return auth.Tokens and error.
func (s *AuthService) SignIn(ctx context.Context, signIn core.SignIn) (auth.Tokens, error) {
	pass, err := s.hasher.Hash(signIn.Password)
	if err != nil {
		return auth.Tokens{}, err
	}
	signIn.Password = pass
	userId, err := s.repo.GetUser(ctx, signIn)
	if err != nil {
		return auth.Tokens{}, err
	}

	return s.createSession(ctx, userId)
}

// Uppdate core.Session in repo with emty tokens value.
// Return empty auth.Tokens struct
func (s *AuthService) LogOut(ctx context.Context, refreshToken string) (auth.Tokens, error) {

	tkns := auth.Tokens{Access: "", Refresh: auth.RTInfo{Token: "", ExpiresAt: time.Now(), TTLInSec: 0}}
	session, err := s.repo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return tkns, err
	}
	session.RefreshToken = tkns.Refresh.Token
	session.ExpiresAt = tkns.Refresh.ExpiresAt

	s.repo.SetSession(ctx, session)
	return tkns, nil
}

// Generate new jwt and refresh tokens, create new core.Session struct.
// And write this session in repo by repo.SetSession method.
// Return auth.Tokens and error.
func (s *AuthService) createSession(ctx context.Context, userId string) (auth.Tokens, error) {
	var (
		tokens auth.Tokens
		err    error
	)

	tokens.Access, err = s.tokenManager.NewJWT(userId)
	if err != nil {
		return auth.Tokens{}, err
	}

	RTInfo, err := s.tokenManager.NewRefreshToken()
	if err != nil {
		return auth.Tokens{}, err
	}
	tokens.Refresh = RTInfo

	session := core.Session{
		UserId:       userId,
		RefreshToken: tokens.Refresh.Token,
		ExpiresAt:    tokens.Refresh.ExpiresAt,
	}

	err = s.repo.SetSession(ctx, session)
	if err != nil {
		return auth.Tokens{}, err
	}

	return tokens, nil
}

// If days until refresh token expired is equal or less then var daysForUpdRtToken
// return auth.Tokens struct with uppdated access token and refresh token.
// Else return auth.Tokens struct with uppdate access token only.
func (s *AuthService) RefreshAccessToken(ctx context.Context, refreshToken string) (auth.Tokens, error) {
	var (
		tokens auth.Tokens
	)

	session, err := s.repo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return tokens, err
	}

	rtDayUntilExpire, err := s.tokenManager.ValidateRefreshToken(session.ExpiresAt)
	if err != nil {
		return auth.Tokens{}, err
	}

	if rtDayUntilExpire <= daysForUpdRtToken {
		return s.createSession(ctx, session.UserId)
	}
	tokens.Access, err = s.tokenManager.NewJWT(session.UserId)
	if err != nil {
		return auth.Tokens{}, err
	}

	return tokens, nil
}
