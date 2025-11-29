package handlers_test

import (
	"bytes"
	"encoding/json"
	"mini-wallet/internal/domain/wallet"
	"net/http"
	"net/http/httptest"

	"mini-wallet/internal/domain/wallet/mocks"
	"mini-wallet/internal/infrastructure/http/handlers"
	dto "mini-wallet/internal/infrastructure/http/handlers/dto"
	slogdiscard "mini-wallet/pkg/sl_logger/slog_discard"
	validateresp "mini-wallet/pkg/validator"

	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAddTransactionHandler(t *testing.T) {
	fixedTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)

	cases := []struct {
		name           string
		reqBody        string
		mockReturnOK   *wallet.Transaction
		mockReturnErr  error
		expectedStatus int
		expectedResp   dto.AddTransactionResponse
	}{
		{
			name:    "Success",
			reqBody: `{"wallet_id":1,"type":"deposit","amount":100.0}`,
			mockReturnOK: &wallet.Transaction{
				WalletID:  1,
				TrType:    "deposit",
				Amount:    100,
				CreatedAt: fixedTime,
			},
			mockReturnErr:  nil,
			expectedStatus: http.StatusOK,
			expectedResp: dto.AddTransactionResponse{
				ValidationResponse: validateresp.OK(),
				TrType:             "deposit",
				Amount:             100,
				CreatedAt:          fixedTime,
			},
		}, {
			name:           "Invalid JSON",
			reqBody:        `{"wallet_id":1}`,
			expectedStatus: http.StatusBadRequest,
			expectedResp: dto.AddTransactionResponse{
				ValidationResponse: validateresp.ValidationResponse{
					Status: "Error",
					Errors: map[string]string{
						"Amount": "Это поле обязательно",
						"TrType": "Это поле обязательно",
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svcMock := mocks.NewService(t)

			if tc.name == "Success" || tc.name == "Service error" {
				svcMock.
					On("AddTransaction",
						mock.Anything,
						mock.AnythingOfType("wallet.Transaction"),
					).
					Return(*tc.mockReturnOK, tc.mockReturnErr).
					Once()
			}

			handler := handlers.NewAddTransaction(slogdiscard.NewDiscardLogger(), svcMock)

			req, err := http.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(tc.reqBody)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, tc.expectedStatus, rr.Code)

			var resp dto.AddTransactionResponse
			require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))

			require.Equal(t, tc.expectedResp.Status, resp.Status)
			require.Equal(t, tc.expectedResp.Errors, resp.Errors)
			require.Equal(t, tc.expectedResp.Amount, resp.Amount)

			if tc.expectedStatus == http.StatusOK {
				require.WithinDuration(t, tc.expectedResp.CreatedAt, resp.CreatedAt, time.Second)
			}

			svcMock.AssertExpectations(t)
		})
	}
}
