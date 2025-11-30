package handlers

import (
	"errors"
	"log/slog"
	"mini-wallet/internal/domain/wallet"
	dto "mini-wallet/internal/infrastructure/http/handlers/dto"
	"mini-wallet/internal/infrastructure/http/transport"
	"mini-wallet/internal/infrastructure/storage/postgres"
	"mini-wallet/pkg/sl_logger/sl"
	valResp "mini-wallet/pkg/validator"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewGetWallet(log *slog.Logger, svc wallet.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.transaction.get"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		idStr := chi.URLParam(r, "WALLET_UUID")
		log.Info("wallet getted", slog.Any("idSr", idStr)) 
		ctx := r.Context()
		
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Error("failed to convert string",
				sl.Err(err),
			)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		wallet, err := svc.GetWallet(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, postgres.ErrWalletNotFound):
				log.Error(postgres.ErrWalletNotFound.Error(), sl.Err(err))
				getWalletResponseErr(w, http.StatusNotFound, postgres.ErrWalletNotFound.Error())
				return
			default:
				log.Error("unexpected error getting event", sl.Err(err))
				return
			}
		}
		log.Info("wallet getted", slog.Any("walletID", wallet.ID), slog.Any("balance", wallet.Balance), slog.Any("createdAt", wallet.CreatedAt))

		getWalletResponseOK(w, wallet)
	}
}

func getWalletResponseOK(w http.ResponseWriter, wal wallet.Wallet) {
	r := dto.GetWalletResponse{
		ValidationResponse: valResp.OK(),
		ID:                 wal.ID,
		Balance:            wal.Balance,
		CreatedAt:          wal.CreatedAt,
		UpdatedAt:          wal.UpdatedAt,
	}
	if err := transport.WriteJSON(w, http.StatusOK, r); err != nil {
		http.Error(w, "failed to encode JSON", http.StatusInternalServerError)
		return
	}
}

func getWalletResponseErr(w http.ResponseWriter, status int, e string) {
	r := dto.GetWalletResponse{
		ValidationResponse: valResp.Error(e),
	}
	if err := transport.WriteJSON(w, status, r); err != nil {
		http.Error(w, "failed to encode JSON", http.StatusInternalServerError)
		return
	}
}
