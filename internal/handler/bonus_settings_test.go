package handler_test

import (
	"encoding/json"
	"errors"
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

// Upsert POST api/v1/businesses/:{id}/settings
func TestUpsertBonusSettings_Success_Returns200(t *testing.T) {
	svc := mocks.NewMockBonusSettingsService(t)
	svc.EXPECT().Upsert(mock.Anything, testBusinessID, 5.0, true).
		Return(testBonusSettings(), nil).Once()
	router := newHandler(nil, nil, nil, svc)

	body := `{"bonus_percent":5,"is_active":true}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/businesses/"+testBusinessID+"/settings", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp dto.BonusSettingsResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, testBusinessID, resp.BusinessID)
	assert.NotEmpty(t, resp.ID)
	assert.True(t, resp.IsActive)
}

func TestUpsertBonusSettings_DeactivateSettings_Returns200(t *testing.T) {
	svc := mocks.NewMockBonusSettingsService(t)
	svc.EXPECT().Upsert(mock.Anything, testBusinessID, 5.0, false).
		Return(testInactiveBonusSettings(), nil).Once()
	router := newHandler(nil, nil, nil, svc)

	body := `{"bonus_percent":5,"is_active":false}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/businesses/"+testBusinessID+"/settings", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp dto.BonusSettingsResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.NotEmpty(t, resp.ID)
	assert.False(t, resp.IsActive)
}

func TestUpsertBonusSettings_InvalidBody_Returns400(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{name: "invalid json", body: `{"not json`},
		{name: "zero percent", body: `{"bonus_percent":0,"is_active":true}`},
		{name: "negative percent", body: `{"bonus_percent":-5,"is_active":false}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := mocks.NewMockBonusSettingsService(t)
			router := newHandler(nil, nil, nil, svc)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/businesses/"+testBusinessID+"/settings", strings.NewReader(tt.body))
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

func TestUpsertBonusSettings_BusinessNotFound_Returns404(t *testing.T) {
	svc := mocks.NewMockBonusSettingsService(t)
	svc.EXPECT().Upsert(mock.Anything, testBusinessID, 5.0, true).
		Return(domain.BonusSettings{}, domain.ErrBusinessNotFound).Once()
	router := newHandler(nil, nil, nil, svc)

	body := `{"bonus_percent":5,"is_active":true}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/businesses/"+testBusinessID+"/settings", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	var resp errorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.NotEmpty(t, resp.Error)
}

// GetByBusinessID - GET api/v1/businesses/:{id}/settings
func TestBonusSettingGetByBusinessID(t *testing.T) {
	tests := []struct {
		name string
		svcReturn domain.BonusSettings
		svcErr error
		wantCode int
		wantBusinessID string
	}{
		{
			name: "success",
			svcReturn: testBonusSettings(),
			wantCode: http.StatusOK,
			wantBusinessID: testBusinessID,
		},
		{
			name: "not found",
			svcErr: domain.ErrBonusSettingsNotFound,
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
			svc := mocks.NewMockBonusSettingsService(t)
			svc.EXPECT().GetByBusinessID(mock.Anything, testBusinessID).
				Return(tt.svcReturn, tt.svcErr).Once()
			router := newHandler(nil, nil, nil, svc)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/businesses/"+testBusinessID+"/settings", nil)
			rec := httptest.NewRecorder()
		
			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantCode, rec.Code)

			if tt.wantBusinessID != "" {
				var resp dto.BonusSettingsResponse
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
				assert.Equal(t, tt.wantBusinessID, resp.BusinessID)
			}
		})
	}
}