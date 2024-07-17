package integration_test

import (
	"bytes"
	"context"
	"fmt"
	"hash"
	"io"
	"net/http"
	"testing"
	"time"

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

	hasher       hash.Hash
	tokenManager auth.TokenManager
	logger       logger.Logger

	createUser bool
	userCookie string
}

const pgURL = "postgres://postgres:qwerty@localhost:5432/postgres?sslmode=disable"
const schemaUrl = "../schema"

func (s *APITestSuite) SetupSuite() {
	l := logger.New("info")

	psql, err := postgres.NewPostgressDB(config.PG{PoolMax: 4, URL: pgURL})
	if err != nil {
		l.Fatal("failed initialize db: %s", err.Error())
	}

	app.DBMigrate(schemaUrl, pgURL, l)

	hasher := hasher.NewSHA1Hasher(config.Hasher{Salt: "salt"})
	tokenManager, err := auth.NewManager(config.TokenManager{
		SigningKey:      "siginKey",
		AccessTokenTTL:  time.Hour,
		RefreshTokenTTL: time.Hour * 48,
	})
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
		ConfigHTTP:   config.HTTP{Host: "127.0.0.1", Port: "8080"},
		Log:          l,
	})

	server := httpserver.NewServer(config.HTTP{Host: "127.0.0.1", Port: "8080"}, handlers.Init())
	s.logger = l
	s.db = psql
	s.repos = repos
	s.services = services
	s.handlers = handlers
	s.server = server
	// s.createUser = true

}
func (s *APITestSuite) SetupTest() {
	s.db.Pool.Exec(context.Background(), "truncate users, users_sessions")
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
		r.Contains(fmt.Sprintf("%s", st), `{"accessToken":`)
		r.Equal(resp.Cookies()[0].Name, "RefreshToken")

}
func (s *APITestSuite) TearDownTest() {

}
func (s *APITestSuite) TearDownSuite() {
	s.server.Shutdown()
	s.db.Close()
}

func (s *APITestSuite) initDeps(psql *postgres.Postgres) {

}

func TestAPISuite(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	suite.Run(t, new(APITestSuite))
}
