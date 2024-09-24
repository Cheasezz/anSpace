package v1

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/Cheasezz/anSpace/backend/config"
	"github.com/Cheasezz/anSpace/backend/internal/core"
	"github.com/Cheasezz/anSpace/backend/internal/service"
	mock_service "github.com/Cheasezz/anSpace/backend/internal/service/mocks"
	"github.com/Cheasezz/anSpace/backend/pkg/auth"
	mock_auth "github.com/Cheasezz/anSpace/backend/pkg/auth/mocks"
	"github.com/Cheasezz/anSpace/backend/pkg/logger"
	mock_logger "github.com/Cheasezz/anSpace/backend/pkg/logger/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var (
	tokens auth.Tokens

	errServiceSignUp             = fmt.Errorf("service sign up error")
	errServiceSignIn             = fmt.Errorf("service sign in error")
	errServiceLogOut             = fmt.Errorf("service logout error")
	errServiceRefreshAccessToken = fmt.Errorf("service refreshAccessToken error")
	errServiceGetUser            = fmt.Errorf("service getUser error")
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	tokens = initTokens()

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestAuthHandler_validateLoginAndPass(t *testing.T) {
	type args struct {
		l string
		p string
	}
	tests := []struct {
		name   string
		h      *Auth
		args   args
		expErr error
	}{
		{
			name: "correct login and password",
			args: args{
				l: "kappa",
				p: "qwerty123456",
			},
		},
		{
			name: "empty login after trim",
			args: args{
				l: " ",
				p: "qwerty123456",
			},
			expErr: errEmptyLoginOrPass,
		},
		{
			name: "empty password after trim",
			args: args{
				l: "kappa",
				p: " ",
			},
			expErr: errEmptyLoginOrPass,
		},
		{
			name: "short password",
			args: args{
				l: "kappa",
				p: "qwerty",
			},
			expErr: errShortPass,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.h.validateLoginAndPass(tt.args.l, tt.args.p)

			if tt.expErr == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.EqualError(t, tt.expErr, err.Error())
			}
		})

	}
}

