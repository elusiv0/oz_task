package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type (
	Config struct {
		App      App
		Http     Http
		Postgres Postgres
	}

	App struct {
		Env string `envconfig:"env" required:"true"`
		Db  string `envconfig:"db" required:"true"`
	}

	Http struct {
		Host            string        `envconfig:"HTTP_HOST" default:"localhost"`
		Port            string        `envconfig:"HTTP_PORT" default:"8080"`
		ReadTimeout     time.Duration `envconfig:"HTTP_READTIMEOUT" default:"5s"`
		WriteTimeout    time.Duration `envconfig:"HTTP_WRITETIMEOUT" default:"5s"`
		ShutdownTimeout time.Duration `envconfig:"HTTP_SHUTDOWNTIMEOUT" default:"3s"`
	}

	Postgres struct {
		MaxPoolSz          int           `envconfig:"PG_MAX_POOL_SIZE" default:"1"`
		ConnectionTimeout  time.Duration `envconfig:"PG_CONNECTION_TIMEOUT" default:"3s"`
		ConnectionAttempts int           `envconfig:"PG_CONNECTION_ATTEMPTS" default:"10"`
		Host               string        `envconfig:"PG_HOST" required:"true"`
		Port               string        `envconfig:"PG_PORT" required:"true"`
		User               string        `envconfig:"PG_USER" required:"true"`
		Password           string        `envconfig:"PG_PASSWORD" required:"true"`
		Name               string        `envconfig:"PG_NAME" required:"true"`
	}
)

func NewConfig() (*Config, error) {
	app := App{}
	if err := envconfig.Process("", &app); err != nil {
		return nil, fmt.Errorf("Config - NewConfig: %w", err)
	}
	config := Config{}
	pg := Postgres{}
	if app.Db == "postgres" {
		if err := envconfig.Process("", &pg); err != nil {
			return nil, fmt.Errorf("Config - NewConfig: %w", err)
		}
	}
	httpcfg := Http{}
	if err := envconfig.Process("", &httpcfg); err != nil {
		return nil, fmt.Errorf("Config - NewConfig: %w", err)
	}
	config.App = app
	config.Postgres = pg
	config.Http = httpcfg
	return &config, nil
}
