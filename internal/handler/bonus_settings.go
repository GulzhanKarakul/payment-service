package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/GulzhanKarakul/payment-service/internal/dto"
	"github.com/go-chi/chi/v5"
)

type upsertBonusSettingsRequest struct {
	BonusPercent float64 `json:"bonus_percent"`
	IsActive     bool    `json:"is_active"`
}

func (req upsertBonusSettingsRequest) Validate() error {
	if req.BonusPercent <= 0 {
		return errors.New("bonus percent must be positive")
	}
	return nil
}


// Upsert POST api/v1/businesses/:{id}/settings
func (h *Handler) upsertBonusSettings(w http.ResponseWriter, r *http.Request) {
	businessID := chi.URLParam(r, "id")
	if businessID == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	var req upsertBonusSettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := req.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	bonusSettings, err := h.bonusSettings.Upsert(r.Context(), businessID, req.BonusPercent, req.IsActive)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.ToBonusSettingsResponse(bonusSettings))
}

// GetByBusinessID - GET api/v1/businesses/:{id}/settings
func (h *Handler) getBonusSettingsByBusinessID(w http.ResponseWriter, r *http.Request) {
	businessID := chi.URLParam(r, "id")
	if businessID == "" {
		writeError(w, http.StatusBadRequest, "business id is required")
		return
	}

	bonusSettings, err := h.bonusSettings.GetByBusinessID(r.Context(), businessID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.ToBonusSettingsResponse(bonusSettings))
}