func initDeps(s *service.Services, tm auth.TokenManager, l logger.Logger) Deps {
	return Deps{
		Services:     s,
		TokenManager: tm,
		ConfigHTTP:   config.HTTP{Host: "127.0.0.1", Port: "8000"},
		Log:          l,
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

func TestAuthHandler_signUp(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuth, l *mock_logger.MockLogger, input core.SignUp, signIn core.SignIn)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authSrv := mock_service.NewMockAuth(ctrl)
	tm := mock_auth.NewMockTokenManager(ctrl)
	l := mock_logger.NewMockLogger(ctrl)

	services := &service.Services{Auth: authSrv}
	deps := initDeps(services, tm, l)
	handler := NewAuthHandler(deps)

	r := gin.New()
	v1 := r.Group("/v1")
	handler.initAuthRoutes(v1)

	tests := []struct {
		name         string
		mockBehavior mockBehavior
		inputBody    string
		SignUp       core.SignUp
		SignIn       core.SignIn
		expStatCode  int
		expReqBody   string
	}{
		{
			name:      "OK",
			inputBody: `{"name":"Iurii","username": "Cheasezz","password":"qwerty123456"}`,
			SignUp:    core.SignUp{Name: "Iurii", Username: "Cheasezz", Password: "qwerty123456"},
			SignIn:    core.SignIn{Username: "Cheasezz", Password: "qwerty123456"},
			mockBehavior: func(s *mock_service.MockAuth, l *mock_logger.MockLogger, signUp core.SignUp, signIn core.SignIn) {
				s.EXPECT().SignUp(gomock.Any(), signUp).Return("fwfgq134", nil)
				s.EXPECT().SignIn(gomock.Any(), signIn).Return(tokens, nil)
			},
			expStatCode: 200,
			expReqBody:  fmt.Sprintf(`{"accessToken":"%s"}`, tokens.Access),
		},
		{
			name:      "Bad request: empty username",
			inputBody: `{"name":"Iurii","username": "","password":"qwerty123456"}`,
			SignUp:    core.SignUp{Name: "Iurii", Username: "", Password: "qwerty123456"},
			mockBehavior: func(s *mock_service.MockAuth, l *mock_logger.MockLogger, signUp core.SignUp, signIn core.SignIn) {
				l.EXPECT().Error(gomock.Any())
			},
			expStatCode: 400,
			expReqBody:  `{"message":"Key: 'SignUp.Username' Error:Field validation for 'Username' failed on the 'required' tag"}`,
		},
		{
			name:      "Bad request: empty name",
			inputBody: `{"name":"","username": "Cheasezz","password":"qwerty123456"}`,
			SignUp:    core.SignUp{Name: "", Username: "Cheasezz", Password: "qwerty123456"},
			mockBehavior: func(s *mock_service.MockAuth, l *mock_logger.MockLogger, signUp core.SignUp, signIn core.SignIn) {
				l.EXPECT().Error(gomock.Any())
			},
			expStatCode: 400,
			expReqBody:  `{"message":"Key: 'SignUp.Name' Error:Field validation for 'Name' failed on the 'required' tag"}`,
		},
		{
			name:      "Bad request: empty password",
			inputBody: `{"name":"Iurii","username": "Cheasezz","password":""}`,
			SignUp:    core.SignUp{Name: "Iurii", Username: "Cheasezz", Password: ""},
			mockBehavior: func(s *mock_service.MockAuth, l *mock_logger.MockLogger, signUp core.SignUp, signIn core.SignIn) {
				l.EXPECT().Error(gomock.Any())
			},
			expStatCode: 400,
			expReqBody:  `{"message":"Key: 'SignUp.Password' Error:Field validation for 'Password' failed on the 'required' tag"}`,
		},
		{
			name:      "Bad request: empty login after trim",
			inputBody: `{"name":"Iurii","username":" ","password":"qwerty123456"}`,
			SignUp:    core.SignUp{Name: "Iurii", Username: " ", Password: "qwerty123456"},
			mockBehavior: func(s *mock_service.MockAuth, l *mock_logger.MockLogger, signUp core.SignUp, signIn core.SignIn) {
				l.EXPECT().Error(errEmptyLoginOrPass)
			},
			expStatCode: 400,
			expReqBody:  fmt.Sprintf(`{"message":"%s"}`, errEmptyLoginOrPass),
		},
		{
			name:      "Bad request: empty passwrod after trim",
			inputBody: `{"name":"Iurii","username":"Cheasezz","password":" "}`,
			SignUp:    core.SignUp{Name: "Iurii", Username: "Cheasezz", Password: " "},
			mockBehavior: func(s *mock_service.MockAuth, l *mock_logger.MockLogger, signUp core.SignUp, signIn core.SignIn) {
				l.EXPECT().Error(errEmptyLoginOrPass)
			},
			expStatCode: 400,
			expReqBody:  fmt.Sprintf(`{"message":"%s"}`, errEmptyLoginOrPass),
		},
		{
			name:      "Bad request: password less then 12 char",
			inputBody: `{"name":"Iurii","username":"Cheasezz","password":"qwerty12345"}`,
			SignUp:    core.SignUp{Name: "Iurii", Username: "Cheasezz", Password: "qwerty12345"},
			mockBehavior: func(s *mock_service.MockAuth, l *mock_logger.MockLogger, signUp core.SignUp, signIn core.SignIn) {
				l.EXPECT().Error(errShortPass)
			},
			expStatCode: 400,
			expReqBody:  fmt.Sprintf(`{"message":"%s"}`, errShortPass),
		},
		{
			name:      "Server error: Service Sign Up error",
			inputBody: `{"name":"Iurii","username": "Cheasezz","password":"qwerty123456"}`,
			SignUp:    core.SignUp{Name: "Iurii", Username: "Cheasezz", Password: "qwerty123456"},
			mockBehavior: func(s *mock_service.MockAuth, l *mock_logger.MockLogger, signUp core.SignUp, signIn core.SignIn) {
				s.EXPECT().SignUp(gomock.Any(), signUp).Return(" ", errServiceSignUp)
				l.EXPECT().Error(errServiceSignUp)
			},
			expStatCode: 500,
			expReqBody:  fmt.Sprintf(`{"message":"%s"}`, errServiceSignUp),
		},
		{
			name:      "Server error: Service Sign In error",
			inputBody: `{"name":"Iurii","username": "Cheasezz","password":"qwerty123456"}`,
			SignUp:    core.SignUp{Name: "Iurii", Username: "Cheasezz", Password: "qwerty123456"},
			SignIn:    core.SignIn{Username: "Cheasezz", Password: "qwerty123456"},
			mockBehavior: func(s *mock_service.MockAuth, l *mock_logger.MockLogger, signUp core.SignUp, signIn core.SignIn) {
				s.EXPECT().SignUp(gomock.Any(), signUp).Return(" ", nil)
				s.EXPECT().SignIn(gomock.Any(), signIn).Return(auth.Tokens{}, errServiceSignIn)
				l.EXPECT().Error(errServiceSignIn)
			},
			expStatCode: 500,
			expReqBody:  fmt.Sprintf(`{"message":"%s"}`, errServiceSignIn),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(authSrv, l, tt.SignUp, tt.SignIn)

			req := httptest.NewRequest(http.MethodPost, "/v1/auth/signup", bytes.NewBufferString(tt.inputBody))
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			require.Equal(t, tt.expStatCode, w.Code)
			require.Equal(t, w.Body.String(), tt.expReqBody)
		})
	}
}

