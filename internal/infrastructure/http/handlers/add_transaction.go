package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"mini-wallet/internal/domain/wallet"
	dto "mini-wallet/internal/infrastructure/http/handlers/dto"
	"mini-wallet/internal/infrastructure/http/transport"
	"mini-wallet/internal/infrastructure/storage/postgres"
	"mini-wallet/pkg/sl_logger/sl"
	valResp "mini-wallet/pkg/validator"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator"
)

func NewAddTransaction(log *slog.Logger, svc wallet.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.transacton.add"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		ctx := r.Context()
		var req dto.AddTransactionRequest

		err := json.NewDecoder(r.Body).Decode(&req)

		if errors.Is(err, io.EOF) {
			log.Error("bad request",
				slog.String("type", transport.ErrEmptyReqBody.Error()),
				sl.Err(err),
			)
			addTransactionResponseErr(w, http.StatusBadRequest, transport.ErrEmptyReqBody.Error())
			return
		}
		if err != nil {
			log.Error("bad request",
				slog.String("type", transport.ErrFailedToDecodeReqBody.Error()),
				sl.Err(err),
			)
			addTransactionResponseErr(w, http.StatusBadRequest, transport.ErrFailedToDecodeReqBody.Error())
			return
		}

		log.Info("request body decoded", slog.Any("req", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))
			if err = transport.WriteJSON(w, http.StatusBadRequest, valResp.ValidationError(validateErr)); err != nil {
				http.Error(w, "failed to encode JSON", http.StatusInternalServerError)
				return
			}
			return
		}

		t := dto.AddTransactionMapToModel(req)

		transaction, err := svc.AddTransaction(ctx, t)
		if err != nil {
			switch {
			case errors.Is(err, postgres.ErrInsFunds):
				log.Error(postgres.ErrInsFunds.Error(), sl.Err(err))
				addTransactionResponseErr(w, http.StatusBadRequest, postgres.ErrInsFunds.Error())
				return
			case errors.Is(err, postgres.ErrWalletNotFound):
				log.Error(postgres.ErrWalletNotFound.Error(), sl.Err(err))
				addTransactionResponseErr(w, http.StatusBadRequest, postgres.ErrWalletNotFound.Error())
				return
			default:
				log.Error("unexpected error adding event", sl.Err(err))
				return
			}
		}

		log.Info("transaction added", slog.Any("title", transaction.WalletID))

		addTransactionResponseOK(w, transaction)
	}
}

func addTransactionResponseOK(w http.ResponseWriter, t wallet.Transaction) {
	r := dto.AddTransactionResponse{
		ValidationResponse: valResp.OK(),
		TrType:             t.TrType,
		Amount:             t.Amount,
		CreatedAt:          t.CreatedAt,
	}
	if err := transport.WriteJSON(w, http.StatusOK, r); err != nil {
		http.Error(w, "failed to encode JSON", http.StatusInternalServerError)
		return
	}
}

func addTransactionResponseErr(w http.ResponseWriter, status int, e string) {
	r := dto.AddTransactionResponse{
		ValidationResponse: valResp.Error(e),
	}
	if err := transport.WriteJSON(w, status, r); err != nil {
		http.Error(w, "failed to encode JSON", http.StatusInternalServerError)
		return
	}
}
