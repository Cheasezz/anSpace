package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Cheasezz/anSpace/backend/internal/core"
	mock_psql "github.com/Cheasezz/anSpace/backend/internal/repository/psql/mocks"
	"github.com/Cheasezz/anSpace/backend/pkg/auth"
	mock_auth "github.com/Cheasezz/anSpace/backend/pkg/auth/mocks"
	mock_hash "github.com/Cheasezz/anSpace/backend/pkg/hasher/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var (
	errHasher                = fmt.Errorf("hasher error")
	errRepo                  = fmt.Errorf("repo error")
	errNewJwt                = fmt.Errorf("new jwt error")
	errNewRefreshToken       = fmt.Errorf("new refresh token error")
	errRepoGetByRefreshToken = fmt.Errorf("repo get by refresh token error")
	errRefreshTokenIsExpired = fmt.Errorf("tm refresh token is expired")
)

func TestAuth_SignUp(t *testing.T) {
	type mockBehavior func(h *mock_hash.MockPasswordHasher, r *mock_psql.MockAuth, input core.SignUp)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_psql.NewMockAuth(ctrl)
	hash := mock_hash.NewMockPasswordHasher(ctrl)
	tm := mock_auth.NewMockTokenManager(ctrl)

	authSrv := newAuthService(repo, hash, tm)

	tests := []struct {
		name         string
		inputUser    core.SignUp
		userId       string
		expErr       error
		mockBehavior mockBehavior
	}{
		{
			name:      "OK",
			inputUser: core.SignUp{Name: "Iurii", Username: "Cheasezz", Password: "qwerty123456"},
			userId:    "1",
			mockBehavior: func(h *mock_hash.MockPasswordHasher, r *mock_psql.MockAuth, input core.SignUp) {
				h.EXPECT().Hash(input.Password).Return(input.Password, nil)
				r.EXPECT().CreateUser(gomock.Any(), input).Return("1", nil)
			},
		},
		{
			name:      "Hasher error",
			inputUser: core.SignUp{Name: "Iurii", Username: "Cheasezz", Password: "qwerty123456"},
			userId:    "",
			expErr:    errHasher,
			mockBehavior: func(h *mock_hash.MockPasswordHasher, r *mock_psql.MockAuth, input core.SignUp) {
				h.EXPECT().Hash(input.Password).Return("", errHasher)
			},
		},
		{
			name:      "Repo error",
			inputUser: core.SignUp{Name: "Iurii", Username: "Cheasezz", Password: "qwerty123456"},
			userId:    "",
			expErr:    errRepo,
			mockBehavior: func(h *mock_hash.MockPasswordHasher, r *mock_psql.MockAuth, input core.SignUp) {
				h.EXPECT().Hash(input.Password).Return(input.Password, nil)
				r.EXPECT().CreateUser(gomock.Any(), input).Return("", errRepo)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(hash, repo, tt.inputUser)

			id, err := authSrv.SignUp(context.Background(), tt.inputUser)
			if err != nil {
				require.Empty(t, id)
				require.Error(t, err)
				require.EqualError(t, tt.expErr, err.Error())
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.userId, id)

		})
	}
}

func initTokens() auth.Tokens {
	tokenStr := "token"

	return auth.Tokens{
		Access: tokenStr,
		Refresh: auth.RTInfo{
			Token:     tokenStr,
			ExpiresAt: time.Now().Add(time.Hour),
			TTLInSec:  43000},
	}
}
func initSession() core.Session {
	return core.Session{UserId: "1", RefreshToken: "token", ExpiresAt: time.Now().Add(time.Hour)}
}

type deps struct {
	h  *mock_hash.MockPasswordHasher
	r  *mock_psql.MockAuth
	tm *mock_auth.MockTokenManager
}

