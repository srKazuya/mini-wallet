package hdto

import (
	"mini-wallet/internal/domain/wallet"
	resp "mini-wallet/pkg/validator"
	"time"
)

type AddTransactionRequest struct {
	WalletID uint64  `json:"wallet_id" validate:"required"`
	TrType   string  `json:"type" validate:"required, oneof=deposit withdraw"`
	Amount   float64 `json:"amount" validate:"required,numeric"`
}

type AddTransactionResponse struct {
	TrType    string    `json:"type"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
	resp.ValidationResponse
}

//AddTransactionMapToModel
func AddTransactionMapToModel(t AddTransactionRequest) wallet.Transaction {
	return wallet.Transaction{
		WalletID: t.WalletID,
		TrType:   t.TrType,
		Amount:   t.Amount,
	}
}
