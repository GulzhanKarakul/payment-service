package handler_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/GulzhanKarakul/payment-service/internal/domain"
	"github.com/GulzhanKarakul/payment-service/internal/dto"
	"github.com/GulzhanKarakul/payment-service/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TestCreateTransaction_Success_Returning201 - POST api/v1/transactions/
func TestCreateTransaction_Success_Returning201(t *testing.T) {
	svc := mocks.NewMockTransactionService(t)
	svc.EXPECT().Create(mock.Anything, testClientID, testBusinessID, testAmount, mock.Anything).
		Return(testTransaction(), nil).Once()
	router := newHandler(nil, nil, svc, nil)

	body := fmt.Sprintf(`{"client_id":"%s","business_id":"%s","amount":%d}`, testClientID, testBusinessID, testAmount)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/transactions", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var resp dto.TransactionResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, testClientID, resp.ClientID)
	assert.Equal(t, testBusinessID, resp.BusinessID)
	assert.NotEmpty(t, resp.ID)
}

func TestCreateTransaction_WithDescription_Returning201(t *testing.T) {
	svc := mocks.NewMockTransactionService(t)
	svc.EXPECT().Create(mock.Anything, testClientID, testBusinessID, testAmount, mock.MatchedBy(func(d *string) bool {
		return d != nil && *d == "оплата кофе"
	})).Return(testTransaction(), nil).Once()
	router := newHandler(nil, nil, svc, nil)

	body := fmt.Sprintf(`{"client_id":"%s","business_id":"%s","amount":%d,"description":"оплата кофе"}`, testClientID, testBusinessID, testAmount)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/transactions", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var resp dto.TransactionResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.NotNil(t, resp.Description)
	assert.Equal(t, "оплата кофе", *resp.Description)
}

