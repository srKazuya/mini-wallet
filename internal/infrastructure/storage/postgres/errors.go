package postgres

import (
	"errors"
)

var (
	ErrOpenDB          = errors.New("failed to open database")
	ErrMigration       = errors.New("failed to run migrations")
	ErrGormOpen        = errors.New("failed to gorm open")
	ErrReviewerNotInPR = errors.New("reviewer not in pr")
)