func initDeps(h *mock_hash.MockPasswordHasher, r *mock_psql.MockAuth, tm *mock_auth.MockTokenManager) deps {
	return deps{h, r, tm}
}
func TestAuth_SignIn(t *testing.T) {
	type mockBehavior func(d deps, i core.SignIn, s core.Session)
	tokens := initTokens()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hash := mock_hash.NewMockPasswordHasher(ctrl)
	repo := mock_psql.NewMockAuth(ctrl)
	tm := mock_auth.NewMockTokenManager(ctrl)
	d := initDeps(hash, repo, tm)

	authSrv := newAuthService(repo, hash, tm)

	tests := []struct {
		name         string
		inputUser    core.SignIn
		session      core.Session
		expErr       error
		mockBehavior mockBehavior
	}{
		{
			name:      "OK",
			inputUser: core.SignIn{Username: "Cheasezz", Password: "qwerty123456"},
			session:   core.Session{UserId: "1", RefreshToken: tokens.Refresh.Token, ExpiresAt: tokens.Refresh.ExpiresAt},
			mockBehavior: func(d deps, i core.SignIn, s core.Session) {
				d.h.EXPECT().Hash(i.Password).Return(i.Password, nil)
				d.r.EXPECT().GetUser(gomock.Any(), i).Return("1", nil)
				d.tm.EXPECT().NewJWT("1").Return(tokens.Access, nil)
				d.tm.EXPECT().NewRefreshToken().Return(tokens.Refresh, nil)
				d.r.EXPECT().SetSession(gomock.Any(), s).Return(nil)
			},
		},
		{
			name:      "hasher error",
			inputUser: core.SignIn{Username: "Cheasezz", Password: "qwerty123456"},
			expErr:    errHasher,
			mockBehavior: func(d deps, i core.SignIn, s core.Session) {
				d.h.EXPECT().Hash(i.Password).Return("", errHasher)
			},
		},
		{
			name:      "repo get user error",
			inputUser: core.SignIn{Username: "Cheasezz", Password: "qwerty123456"},
			expErr:    errRepo,
			mockBehavior: func(d deps, i core.SignIn, s core.Session) {
				d.h.EXPECT().Hash(i.Password).Return(i.Password, nil)
				d.r.EXPECT().GetUser(gomock.Any(), i).Return("", errRepo)

			},
		},
		{
			name:      "tm new jwt error",
			inputUser: core.SignIn{Username: "Cheasezz", Password: "qwerty123456"},
			expErr:    errNewJwt,
			mockBehavior: func(d deps, i core.SignIn, s core.Session) {
				d.h.EXPECT().Hash(i.Password).Return(i.Password, nil)
				d.r.EXPECT().GetUser(gomock.Any(), i).Return("1", nil)
				d.tm.EXPECT().NewJWT("1").Return("", errNewJwt)
			},
		},
		{
			name:      "tm new refresh token error",
			inputUser: core.SignIn{Username: "Cheasezz", Password: "qwerty123456"},
			expErr:    errNewRefreshToken,
			mockBehavior: func(d deps, i core.SignIn, s core.Session) {
				d.h.EXPECT().Hash(i.Password).Return(i.Password, nil)
				d.r.EXPECT().GetUser(gomock.Any(), i).Return("1", nil)
				d.tm.EXPECT().NewJWT("1").Return(tokens.Access, nil)
				d.tm.EXPECT().NewRefreshToken().Return(auth.RTInfo{}, errNewRefreshToken)
			},
		},
		{
			name:      "repo set session error",
			inputUser: core.SignIn{Username: "Cheasezz", Password: "qwerty123456"},
			session:   core.Session{UserId: "1", RefreshToken: tokens.Refresh.Token, ExpiresAt: tokens.Refresh.ExpiresAt},
			expErr:    errRepo,
			mockBehavior: func(d deps, i core.SignIn, s core.Session) {
				d.h.EXPECT().Hash(i.Password).Return(i.Password, nil)
				d.r.EXPECT().GetUser(gomock.Any(), i).Return("1", nil)
				d.tm.EXPECT().NewJWT("1").Return(tokens.Access, nil)
				d.tm.EXPECT().NewRefreshToken().Return(tokens.Refresh, nil)
				d.r.EXPECT().SetSession(gomock.Any(), s).Return(errRepo)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(d, tt.inputUser, tt.session)

			tkns, err := authSrv.SignIn(context.Background(), tt.inputUser)
			if err != nil {
				require.Empty(t, tkns)
				require.Error(t, err)
				require.EqualError(t, tt.expErr, err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tokens, tkns)
			}
		})
	}
}

func TestAuth_LogOut(t *testing.T) {
	type mockBehavior func(d deps, rt string, s core.Session)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hash := mock_hash.NewMockPasswordHasher(ctrl)
	repo := mock_psql.NewMockAuth(ctrl)
	tm := mock_auth.NewMockTokenManager(ctrl)
	d := initDeps(hash, repo, tm)

	authSrv := newAuthService(repo, hash, tm)

	tests := []struct {
		name         string
		rToken       string
		session      core.Session
		expErr       error
		expTokens    auth.Tokens
		mockBehavior mockBehavior
	}{
		{
			name:      "ok",
			rToken:    "token",
			session:   initSession(),
			expTokens: auth.Tokens{Refresh: auth.RTInfo{ExpiresAt: time.Now()}},
			mockBehavior: func(d deps, rt string, s core.Session) {
				d.r.EXPECT().GetByRefreshToken(gomock.Any(), rt).Return(s, nil)
				d.r.EXPECT().SetSession(gomock.Any(), gomock.Any())
			},
		},
		{
			name:      "repo get by refresh token error",
			rToken:    "token",
			session:   core.Session{},
			expTokens: auth.Tokens{Refresh: auth.RTInfo{ExpiresAt: time.Now()}},
			expErr:    errRepoGetByRefreshToken,
			mockBehavior: func(d deps, rt string, s core.Session) {
				d.r.EXPECT().GetByRefreshToken(gomock.Any(), rt).Return(s, errRepoGetByRefreshToken)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(d, tt.rToken, tt.session)

			tkns, err := authSrv.LogOut(context.Background(), tt.rToken)
			if err != nil {
				require.Error(t, err)
				require.EqualError(t, tt.expErr, err.Error())
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.expTokens.Access, tkns.Access)
			require.Equal(t, tt.expTokens.Refresh.Token, tkns.Refresh.Token)
		})
	}
}

