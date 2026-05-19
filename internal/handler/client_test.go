package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/GulzhanKarakul/payment-service/internal/domain"
	"github.com/GulzhanKarakul/payment-service/internal/dto"
	"github.com/GulzhanKarakul/payment-service/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Post /api/v1/clients - create client
func TestCreateClient_Success_Returns201(t *testing.T) {
	svc := mocks.NewMockClientService(t)
	svc.EXPECT().Create(mock.Anything, testPhone, testName).
		Return(testClient(), nil).Once()
	router := newHandler(svc, nil, nil, nil)

	body := `{"phone":"`+testPhone+`","name":"`+testName+`"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/clients", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var resp dto.ClientResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, testPhone, resp.Phone)
	assert.Equal(t, testName, resp.Name)
	assert.NotEmpty(t, resp.ID)
	assert.True(t, resp.IsActive)
}

// валидные ошибки: инвалидный боди (тдд), ответ 400 всегда, сервис не вызывается
func TestCreateClient_InvalidBody_Returns400(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{name: "invalid json", body: `{not json`},
		{name: "empty phone", body: `{"phone":"","name":"Gulzhan"}`},
		{name: "missing phone", body: `{"name":"Gulzhan"}`},
		{name: "short phone", body: `{"phone":"+7777","name":"Gulzhan"}`},
		{name: "empty name", body: `{"phone":"+77771156580","name":""}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := mocks.NewMockClientService(t)
			router := newHandler(svc, nil, nil, nil)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/clients", strings.NewReader(tt.body))
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

func TestCreateClient_AlreadyExist_Returns409(t *testing.T) {
	svc := mocks.NewMockClientService(t)
	svc.EXPECT().Create(mock.Anything, testPhone, testName).
		Return(domain.Client{}, domain.ErrClientAlreadyExist).Once()
	router := newHandler(svc, nil, nil, nil)

	body := `{"phone":"`+testPhone+`","name":"`+testName+`"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/clients", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusConflict, rec.Code)

	var resp errorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.NotEmpty(t, resp.Error)
}

// Get /api/v1/clients/{id} get cliert by id
func TestGetClientByID(t *testing.T) {
	tests := []struct {
		name string
		svcReturn domain.Client
		svcErr error
		wantCode int
		wantPhone string
	}{
		{
			name: "success",
			svcReturn: testClient(),
			wantCode: http.StatusOK,
			wantPhone: testPhone,
		},
		{
			name: "not found",
			svcErr: domain.ErrClientNotFound,
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
			svc := mocks.NewMockClientService(t)
			svc.EXPECT().GetByID(mock.Anything, testClientID).
				Return(tt.svcReturn, tt.svcErr).Once()
			router := newHandler(svc, nil, nil, nil)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/clients/"+testClientID, nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantCode, rec.Code)
			
			if tt.wantPhone != "" {
				var resp dto.ClientResponse
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
				assert.Equal(t, tt.wantPhone, resp.Phone)
			}
		})
	}
}

// Get /api/v1/clients/phone/{phone}
func TestGetClientByPhone(t *testing.T) {
	tests := []struct {
		name string
		svcReturn domain.Client
		svcErr error
		wantCode int
	}{
		{
			name: "success",
			svcReturn: testClient(),
			wantCode: http.StatusOK,
		},
		{
			name: "not found",
			svcErr: domain.ErrClientNotFound,
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := mocks.NewMockClientService(t)
			svc.EXPECT().GetByPhone(mock.Anything, testPhone).
				Return(tt.svcReturn, tt.svcErr).Once()
			router := newHandler(svc, nil, nil, nil)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/clients/phone/"+url.PathEscape(testPhone), nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantCode, rec.Code)
		})
	}
}