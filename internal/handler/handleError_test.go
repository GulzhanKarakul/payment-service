package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/GulzhanKarakul/payment-service/internal/domain"
	"github.com/GulzhanKarakul/payment-service/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandleError_MapsCorrectHTTPCodes(t *testing.T) {
	tests := []struct {
		name string
		wantCode int
		setup func(t *testing.T) (*http.Request, http.Handler)
	} {
		// not found 404
		{
			name: "client not found",
			wantCode: http.StatusNotFound,
			setup: func(t *testing.T) (*http.Request, http.Handler) {
				svc := mocks.NewMockClientService(t)
				svc.EXPECT().GetByID(mock.Anything, testClientID).
					Return(domain.Client{}, domain.ErrClientNotFound).Once()
				req := httptest.NewRequest(http.MethodGet, "/api/v1/clients/"+testClientID, nil)
				return req, newHandler(svc, nil, nil, nil)
			},
		},
		{
			name: "business not found",
			wantCode: http.StatusNotFound,
			setup: func(t *testing.T) (*http.Request, http.Handler) {
				svc := mocks.NewMockBusinessService(t)
				svc.EXPECT().GetByID(mock.Anything, testBusinessID).
					Return(domain.Business{}, domain.ErrBusinessNotFound).Once()
				req := httptest.NewRequest(http.MethodGet, "/api/v1/businesses/"+testBusinessID, nil)
				return req, newHandler(nil, svc, nil, nil)
			},
		},
		{
			name: "transaction not found",
			wantCode: http.StatusNotFound,
			setup: func(t *testing.T) (*http.Request, http.Handler) {
				svc := mocks.NewMockTransactionService(t)
				svc.EXPECT().GetByID(mock.Anything, testTransactionID).
					Return(domain.Transaction{}, domain.ErrTransactionNotFound).Once()
				req := httptest.NewRequest(http.MethodGet, "/api/v1/transactions/"+testTransactionID, nil)
				return req, newHandler(nil, nil, svc, nil)
			},
		},
		{
			name: "bonus settings not found",
			wantCode: http.StatusNotFound,
			setup: func(t *testing.T) (*http.Request, http.Handler) {
				svc := mocks.NewMockBonusSettingsService(t)
				svc.EXPECT().GetByBusinessID(mock.Anything, testBusinessID).
					Return(domain.BonusSettings{}, domain.ErrBonusSettingsNotFound).Once()
				req := httptest.NewRequest(http.MethodGet, "/api/v1/businesses/"+testBusinessID+"/settings", nil)
				return req, newHandler(nil, nil, nil, svc)
			},
		},
		// 409 conflict
		{
			name: "client already exist",
			wantCode: http.StatusConflict,
			setup: func(t *testing.T) (*http.Request, http.Handler) {
				svc := mocks.NewMockClientService(t)
				svc.EXPECT().Create(mock.Anything, testPhone, testName).
					Return(domain.Client{}, domain.ErrClientAlreadyExist).Once()
				body := `{"phone":"` + testPhone + `","name":"` +testName + `"}`
				req := httptest.NewRequest(http.MethodPost, "/api/v1/clients", strings.NewReader(body))
				return req, newHandler(svc, nil, nil, nil)
			},
		},
		{
			name: "transaction cancelled",
			wantCode: http.StatusConflict,
			setup: func(t *testing.T) (*http.Request, http.Handler) {
				svc := mocks.NewMockTransactionService(t)
				svc.EXPECT().Cancel(mock.Anything, testTransactionID).
					Return(domain.ErrTransactionCancelled).Once()
				req := httptest.NewRequest(http.MethodPatch, "/api/v1/transactions/"+testTransactionID+"/cancel", nil)
				return req, newHandler(nil, nil, svc, nil)
			},
		},
		// 422 unprocessable entity
		{
			name: "client is not active",
			wantCode: http.StatusUnprocessableEntity,
			setup: func(t *testing.T) (*http.Request, http.Handler) {
				svc := mocks.NewMockTransactionService(t)
				svc.EXPECT().Create(mock.Anything, testClientID, testBusinessID, testAmount, mock.Anything).
					Return(domain.Transaction{}, domain.ErrClientIsNotActive).Once()
				body := `{"client_id":"`+testClientID+`","business_id":"`+testBusinessID+`","amount":500000}`
				req := httptest.NewRequest(http.MethodPost, "/api/v1/transactions", strings.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				return req, newHandler(nil, nil, svc, nil)
			},
		},
		{
			name: "business is not active",
			wantCode: http.StatusUnprocessableEntity,
			setup: func(t *testing.T) (*http.Request, http.Handler) {
				svc := mocks.NewMockTransactionService(t)
				svc.EXPECT().Create(mock.Anything, testClientID, testBusinessID, testAmount, mock.Anything).
					Return(domain.Transaction{}, domain.ErrBusinessIsNotActive).Once()
				body := `{"client_id":"`+testClientID+`","business_id":"`+testBusinessID+`","amount":500000}`
				req := httptest.NewRequest(http.MethodPost, "/api/v1/transactions", strings.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				return req, newHandler(nil, nil, svc, nil)
			},
		},
		{
			name: "insufficient balance",
			wantCode: http.StatusUnprocessableEntity,
			setup: func(t *testing.T) (*http.Request, http.Handler) {
				svc := mocks.NewMockTransactionService(t)
				svc.EXPECT().Create(mock.Anything, testClientID, testBusinessID, testAmount, mock.Anything).
					Return(domain.Transaction{}, domain.ErrInsufficientBonusBalance).Once()
				body := `{"client_id":"`+testClientID+`","business_id":"`+testBusinessID+`","amount":500000}`
				req := httptest.NewRequest(http.MethodPost, "/api/v1/transactions", strings.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				return req, newHandler(nil, nil, svc, nil)
			},
		},
		// server error 500
		{
			name: "unknown error",
			wantCode: http.StatusInternalServerError,
			setup: func(t *testing.T) (*http.Request, http.Handler) {
				svc := mocks.NewMockClientService(t)
				svc.EXPECT().GetByID(mock.Anything, testClientID).
					Return(domain.Client{}, errors.New("db error")).Once()
				req := httptest.NewRequest(http.MethodGet, "/api/v1/clients/"+testClientID, nil)
				return req, newHandler(svc, nil, nil, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, router := tt.setup(t)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantCode, rec.Code)
			var resp errorResponse
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
			assert.NotEmpty(t, resp.Error)
		})
	}
}