func TestAuthService_RefreshAccessToken(t *testing.T) {
	type mockBehavior func(d deps, tkns auth.Tokens, s core.Session, day int)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hash := mock_hash.NewMockPasswordHasher(ctrl)
	repo := mock_psql.NewMockAuth(ctrl)
	tm := mock_auth.NewMockTokenManager(ctrl)
	d := initDeps(hash, repo, tm)

	authSrv := newAuthService(repo, hash, tm)

	tests := []struct {
		name           string
		tokens         auth.Tokens
		session        core.Session
		dayUntilExpire int
		expErr         error
		mockBehavior   mockBehavior
	}{
		{
			name:           "ok (upd access token only)",
			tokens:         auth.Tokens{Access: "token"},
			session:        initSession(),
			dayUntilExpire: 20,
			mockBehavior: func(d deps, tkns auth.Tokens, s core.Session, day int) {
				d.r.EXPECT().GetByRefreshToken(gomock.Any(), tkns.Refresh.Token).Return(s, nil)
				d.tm.EXPECT().ValidateRefreshToken(s.ExpiresAt).Return(day, nil)
				d.tm.EXPECT().NewJWT(s.UserId).Return(tkns.Access, nil)
			},
		},
		{
			name:           "ok (upd both tokens)",
			tokens:         initTokens(),
			session:        initSession(),
			dayUntilExpire: 5,
			mockBehavior: func(d deps, tkns auth.Tokens, s core.Session, day int) {
				d.r.EXPECT().GetByRefreshToken(gomock.Any(), tkns.Refresh.Token).Return(s, nil)
				d.tm.EXPECT().ValidateRefreshToken(s.ExpiresAt).Return(day, nil)
				d.tm.EXPECT().NewJWT(s.UserId).Return(tkns.Access, nil)
				d.tm.EXPECT().NewRefreshToken().Return(tkns.Refresh, nil)
				d.r.EXPECT().SetSession(gomock.Any(), s).Return(nil)
			},
		},
		{
			name:    "repo get by refresh token error",
			tokens:  initTokens(),
			session: core.Session{},
			expErr:  errRepoGetByRefreshToken,
			mockBehavior: func(d deps, tkns auth.Tokens, s core.Session, day int) {
				d.r.EXPECT().GetByRefreshToken(gomock.Any(), tkns.Refresh.Token).Return(s, errRepoGetByRefreshToken)
			},
		},
		{
			name:    "tm refresh token is expired",
			tokens:  initTokens(),
			session: initSession(),
			expErr:  errRefreshTokenIsExpired,
			mockBehavior: func(d deps, tkns auth.Tokens, s core.Session, day int) {
				d.r.EXPECT().GetByRefreshToken(gomock.Any(), tkns.Refresh.Token).Return(s, nil)
				d.tm.EXPECT().ValidateRefreshToken(s.ExpiresAt).Return(0, errRefreshTokenIsExpired)
			},
		},
		{
			name:           "tm new jwt error",
			tokens:         initTokens(),
			session:        initSession(),
			dayUntilExpire: 20,
			expErr:         errNewJwt,
			mockBehavior: func(d deps, tkns auth.Tokens, s core.Session, day int) {
				d.r.EXPECT().GetByRefreshToken(gomock.Any(), tkns.Refresh.Token).Return(s, nil)
				d.tm.EXPECT().ValidateRefreshToken(s.ExpiresAt).Return(day, nil)
				d.tm.EXPECT().NewJWT(s.UserId).Return("", errNewJwt)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(d, tt.tokens, tt.session, tt.dayUntilExpire)
			tkns, err := authSrv.RefreshAccessToken(context.Background(), tt.tokens.Refresh.Token)
			if err != nil {
				require.Empty(t, tkns)
				require.Error(t, err)
				require.EqualError(t, tt.expErr, err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.tokens, tkns)
			}
		})
	}
}
