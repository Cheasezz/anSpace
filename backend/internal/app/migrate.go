package app

import (
	"errors"
	"time"

	"github.com/Cheasezz/anSpace/backend/pkg/logger"
	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/golang-migrate/migrate/source/github"
)

const (
	_defaultAttempts = 20
	_defaultTimeout  = time.Second
)

func DBMigrate(schemaUrl, pgUrl string, l logger.Logger) {
	if len(pgUrl) == 0 {
		l.Fatal("migrate: environment variable not declared: PG_URL")
	}

	var (
		attempts = _defaultAttempts
		err      error
		m        *migrate.Migrate
	)

	for attempts > 0 {
		m, err = migrate.New("file://"+schemaUrl, pgUrl)
		if err == nil {
			break
		}

		l.Info("Migrate: postgres is trying to connect, attempts left: %d", attempts)
		l.Error("err from migrate new error : %s", err)
		time.Sleep(_defaultTimeout)
		attempts--
	}

	if err != nil {
		l.Fatal("Migrate: postgres connect error: %s", err)
	}

	err = m.Up()
	defer m.Close()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		l.Fatal("Migrate: up error: %s", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		l.Info("Migrate: no change")
		return
	}

	l.Info("Migrate: up success")
}
