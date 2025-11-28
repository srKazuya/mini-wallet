package wallet

type Transaction struct {
	UUID     uint64
	WalletID uint64
	Type     string
	amount   float64
}
