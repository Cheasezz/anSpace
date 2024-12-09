package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/Cheasezz/anSpace/backend/internal/core"
	"github.com/Cheasezz/anSpace/backend/internal/repository/psql"
	"github.com/Cheasezz/anSpace/backend/pkg/auth"
	"github.com/Cheasezz/anSpace/backend/pkg/email"
	"github.com/Cheasezz/anSpace/backend/pkg/hasher"
	"github.com/google/uuid"
)

type Auth interface {
	SignUp(ctx context.Context, signUp core.AuthCredentials) (auth.Tokens, error)
	SignIn(ctx context.Context, signIn core.AuthCredentials) (auth.Tokens, error)
	LogOut(ctx context.Context, refreshToken string) (auth.Tokens, error)
	RefreshAccessToken(ctx context.Context, refreshToken string) (auth.Tokens, error)
	GetUser(ctx context.Context, userId uuid.UUID) (core.User, error)
	GenPassResetCode(ctx context.Context, email string) error
}

type AuthService struct {
	repo         psql.Auth
	hasher       hasher.PasswordHasher
	tokenManager auth.TokenManager
	emailSender  email.Sender
}

func newAuthService(r psql.Auth, h hasher.PasswordHasher, tm auth.TokenManager, es email.Sender) *AuthService {
	return &AuthService{
		repo:         r,
		hasher:       h,
		tokenManager: tm,
		emailSender:  es,
	}
}

// Hash password and write new user into db.
// With method repo.CreateUser.
// Return userId and error.
func (s *AuthService) SignUp(ctx context.Context, signUp core.AuthCredentials) (auth.Tokens, error) {
	pass, err := s.hasher.Hash(signUp.Password)
	if err != nil {
		return auth.Tokens{}, err
	}
	signUp.Password = pass
	userId, err := s.repo.CreateUser(ctx, signUp)
	if err != nil {
		return auth.Tokens{}, err
	}

	return s.createSession(ctx, userId)
}

// Hash password and search ind db userId with method repo.GetUser.
// Pass userId into createSession method.
// Return auth.Tokens and error.
func (s *AuthService) SignIn(ctx context.Context, signIn core.AuthCredentials) (auth.Tokens, error) {
	pass, err := s.hasher.Hash(signIn.Password)
	if err != nil {
		return auth.Tokens{}, err
	}
	signIn.Password = pass
	userId, err := s.repo.GetUserIdByLogPas(ctx, signIn)
	if err != nil {
		return auth.Tokens{}, err
	}

	return s.createSession(ctx, userId)
}

// Uppdate core.Session in repo with emty tokens value.
// Return empty auth.Tokens struct
func (s *AuthService) LogOut(ctx context.Context, refreshToken string) (auth.Tokens, error) {

	tkns := auth.Tokens{Access: auth.ATknInfo{Token: ""}, Refresh: auth.RTknInfo{Token: ""}}
	session, err := s.repo.GetUserSessionByRefreshToken(ctx, refreshToken)
	if err != nil {
		return tkns, err
	}
	session.RefreshToken = tkns.Refresh.Token

	if err = s.repo.SetSession(ctx, session); err != nil {
		return auth.Tokens{}, err
	}
	return tkns, nil
}

// Generate new jwt and refresh tokens, create new core.Session struct.
// And write this session in repo by repo.SetSession method.
// Return auth.Tokens and error.
func (s *AuthService) createSession(ctx context.Context, userId uuid.UUID) (auth.Tokens, error) {
	var (
		tokens auth.Tokens
		err    error
	)

	tokens.Access.Token, err = s.tokenManager.NewJWT(userId.String())
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
	}

	err = s.repo.SetSession(ctx, session)
	if err != nil {
		return auth.Tokens{}, err
	}

	return tokens, nil
}

// Return auth.Tokens struct with uppdate access token only.
func (s *AuthService) RefreshAccessToken(ctx context.Context, refreshToken string) (auth.Tokens, error) {
	var (
		tokens auth.Tokens
	)

	session, err := s.repo.GetUserSessionByRefreshToken(ctx, refreshToken)
	if err != nil {
		return tokens, err
	}

	tokens.Access.Token, err = s.tokenManager.NewJWT(session.UserId.String())
	if err != nil {
		return auth.Tokens{}, err
	}

	return tokens, nil
}

// Return user by userid
func (s *AuthService) GetUser(c context.Context, userId uuid.UUID) (core.User, error) {
	user, err := s.repo.GetUserById(c, userId)
	if err != nil {
		return core.User{}, err
	}
	return user, nil
}

// Genereate password reset code. Set them in db and send on email in param
func (s *AuthService) GenPassResetCode(c context.Context, email string) error {
	user, err := s.repo.GetUserByEmail(c, email)
	if err != nil {
		return err
	}

	passResetCode := make([]byte, 32)
	src := rand.NewSource(time.Now().UTC().Unix() + rand.Int63())
	r := rand.New(src)
	r.Read(passResetCode)

	code := core.CodeCredentials{
		Email:     email,
		Code:      fmt.Sprintf("%x", passResetCode),
		CodeType:  "passReset",
		ExpiresAt: time.Now().Add(time.Minute * 30),
	}

	if err := s.repo.SetPassResetCode(c, code); err != nil {
		return err
	}
	message := fmt.Sprintf("This is your password reset code:%x", passResetCode)
	if err := s.emailSender.Send(user.Email, message); err != nil {
		return err
	}
	return nil
}
