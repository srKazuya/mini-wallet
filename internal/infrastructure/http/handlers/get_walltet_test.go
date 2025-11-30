package handlers_test

import (
	"encoding/json"
	"mini-wallet/internal/domain/wallet"
	"net/http"
	"net/http/httptest"
	"strconv"

	"mini-wallet/internal/domain/wallet/mocks"
	"mini-wallet/internal/infrastructure/http/handlers"
	dto "mini-wallet/internal/infrastructure/http/handlers/dto"
	"mini-wallet/internal/infrastructure/storage/postgres"
	slogdiscard "mini-wallet/pkg/sl_logger/slog_discard"
	validateresp "mini-wallet/pkg/validator"

	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetWalletHandler(t *testing.T) {
	fixedTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)

	cases := []struct {
		name           string
		mockReturnOK   wallet.Wallet
		mockReturnErr  error
		expectedStatus int
		expectedResp   dto.GetWalletResponse
	}{
		{
			name: "Success",
			mockReturnOK: wallet.Wallet{
				ID:        1,
				Balance:   100,
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockReturnErr:  nil,
			expectedStatus: http.StatusOK,
			expectedResp: dto.GetWalletResponse{
				ValidationResponse: validateresp.OK(),
				ID:                 1,
				Balance:            100,
				CreatedAt:          fixedTime,
				UpdatedAt:          fixedTime,
			},
		},
		{
			name:           "No wallet",
			expectedStatus: http.StatusNotFound,
			mockReturnErr: postgres.ErrWalletNotFound,
			expectedResp: dto.GetWalletResponse{
				ValidationResponse: validateresp.ValidationResponse{
					Status: "Error",
					Errors: map[string]string{
						"error": "Wallet not Found",
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svcMock := mocks.NewService(t)

			svcMock.On("GetWallet", mock.Anything, 1).
				Return(tc.mockReturnOK, tc.mockReturnErr).
				Once()

			handler := handlers.NewGetWallet(slogdiscard.NewDiscardLogger(), svcMock)

			id := 1
			r := chi.NewRouter()
			r.Get("/api/v1/wallets/{WALLET_UUID}", handler.ServeHTTP)
			req, err := http.NewRequest(http.MethodGet, "/api/v1/wallets/"+strconv.Itoa(id), nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			require.Equal(t, tc.expectedStatus, rr.Code)

			var resp dto.GetWalletResponse
			require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))

			require.Equal(t, tc.expectedResp.Status, resp.Status)
			require.Equal(t, tc.expectedResp.Errors, resp.Errors)
			require.Equal(t, tc.expectedResp.Balance, resp.Balance)
			require.Equal(t, tc.expectedResp.ID, resp.ID)

			if tc.expectedStatus == http.StatusOK {
				require.WithinDuration(t, tc.expectedResp.CreatedAt, resp.CreatedAt, time.Second)
				require.WithinDuration(t, tc.expectedResp.UpdatedAt, resp.UpdatedAt, time.Second)
			}

			svcMock.AssertExpectations(t)
		})
	}
}
