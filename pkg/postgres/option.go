package postgres

import "time"

type Option func(p *Postgres)

func ConnAttempts(connAttempts int) Option {
	return func(p *Postgres) {
		p.connAttempts = connAttempts
	}
}

func ConnTimeout(connTimeout time.Duration) Option {
	return func(p *Postgres) {
		p.connTimeout = connTimeout
	}
}

func MaxPoolSz(maxPoolSz int) Option {
	return func(p *Postgres) {
		p.maxPoolSz = maxPoolSz
	}
}
