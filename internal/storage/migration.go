package storage

import (
	"context"
	"io/fs"

	"github.com/fungicibus/order/internal/logger"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func (s *service) MigrationUp(ctx context.Context, fs fs.FS, log *logger.Logger) error {
	goose.SetLogger(logger.WtihSource(log, "goose"))
	goose.SetBaseFS(fs)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	conn := stdlib.OpenDBFromPool(s.rwPool)

	return goose.Up(conn, "migrations")
}
