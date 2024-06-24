package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Postgres struct {
	connAttempts int
	connTimeout  time.Duration
	maxPoolSz    int
	Builder      squirrel.StatementBuilderType
	PgxPool      *pgxpool.Pool
	poolCfg      *pgxpool.Config
	logger       *slog.Logger
	Url          string
}

type ConnConfig struct {
	host     string
	port     string
	user     string
	password string
	dbname   string
}

func NewConnectionConfig(
	host,
	port,
	user,
	password,
	name string) *ConnConfig {
	return &ConnConfig{
		host:     host,
		port:     port,
		user:     user,
		password: password,
		dbname:   name,
	}
}

func (c *ConnConfig) ParseUrl() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		c.user,
		c.password,
		c.host,
		c.port,
		c.dbname,
	)
}

const (
	defaultMaxPoolSz    = 1
	defaultConnAttempts = 10
	defaultConnTimeout  = 2 * time.Second
)

func New(c *ConnConfig, logger *slog.Logger, opts ...Option) (*Postgres, error) {
	pg := &Postgres{
		connAttempts: defaultConnAttempts,
		connTimeout:  defaultConnTimeout,
		maxPoolSz:    defaultMaxPoolSz,
		logger:       logger,
	}

	pg.Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	for _, opt := range opts {
		opt(pg)
	}

	url := c.ParseUrl()
	pg.Url = url

	var err error
	pg.poolCfg, err = pgxpool.ParseConfig(url)

	if err != nil {
		return nil, fmt.Errorf("error with parse url for postgres pool config: %w", err)
	}

	if err = pg.connectWithAttempts(); err != nil {
		return nil, err
	}

	return pg, nil
}

func (pg *Postgres) connectWithAttempts() error {
	var err error
	for pg.connAttempts > 0 {
		pg.PgxPool, err = pgxpool.ConnectConfig(context.Background(), pg.poolCfg)

		if err == nil {
			return nil
		}

		pg.logger.Warn("Postgres connection failure,", slog.Int("attempts", pg.connAttempts))
		pg.connAttempts--
		time.Sleep(pg.connTimeout)
	}

	return fmt.Errorf("failed connection to postgres, zero attempts left: %w", err)
}