func TestCreateTransaction_InvalidBody_Returns400(t *testing.T) {
	tests := []struct {
		name string
		body string
	} {
		{name: "invalid json", body: `{not json`},
		{name: "empty client_id", body: `{"client_id":"","business_id":"biz-uuid-1111","amount":500000}`},
		{name: "empty business_id", body: `{"client_id":"cli-uuid-0000","business_id":"","amount":500000}`},
		{name: "empty amount", body: `{"client_id":"cli-uuid-0000","business_id":"biz-uuid-1111","amount":}`},
		{name: "missing client_id", body: `{"business_id":"biz-uuid-1111","amount":500000}`},
		{name: "missing business_id", body: `{"client_id":"cli-uuid-0000","amount":500000}`},
		{name: "missing amount", body: `{"client_id":"cli-uuid-0000","business_id":"biz-uuid-1111"}`},
		{name: "zero amount", body: `{"client_id":"cli-uuid-0000","business_id":"biz-uuid-1111","amount":0}`},
		{name: "negative amount", body: `{"client_id":"cli-uuid-0000","business_id":"biz-uuid-1111","amount":-100000}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := mocks.NewMockTransactionService(t)
			router := newHandler(nil, nil, svc, nil)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/transactions", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusBadRequest, rec.Code)

			var resp errorResponse
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
			assert.NotEmpty(t, resp.Error)
		})
	}
}

func TestCreateTransaction_422Errors(t *testing.T) {
	tests := []struct {
		name string
		svcErr error
	} {
		{name: "client is not active", svcErr: domain.ErrClientIsNotActive},
		{name: "business is not active", svcErr: domain.ErrBusinessIsNotActive},
		{name: "insuffucient bonus balance", svcErr: domain.ErrInsufficientBonusBalance},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := mocks.NewMockTransactionService(t)
			svc.EXPECT().Create(mock.Anything, testClientID, testBusinessID, testAmount, mock.Anything).
				Return(domain.Transaction{}, tt.svcErr).Once()
			router := newHandler(nil, nil, svc, nil)

			body := fmt.Sprintf(`{"client_id":"%s","business_id":"%s","amount":%d}`, testClientID, testBusinessID, testAmount)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/transactions", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

			var resp errorResponse
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
			assert.NotEmpty(t, resp.Error)
		})
	}
}

// GET api/v1/transactions/:{id}
func TestGetTransactionByID(t *testing.T) {
	tests := []struct {
		name string
		svcReturn domain.Transaction
		svcErr error
		wantCode int
		wantAmount int64
	} {
		{
			name: "success",
			svcReturn: testTransaction(),
			wantCode: http.StatusOK,
			wantAmount: testAmount,
		},
		{
			name: "not found",
			svcErr: domain.ErrTransactionNotFound,
			wantCode: http.StatusNotFound,
		},
		{
			name: "internal error",
			svcErr: errors.New("db timeout"),
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := mocks.NewMockTransactionService(t)
			svc.EXPECT().GetByID(mock.Anything, testTransactionID).Return(tt.svcReturn, tt.svcErr).Once()
			router := newHandler(nil, nil, svc, nil)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/transactions/"+testTransactionID, nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantCode, rec.Code)

			if tt.wantAmount != 0 {
				var resp dto.TransactionResponse
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
				assert.NotEmpty(t, resp.ID)
				assert.Equal(t, tt.wantAmount, resp.Amount)
			}
		})
	}
}

func TestGetTransactionsByClientID_Pagination(t *testing.T) {
	tests := []struct {
		name string
		queryParams string
		wantLimit int
		wantOffset int
		svcReturn []domain.Transaction
		wantCode int
		wantLen int
		callService bool
	} {
		{
			name: "default pagination",
			queryParams: "?client_id="+testClientID,
			wantLimit: 20,
			wantOffset: 0,
			svcReturn: makeTransactions(3),
			wantCode: http.StatusOK,
			wantLen: 3,
			callService: true,
		},
		{
			name: "custom pagination",
			queryParams: "?client_id="+testClientID+"&limit=5&offset=10",
			wantLimit: 5,
			wantOffset: 10,
			svcReturn: makeTransactions(5),
			wantCode: http.StatusOK,
			wantLen: 5,
			callService: true,
		},
		{
			name: "empty result returnts empty array not null",
			queryParams: "?client_id="+testClientID,
			wantLimit: 20,
			wantOffset: 0,
			svcReturn: []domain.Transaction{},
			wantCode: http.StatusOK,
			wantLen: 0,
			callService: true,
		},
		{
			name: "missing client id returns 400",
			queryParams: "",
			wantCode: http.StatusBadRequest,
			callService: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			
			svc := mocks.NewMockTransactionService(t)
			if tt.callService {
				svc.EXPECT().GetByClientID(mock.Anything, testClientID, tt.wantLimit, tt.wantOffset).
					Return(tt.svcReturn, nil).Once()
			}
			router := newHandler(nil, nil, svc, nil)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/transactions"+tt.queryParams, nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantCode, rec.Code)

			if tt.wantCode == http.StatusOK {
				var resp []dto.TransactionResponse
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
				assert.NotNil(t, resp)
				assert.Len(t, resp, tt.wantLen)
			}
		})
	}
}

// PATCH /api/v1/transactions/"+{id}+"/cancel
func TestCancelTransaction_Success_Return200(t *testing.T){
	svc := mocks.NewMockTransactionService(t)
	svc.EXPECT().Cancel(mock.Anything, testTransactionID).Return(nil).Once()
	router := newHandler(nil, nil, svc, nil)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/transactions/"+testTransactionID+"/cancel", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var resp messageResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.NotEmpty(t, resp.Message)
}

func TestCancelTransaction_AlreadyCancelled_Returns409(t *testing.T) {
	svc := mocks.NewMockTransactionService(t)
	svc.EXPECT().Cancel(mock.Anything, testTransactionID).Return(domain.ErrTransactionCancelled).Once()
	router := newHandler(nil, nil, svc, nil)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/transactions/"+testTransactionID+"/cancel", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestCancelTransaction_NotFound_Returns404(t *testing.T) {
	svc := mocks.NewMockTransactionService(t)
	svc.EXPECT().Cancel(mock.Anything, testTransactionID).Return(domain.ErrTransactionNotFound).Once()
	router := newHandler(nil, nil, svc, nil)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/transactions/"+testTransactionID+"/cancel", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}