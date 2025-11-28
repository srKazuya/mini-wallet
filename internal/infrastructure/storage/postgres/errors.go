package postgres

import (
	"errors"
)

var (
	ErrOpenDB         = errors.New("failed to open database")
	ErrMigration      = errors.New("failed to run migrations")
	ErrInsFunds       = errors.New("Insufficient funds")
	ErrWalletNotFound = errors.New("Wallet not Found")
)
