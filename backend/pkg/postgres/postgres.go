package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Cheasezz/anSpace/backend/config"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	pgxUUID "github.com/vgarvardt/pgx-google-uuid/v5"
)

const (
	connAttempts = 10
	connTimeout  = time.Second
)

type Postgres struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration

	Pool  *pgxpool.Pool
	Scany *pgxscan.API
}

func NewPostgressDB(cfg config.PG) (*Postgres, error) {
	pg := &Postgres{
		maxPoolSize:  cfg.PoolMax,
		connAttempts: connAttempts,
		connTimeout:  connTimeout,
		Scany:        pgxscan.DefaultAPI,
	}

	poolConfig, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - pgxpool.ParseConfig: %w", err)
	}

	poolConfig.MaxConns = int32(pg.maxPoolSize)

	poolConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxUUID.Register(conn.TypeMap())

		return nil
	}

	for pg.connAttempts > 0 {
		pg.Pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err == nil {
			break
		}

		log.Printf("Postgres is trying to connect, attempts left: %d", pg.connAttempts)

		time.Sleep(pg.connTimeout)

		pg.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - connAttempts == 0: %w", err)
	}

	if err = pg.Pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("could not ping postgres: %w", err)
	}

	log.Printf("Postgres connected, connAttempts: %d", pg.connAttempts)

	return pg, nil
}

func (p *Postgres) Close() {
	if p.Pool != nil {
		log.Print("Postgres closed")
		p.Pool.Close()
	}
}