func TestAuthHandler_signIn(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuth, l *mock_logger.MockLogger, input core.SignIn)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authSrv := mock_service.NewMockAuth(ctrl)
	tm := mock_auth.NewMockTokenManager(ctrl)
	l := mock_logger.NewMockLogger(ctrl)

	services := &service.Services{Auth: authSrv}
	deps := initDeps(services, tm, l)
	handler := NewAuthHandler(deps)

	r := gin.New()
	v1 := r.Group("/v1")
	handler.initAuthRoutes(v1)

	tests := []struct {
		name         string
		mockBehavior mockBehavior
		inputBody    string
		SignUp       core.SignIn
		expStatCode  int
		expReqBody   string
	}{
		{
			name:      "OK",
			inputBody: `{"username":"Cheasezz","password":"qwerty123456"}`,
			SignUp:    core.SignIn{Username: "Cheasezz", Password: "qwerty123456"},
			mockBehavior: func(s *mock_service.MockAuth, l *mock_logger.MockLogger, input core.SignIn) {
				s.EXPECT().SignIn(gomock.Any(), input).Return(tokens, nil)
			},
			expStatCode: 200,
			expReqBody:  fmt.Sprintf(`{"accessToken":"%s"}`, tokens.Access),
		},
		{
			name:      "Bad request: empty username",
			inputBody: `{"username":"","password":"qwerty123456"}`,
			SignUp:    core.SignIn{Username: "", Password: "qwerty123456"},
			mockBehavior: func(s *mock_service.MockAuth, l *mock_logger.MockLogger, input core.SignIn) {
				l.EXPECT().Error(gomock.Any())
			},
			expStatCode: 400,
			expReqBody:  `{"message":"Key: 'SignIn.Username' Error:Field validation for 'Username' failed on the 'required' tag"}`,
		},
		{
			name:      "Bad request: empty password",
			inputBody: `{"username":"Cheasezz","password":""}`,
			SignUp:    core.SignIn{Username: "Cheasezz", Password: ""},
			mockBehavior: func(s *mock_service.MockAuth, l *mock_logger.MockLogger, input core.SignIn) {
				l.EXPECT().Error(gomock.Any())
			},
			expStatCode: 400,
			expReqBody:  `{"message":"Key: 'SignIn.Password' Error:Field validation for 'Password' failed on the 'required' tag"}`,
		},
		{
			name:      "Server error: Service Sign In error",
			inputBody: `{"username":"Cheasezz","password":"qwerty123456"}`,
			SignUp:    core.SignIn{Username: "Cheasezz", Password: "qwerty123456"},
			mockBehavior: func(s *mock_service.MockAuth, l *mock_logger.MockLogger, input core.SignIn) {
				s.EXPECT().SignIn(gomock.Any(), input).Return(auth.Tokens{}, errServiceSignIn)
				l.EXPECT().Error(errServiceSignIn)
			},
			expStatCode: 500,
			expReqBody:  fmt.Sprintf(`{"message":"%s"}`, errServiceSignIn),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(authSrv, l, tt.SignUp)

			req := httptest.NewRequest(http.MethodPost, "/v1/auth/signin", bytes.NewBufferString(tt.inputBody))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			require.Equal(t, tt.expStatCode, w.Code)
			require.Equal(t, w.Body.String(), tt.expReqBody)
		})
	}
}

