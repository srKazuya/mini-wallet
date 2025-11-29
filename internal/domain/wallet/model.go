package wallet

import "time"

type Transaction struct {
	ID        uint64
	WalletID  uint64
	TrType    string
	Amount    float64
	CreatedAt time.Time
}

type Wallet struct {
	ID        uint64
	Balance   float64
	CreatedAt time.Time
	UpdatedAt time.Time
}
