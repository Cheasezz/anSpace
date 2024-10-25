package v1

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Cheasezz/anSpace/backend/internal/service"
	mock_service "github.com/Cheasezz/anSpace/backend/internal/service/mocks"
	mock_auth "github.com/Cheasezz/anSpace/backend/pkg/auth/mocks"
	mock_logger "github.com/Cheasezz/anSpace/backend/pkg/logger/mocks"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var (
	errTmParse = fmt.Errorf("token manager parse error")
)

func TestAuth_userIdentity(t *testing.T) {
	type mockBehavior func(s *mock_auth.MockTokenManager, l *mock_logger.MockLogger, token string)

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

	r.GET("/protected", handler.userIdentity, func(ctx *gin.Context) {
		id, _ := ctx.Get(userCtx)
		ctx.String(200, id.(string))
	})

	tests := []struct {
		name         string
		headerName   string
		headerValue  string
		token        string
		mockBehavior mockBehavior
		expStatCode  int
		expReqBody   string
	}{
		{
			name:        "OK",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(s *mock_auth.MockTokenManager, l *mock_logger.MockLogger, token string) {
				s.EXPECT().Parse(token).Return("1", nil)
			},
			expStatCode: 200,
			expReqBody:  "1",
		},
		{
			name:       "Empty auth header",
			headerName: "",
			mockBehavior: func(s *mock_auth.MockTokenManager, l *mock_logger.MockLogger, token string) {
				l.EXPECT().Error(errEmptyAuthHeader)
			},
			expStatCode: 401,
			expReqBody:  fmt.Sprintf(`{"message":"%s"}`, errEmptyAuthHeader),
		},
		{
			name:        "One word in auth header value",
			headerName:  "Authorization",
			headerValue: "token",
			token:       "token",
			mockBehavior: func(s *mock_auth.MockTokenManager, l *mock_logger.MockLogger, token string) {
				l.EXPECT().Error(errInvalidAuthHeader)
			},
			expStatCode: 401,
			expReqBody:  fmt.Sprintf(`{"message":"%s"}`, errInvalidAuthHeader),
		},
		{
			name:        "Uncorrect Bearer key word in auth header",
			headerName:  "Authorization",
			headerValue: "Bearrer token",
			token:       "token",
			mockBehavior: func(s *mock_auth.MockTokenManager, l *mock_logger.MockLogger, token string) {
				l.EXPECT().Error(errInvalidAuthHeader)
			},
			expStatCode: 401,
			expReqBody:  fmt.Sprintf(`{"message":"%s"}`, errInvalidAuthHeader),
		},
		{
			name:        "Empty token in auth header",
			headerName:  "Authorization",
			headerValue: "Bearer",
			token:       "token",
			mockBehavior: func(s *mock_auth.MockTokenManager, l *mock_logger.MockLogger, token string) {
				l.EXPECT().Error(errInvalidAuthHeader)
			},
			expStatCode: 401,
			expReqBody:  fmt.Sprintf(`{"message":"%s"}`, errInvalidAuthHeader),
		},
		{
			name:        "Token manager parse error",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(s *mock_auth.MockTokenManager, l *mock_logger.MockLogger, token string) {
				s.EXPECT().Parse(token).Return("", errTmParse)
				l.EXPECT().Error(errTmParse)
			},
			expStatCode: 500,
			expReqBody:  fmt.Sprintf(`{"message":"%s"}`, errTmParse),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tm, l, tt.token)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/protected", nil)
			req.Header.Set(tt.headerName, tt.headerValue)

			r.ServeHTTP(w, req)

			require.Equal(t, w.Code, tt.expStatCode)
			require.Equal(t, w.Body.String(), tt.expReqBody)
		})
	}
}

func Test_getUserIdFrmCtx(t *testing.T) {
	var getContext = func(id uuid.UUID) *gin.Context {
		c := &gin.Context{}
		c.Set(userCtx, id.String())
		return c
	}
	testUUID := uuid.New()
	tests := []struct {
		name string
		c    *gin.Context
		id   uuid.UUID
		fail bool
	}{
		{
			name: "OK",
			c:    getContext(testUUID),
			id:   testUUID,
		},
		{
			name: "Empty",
			c:    &gin.Context{},
			fail: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := getUserIdFrmCtx(tt.c)
			if tt.fail {
				require.Error(t, err)
				require.EqualError(t, errUserIdNotFound, err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, id, tt.id)
			}
		})
	}
}
