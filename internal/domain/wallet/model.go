package wallet

import "time"

type Transaction struct {
	ID      uint64
	WalletID  uint64
	TrType    string
	Amount    float64
	CreatedAt time.Time
}
