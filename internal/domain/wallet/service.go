package wallet

import (
	"context"
	"log/slog"
)

type Service interface {
	AddTransaction(ctx context.Context, t Transaction) (Transaction, error)
	GetWallet(ctx context.Context, id int) (Wallet, error)
}
type service struct {
	log     *slog.Logger
	storage Storage
}

func NewService(log *slog.Logger, storage Storage) Service {
	return &service{log: log, storage: storage}
}

func (s *service) AddTransaction(ctx context.Context, t Transaction) (Transaction, error) {
	return s.storage.AddTransaction(ctx, t)
}
func (s *service) GetWallet(ctx context.Context, id int) (Wallet, error) {
	return s.storage.GetWallet(ctx, id)
}
