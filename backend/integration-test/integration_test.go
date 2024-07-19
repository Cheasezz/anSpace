package integration_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/Cheasezz/anSpace/backend/config"
	"github.com/Cheasezz/anSpace/backend/internal/app"
	repositories "github.com/Cheasezz/anSpace/backend/internal/repository"
	"github.com/Cheasezz/anSpace/backend/internal/service"
	httpHandlers "github.com/Cheasezz/anSpace/backend/internal/transport/http"
	v1 "github.com/Cheasezz/anSpace/backend/internal/transport/http/v1"
	"github.com/Cheasezz/anSpace/backend/pkg/auth"
	hasher "github.com/Cheasezz/anSpace/backend/pkg/hasher"
	"github.com/Cheasezz/anSpace/backend/pkg/logger"
	"github.com/Cheasezz/anSpace/backend/pkg/postgres"
	httpserver "github.com/Cheasezz/anSpace/backend/pkg/server"
	"github.com/stretchr/testify/suite"
)

type APITestSuite struct {
	suite.Suite

	db       *postgres.Postgres
	repos    *repositories.Repositories
	services *service.Services
	handlers *httpHandlers.Handlers
	server   *httpserver.Server

	hasher       hasher.PasswordHasher
	tokenManager auth.TokenManager
	logger       logger.Logger

	userCookie string
}

func (s *APITestSuite) SetupSuite() {
	l := logger.New("info")

	cfg, err := config.NewConfigIntTest()
	if err != nil {
		l.Fatal("failed initialize config integration: %s", err.Error())
	}

	psql, err := postgres.NewPostgressDB(cfg.PG)
	if err != nil {
		l.Fatal("failed initialize db: %s", err.Error())
	}

	app.DBMigrate(cfg.PG, l)

	hasher := hasher.NewSHA1Hasher(cfg.Hasher)
	tokenManager, err := auth.NewManager(cfg.TokenManager)
	if err != nil {
		l.Fatal("failed initialize tokenManager: %s", err.Error())
	}

	repos := repositories.NewRepositories(psql)

	services := service.NewServices(service.Deps{
		Repos:        repos,
		Hasher:       hasher,
		TokenManager: tokenManager,
	})

	handlers := httpHandlers.NewHandlers(v1.Deps{
		Services:     services,
		TokenManager: tokenManager,
		ConfigHTTP:   cfg.HTTP,
		Log:          l,
	})

	server := httpserver.NewServer(cfg.HTTP, handlers.Init())
	s.logger = l
	s.hasher = hasher
	s.tokenManager = tokenManager
	s.db = psql
	s.repos = repos
	s.services = services
	s.handlers = handlers
	s.server = server

}
func (s *APITestSuite) SetupTest() {
	_, err := s.db.Pool.Exec(context.Background(), "truncate users, users_sessions")
	if err != nil {
		s.logger.Error("db exec error: %s", err.Error())
	}
	var username string
	var st []byte

	inputSignUp := `{"Name": "Iurii", "Username": "Cheasezz", "Password": "qwerty123456"}`
	r := s.Require()

	// SignUp for create new user before test every handlers
	resp, err := http.Post("http://"+s.server.HttpServer.Addr+"/api/v1/auth/signup", "json", bytes.NewBufferString(inputSignUp))
	if err != nil {
		s.logger.Error("http post signup error: %s", err.Error())
	}
	err = s.db.Scany.Get(context.Background(), s.db.Pool, &username, `select username from users where username='Cheasezz' and name='Iurii'`)
	if err != nil {
		s.logger.Error("FromTestSignUp db scany get error: %s", err.Error())
	}
	st, _ = io.ReadAll(resp.Body)
	s.userCookie = resp.Cookies()[0].Value

	r.Equal(http.StatusOK, resp.StatusCode)
	r.Equal("Cheasezz", username)
	r.Contains(string(st), `{"accessToken":`)
	r.Equal(resp.Cookies()[0].Name, "RefreshToken")

}
func (s *APITestSuite) TearDownTest() {

}
func (s *APITestSuite) TearDownSuite() {
	if err := s.server.Shutdown(); err != nil {
		s.logger.Error("error occured on server shutting down: %s", err.Error())
	}
	s.db.Close()
}

func TestAPISuite(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	suite.Run(t, new(APITestSuite))
}
