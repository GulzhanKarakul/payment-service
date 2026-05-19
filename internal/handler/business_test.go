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

// Post /api/v1/businesses create business
func TestCreateBusiness_Success_Returns201(t *testing.T) {
	svc := mocks.NewMockBusinessService(t)
	svc.EXPECT().Create(mock.Anything, testName, testPhone).
		Return(testBusiness(), nil).Once()
	router := newHandler(nil, svc, nil, nil)

	body := `{"name":"`+testName+`","owner_phone":"`+testPhone+`"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/businesses", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var resp dto.BusinessResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, testBusiness().Name, resp.Name)
	assert.Equal(t, testPhone, resp.OwnerPhone)
	assert.NotEmpty(t, resp.ID)
	assert.True(t, resp.IsActive)
}

func TestCreateBusiness_InvalidBody_Returns400(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{name: "invalid json", body: `{"not json"`},
		{name: "empty owner phone", body: `{"name":"Gulzhan","owner_phone":""}`},
		{name: "missing owner phone", body: `{"name":"Gulzhan"}`},
		{name: "short owner phone", body: `{"name":"Gulzhan","owner_phone":"+7777"}`},
		{name: "empty name", body: `{"name":"","owner_phone":"+77771156580"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := mocks.NewMockBusinessService(t)
			router := newHandler(nil, svc, nil, nil)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/businesses", strings.NewReader(tt.body))
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

func TestCreateBusiness_AlreadyExist_Returns409(t *testing.T) {
	svc := mocks.NewMockBusinessService(t)
	svc.EXPECT().Create(mock.Anything, testName, testPhone).
		Return(domain.Business{}, domain.ErrBusinessAlreadyExist).Once()
	router := newHandler(nil, svc, nil, nil)

	body := `{"name":"`+testName+`","owner_phone":"`+testPhone+`"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/businesses", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusConflict, rec.Code)

	var resp errorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.NotEmpty(t, resp.Error)
}

func TestCreateBusiness_InternalError_Returns500(t *testing.T) {
	svc := mocks.NewMockBusinessService(t)
	svc.EXPECT().Create(mock.Anything, testName, testPhone).
		Return(domain.Business{}, errors.New("db connection lost")).Once()
	router := newHandler(nil, svc, nil, nil)

	body := `{"name":"`+testName+`","owner_phone":"`+testPhone+`"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/businesses", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	var resp errorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, "internal server error", resp.Error)
}

// Get /api/v1/businesses/{id} get business by id
func TestGetBusinessByID(t *testing.T) {
	tests := []struct {
		name string
		svcReturn domain.Business
		svcErr error
		wantCode int
		wantPhone string
	}{
		{
			name: "success",
			svcReturn: testBusiness(),
			wantCode: http.StatusOK,
			wantPhone: testPhone,
		},
		{
			name: "not found",
			svcErr: domain.ErrBusinessNotFound,
			wantCode: http.StatusNotFound,
		},
		{
			name: "internal error",
			svcErr: errors.New("db error"),
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := mocks.NewMockBusinessService(t)
			svc.EXPECT().GetByID(mock.Anything, testBusinessID).
				Return(tt.svcReturn, tt.svcErr).Once()
			router := newHandler(nil, svc, nil, nil)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/businesses/"+testBusinessID, nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantCode, rec.Code)

			if tt.wantPhone != "" {
				var resp dto.BusinessResponse
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
				assert.Equal(t, tt.wantPhone, resp.OwnerPhone)
			}
		})
	}
}

// Post /api/v1/businesses/{id}/balance Update business balance
func TestUpdateBusinessBalance_Success_Returns200(t *testing.T) {
	svc := mocks.NewMockBusinessService(t)
	svc.EXPECT().UpdateBonusBalance(mock.Anything, testBusinessID, testBalance).
		Return(testBusiness(), nil).Once()
	router := newHandler(nil, svc, nil, nil)

	body := fmt.Sprintf(`{"amount":%d}`, testBalance)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/businesses/"+testBusinessID+"/balance", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp dto.BusinessResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, testBalance, resp.BonusBalance)
}

func TestBusinessHandler_UpdateBonusBalance_InvalidBody_Returns400(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{name: "invalid json", body: `{"not json`},
		{name: "empty amount", body: `{"amount":""}`},
		{name: "zero amount", body: `{"amount":0}`},
		{name: "negative amount", body: `{"amount":-2000}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := mocks.NewMockBusinessService(t)
			router := newHandler(nil, svc, nil, nil)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/businesses/"+testBusinessID+"/balance", strings.NewReader(tt.body))
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

func TestUpdateBonusBalance_NotFound_Returns404(t *testing.T) {
	svc := mocks.NewMockBusinessService(t)
	svc.EXPECT().UpdateBonusBalance(mock.Anything, testBusinessID, testBalance).
		Return(domain.Business{}, domain.ErrBusinessNotFound).Once()
	router := newHandler(nil, svc, nil, nil)

	body := fmt.Sprintf(`{"amount":%d}`, testBalance)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/businesses/"+testBusinessID+"/balance", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	var resp errorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.NotEmpty(t, resp.Error)
}