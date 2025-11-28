package postgres

import (
	"database/sql"
	"github.com/brianvoe/gofakeit/v7"
	"log/slog"
)

func SeedWallets(db *sql.DB, log *slog.Logger) error {
	gofakeit.Seed(12345)

	for i := 0; i < 10; i++ {

		balance := gofakeit.Float64Range(0, 10000)
		createdAt := gofakeit.Date()
		updatedAt := gofakeit.Date()

		var walletID int
		err := db.QueryRow(`
		INSERT INTO wallets (balance, created_at, updated_at)
		VALUES ($1, $2, $3)
		RETURNING id
		`, balance, createdAt, updatedAt).Scan(&walletID)

		if err != nil {
			return err
		}

		numTx := gofakeit.Number(1, 5)
		for j := 0; j < numTx; j++ {
			txAmount := gofakeit.Float64Range(10, 1000)
			txType := "deposit"
			if gofakeit.Bool() {
				txType = "withdraw"
			}
			txCreated := gofakeit.Date()
			_, err := db.Exec(`
                INSERT INTO transactions (wallet_id, amount, transaction_type, created_at)
                VALUES ($1, $2, $3, $4)
            `, walletID, txAmount, txType, txCreated)
			if err != nil {
				return err
			}
		}
		log.Info("seeding", slog.Int("iteration:", i))
	}

	return nil
}
