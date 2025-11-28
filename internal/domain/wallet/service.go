package wallet

import (
	"log/slog"
)

type Service interface {
}
type service struct {
	log     *slog.Logger
	storage Storage
}

func NewService(log *slog.Logger, storage Storage) Service {
	return &service{log: log, storage: storage}
}

