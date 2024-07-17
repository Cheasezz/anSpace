package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Cheasezz/anSpace/backend/config"
	repositories "github.com/Cheasezz/anSpace/backend/internal/repository"
	"github.com/Cheasezz/anSpace/backend/internal/service"
	httpHandlers "github.com/Cheasezz/anSpace/backend/internal/transport/http"
	v1 "github.com/Cheasezz/anSpace/backend/internal/transport/http/v1"
	"github.com/Cheasezz/anSpace/backend/pkg/auth"
	"github.com/Cheasezz/anSpace/backend/pkg/hasher"
	"github.com/Cheasezz/anSpace/backend/pkg/logger"
	"github.com/Cheasezz/anSpace/backend/pkg/postgres"
	httpserver "github.com/Cheasezz/anSpace/backend/pkg/server"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	psql, err := postgres.NewPostgressDB(cfg.PG)
	if err != nil {
		l.Fatal("failed initialize db: %s", err.Error())
	}
	defer psql.Close()

	DBMigrate(cfg.PG.Schema_Url, cfg.PG.URL, l)

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

	srv := httpserver.NewServer(cfg.HTTP, handlers.Init())
	l.Info("TodoApp Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	select {
	case s := <-quit:
		l.Info("app - Run - signal: " + s.String())
	case err = <-srv.Notify():
		l.Error("app - Run - httpServer.Notify: %s", err)
	}

	if err := srv.Shutdown(); err != nil {
		l.Error("error occured on server shutting down: %s", err.Error())
	}

	l.Info("TodoApp Shutting Down")
}
