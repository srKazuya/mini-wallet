package wallet

import "context"

type Storage interface {
	AddTransaction(ctx context.Context, Transaction Transaction) (Transaction, error)
	GetWallet(ctx context.Context, id int) (Wallet, error)
}
