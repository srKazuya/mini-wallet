package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	w "mini-wallet/internal/domain/wallet"

	"github.com/lib/pq"
	_ "github.com/lib/pq"

	"github.com/pressly/goose/v3"
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

	log.Info("start seeding...")
	if cfg.Seed {
		if err := SeedWallets(sqlDB, log); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	return &PostgresStorage{db: sqlDB}, nil
}

func (s *PostgresStorage) AddTransaction(ctx context.Context, t w.Transaction) (w.Transaction, error) {
	const op = "stroage.pg.AddTransaction"

	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return w.Transaction{}, fmt.Errorf("%s: begin tx:%w", op, err)
	}
	defer tx.Rollback()

	var currentBalance float64
	if err = tx.QueryRowContext(ctx, `
		SELECT balance 
		FROM wallets
		WHERE id = $1
		FOR UPDATE
	`, t.WalletID).Scan(&currentBalance); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return w.Transaction{}, ErrWalletNotFound
		}
		return w.Transaction{}, fmt.Errorf("%s: select wallet: %w", op, err)
	}

	if t.TrType == "withdraw" && currentBalance < t.Amount {
		return w.Transaction{}, ErrInsFunds
	}

	newBalance := currentBalance
	if t.TrType == "deposit" {
		newBalance += t.Amount
	} else {
		newBalance -= t.Amount
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE wallets
		SET balance = $1
		WHERE id = $2
	`, newBalance, t.WalletID)
	if err != nil {
		return w.Transaction{}, fmt.Errorf("%s: update balance:%w", op, err)
	}

	query := `
		INSERT INTO transactions (wallet_id, amount, transaction_type)
		VALUES ($1, $2, $3)
		RETURNING id, wallet_id, amount, transaction_type, created_at
	`

	var output w.Transaction

	err = tx.QueryRowContext(ctx, query, t.WalletID, t.Amount, t.TrType).Scan(
		&output.ID,
		&output.WalletID,
		&output.Amount,
		&output.TrType,
		&output.CreatedAt,
	)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return w.Transaction{}, ErrWalletNotFound
		}
		return w.Transaction{}, fmt.Errorf("%s: insert transaction: %w", op, err)
	}

	if err := tx.Commit(); err != nil {
		return w.Transaction{}, fmt.Errorf("%s: commit: %w", op, err)
	}

	return output, nil
}
