package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/fungicibus/order/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type service struct {
	rwPool *pgxpool.Pool
	roPool *pgxpool.Pool

	pingTimeout time.Duration
}

func New(ctx context.Context, cfg config.Postgres) (*service, error) {
	rwConfig, err := pgxpool.ParseConfig(cfg.RWDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse rw config: %w", err)
	}
	rwPool, err := pgxpool.NewWithConfig(ctx, rwConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create new rw pool: %w", err)
	}

	roConfig, err := pgxpool.ParseConfig(cfg.RODSN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ro config: %w", err)
	}
	roPool, err := pgxpool.NewWithConfig(ctx, roConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create new ro pool: %w", err)
	}

	return &service{
		rwPool:      rwPool,
		roPool:      roPool,
		pingTimeout: cfg.PingTimeout,
	}, nil
}

func (s *service) Ping(ctx context.Context) error {
	pingCtx, cancel := context.WithTimeout(ctx, s.pingTimeout)
	defer cancel()

	if err := s.rwPool.Ping(pingCtx); err != nil {
		return fmt.Errorf("rw pool ping failed: %w", err)
	}
	if err := s.roPool.Ping(pingCtx); err != nil {
		return fmt.Errorf("ro pool ping failed: %w", err)
	}
	return nil
}

func (s *service) Close() {
	s.rwPool.Close()
	s.roPool.Close()
}