func TestAuth_logOut(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuth, l *mock_logger.MockLogger, refreshToken string, expTkn auth.Tokens)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authSrv := mock_service.NewMockAuth(ctrl)
	tm := mock_auth.NewMockTokenManager(ctrl)
	l := mock_logger.NewMockLogger(ctrl)

	services := &service.Services{Auth: authSrv}
	deps := initDeps(services, tm, l)
	handler := NewAuthHandler(deps)

	r := gin.New()
	v1 := r.Group("/v1")
	handler.initAuthRoutes(v1)

	tests := []struct {
		name         string
		cookieName   string
		rToken       string
		expStatCode  int
		expReqBody   string
		expTokens    auth.Tokens
		mockBehavior mockBehavior
	}{
		{
			name:        "ok",
			cookieName:  rtCookieName,
			rToken:      "token",
			expStatCode: 200,
			expReqBody:  `{"accessToken":""}`,
			expTokens:   auth.Tokens{Access: "", Refresh: auth.RTInfo{Token: "", ExpiresAt: time.Now(), TTLInSec: 0}},
			mockBehavior: func(s *mock_service.MockAuth, l *mock_logger.MockLogger, rt string, expTkn auth.Tokens) {
				s.EXPECT().LogOut(gomock.Any(), rt).Return(expTkn, nil)
			},
		},
		{
			name:        "empty cookie name",
			expStatCode: 401,
			expReqBody:  fmt.Sprintf(`{"message":"%s"}`, http.ErrNoCookie),
			mockBehavior: func(s *mock_service.MockAuth, l *mock_logger.MockLogger, rt string, expTkn auth.Tokens) {
				l.EXPECT().Error(http.ErrNoCookie)
			},
		},
		{
			name:        "service error",
			cookieName:  rtCookieName,
			rToken:      "token",
			expStatCode: 500,
			expReqBody:  fmt.Sprintf(`{"message":"%s"}`, errServiceLogOut),
			mockBehavior: func(s *mock_service.MockAuth, l *mock_logger.MockLogger, rt string, expTkn auth.Tokens) {
				s.EXPECT().LogOut(gomock.Any(), rt).Return(expTkn, errServiceLogOut)
				l.EXPECT().Error(errServiceLogOut)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(authSrv, l, tt.rToken, tt.expTokens)

			req := httptest.NewRequest(http.MethodGet, "/v1/auth/logout", nil)
			req.AddCookie(&http.Cookie{
				Name:     tt.cookieName,
				Value:    tt.rToken,
				MaxAge:   tokens.Refresh.TTLInSec,
				Expires:  tokens.Refresh.ExpiresAt,
				Path:     "/",
				Domain:   deps.ConfigHTTP.Host,
				HttpOnly: true,
				SameSite: http.SameSiteNoneMode,
				Secure:   true,
			})
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			require.Equal(t, tt.expStatCode, w.Code)
			require.Equal(t, w.Body.String(), tt.expReqBody)
			require.Empty(t, w.Result().Header.Get(rtCookieName))
		})
	}
}

