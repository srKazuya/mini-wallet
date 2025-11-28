package postgres

import (
	"database/sql"
	"fmt"
	
	_ "github.com/lib/pq"

	"github.com/pressly/goose/v3"
	"log/slog"
)

type PostgresStorage struct {
	db *sql.DB
}

func New(cfg Config, log *slog.Logger) (*PostgresStorage, error) {
	const op = "storage.postgres.NewStrorage"
	log = log.With(
		slog.String("op", op),
	)

	sqlDB, err := sql.Open("postgres", cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %w", op, ErrOpenDB, err)
	}
	log.Info("start migrate...", slog.String("path", cfg.MigrationsPath))
	if err := goose.Up(sqlDB, cfg.MigrationsPath); err != nil {
		return nil, fmt.Errorf("%s: %w: %w", op, ErrMigration, err)
	}

	//DB seed
	// log.Info("start seeding...")

	return &PostgresStorage{db: sqlDB}, nil
}
