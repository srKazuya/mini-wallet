package wallet

import "context"

type Storage interface {
	AddTransaction(ctx context.Context, Transaction Transaction) (Transaction, error)
}