func TestAuth_refreshAccessToken(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuth, l *mock_logger.MockLogger, refreshToken string)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authSrv := mock_service.NewMockAuth(ctrl)
	tm := mock_auth.NewMockTokenManager(ctrl)
	l := mock_logger.NewMockLogger(ctrl)

	services := &service.Services{Auth: authSrv}
	deps := initDeps(services, tm, l)
	handler := NewAuthHandler(deps)

	r := gin.New()
	v1 := r.Group("/v1")
	handler.initAuthRoutes(v1)

	tests := []struct {
		name         string
		mockBehavior mockBehavior
		// tokens       auth.Tokens
		cookieName   string
		refreshToken string
		expStatCode  int
		expReqBody   string
	}{
		{
			name: "ok (upd both tokens)",
			mockBehavior: func(s *mock_service.MockAuth, l *mock_logger.MockLogger, refreshToken string) {
				s.EXPECT().RefreshAccessToken(gomock.Any(), refreshToken).Return(tokens, nil)
			},
			// tokens: initTokens(),
			cookieName:   "RefreshToken",
			refreshToken: "token",
			expStatCode:  200,
			expReqBody:   fmt.Sprintf(`{"accessToken":"%s"}`, tokens.Access),
		},
		{
			name: "ok (upd access token only)",
			mockBehavior: func(s *mock_service.MockAuth, l *mock_logger.MockLogger, refreshToken string) {
				s.EXPECT().RefreshAccessToken(gomock.Any(), refreshToken).Return(auth.Tokens{Access: "token"}, nil)
			},
			cookieName:   "RefreshToken",
			refreshToken: "token",
			expStatCode:  200,
			expReqBody:   fmt.Sprintf(`{"accessToken":"%s"}`, tokens.Access),
		},
		{
			name: "StatusUnauthorized: empty cookie name",
			mockBehavior: func(s *mock_service.MockAuth, l *mock_logger.MockLogger, refreshToken string) {
				l.EXPECT().Error(http.ErrNoCookie)
			},
			expStatCode: 401,
			expReqBody:  fmt.Sprintf(`{"message":"%s"}`, http.ErrNoCookie),
		},
		{
			name: "Server error: Service Refresh Access Token error",
			mockBehavior: func(s *mock_service.MockAuth, l *mock_logger.MockLogger, refreshToken string) {
				s.EXPECT().RefreshAccessToken(gomock.Any(), refreshToken).Return(auth.Tokens{}, errServiceRefreshAccessToken).Times(1)
				l.EXPECT().Error(errServiceRefreshAccessToken)
			},
			cookieName:   "RefreshToken",
			refreshToken: "token",
			expStatCode:  500,
			expReqBody:   fmt.Sprintf(`{"message":"%s"}`, errServiceRefreshAccessToken),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(authSrv, l, tt.refreshToken)

			req := httptest.NewRequest(http.MethodPost, "/v1/auth/refresh", nil)
			req.AddCookie(&http.Cookie{
				Name:     tt.cookieName,
				Value:    tt.refreshToken,
				MaxAge:   tokens.Refresh.TTLInSec,
				Expires:  tokens.Refresh.ExpiresAt,
				Path:     "/",
				Domain:   deps.ConfigHTTP.Host,
				HttpOnly: true,
				SameSite: http.SameSiteNoneMode,
				Secure:   true,
			})
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			require.Equal(t, tt.expStatCode, w.Code)
			require.Equal(t, w.Body.String(), tt.expReqBody)

		})
	}
}

func TestAuth_me(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuth, l *mock_logger.MockLogger, refreshToken string)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authSrv := mock_service.NewMockAuth(ctrl)
	tm := mock_auth.NewMockTokenManager(ctrl)
	l := mock_logger.NewMockLogger(ctrl)

	services := &service.Services{Auth: authSrv}
	deps := initDeps(services, tm, l)
	handler := NewAuthHandler(deps)

	r := gin.New()
	v1 := r.Group("/v1")
	handler.initAuthRoutes(v1)

	tests := []struct {
		name         string
		accessToken  string
		expStatCode  int
		expReqBody   string
		mockBehavior mockBehavior
	}{
		{
			name:        "ok",
			accessToken: "acToken",
			expStatCode: 200,
			expReqBody:  fmt.Sprint(`{"user":"userName"}`),
			mockBehavior: func(s *mock_service.MockAuth, l *mock_logger.MockLogger, accessToken string) {
				tm.EXPECT().Parse(accessToken).Return("userId1", nil).Times(1)
				s.EXPECT().GetUser(gomock.Any(), "userId1").Return("userName", nil).Times(1)
			},
		},
		{
			name:        "error service GetUser",
			accessToken: "acToken",
			expStatCode: 401,
			expReqBody:  fmt.Sprintf(`{"message":"%s"}`, errServiceGetUser),
			mockBehavior: func(s *mock_service.MockAuth, l *mock_logger.MockLogger, accessToken string) {
				tm.EXPECT().Parse(accessToken).Return("userId1", nil).Times(1)
				s.EXPECT().GetUser(gomock.Any(), "userId1").Return("", errServiceGetUser).Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(authSrv, l, tt.accessToken)

			req := httptest.NewRequest(http.MethodGet, "/v1/auth/me", nil)
			req.Header.Add(authorizationHeader, fmt.Sprintf("Bearer %s", tt.accessToken))
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			require.Equal(t, tt.expStatCode, w.Code)
			require.Equal(t, w.Body.String(), tt.expReqBody)
		})
	}
}
