package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/GulzhanKarakul/payment-service/internal/dto"
	"github.com/go-chi/chi/v5"
)

type createBusinessRequest struct {
	Name       string `json:"name"`
	OwnerPhone string `json:"owner_phone"`
}

func (req createBusinessRequest) Validate() error {
	if req.Name == "" {
		return errors.New("name is required")
	}
	if req.OwnerPhone == "" {
		return errors.New("phone is required")
	}
	if len(req.OwnerPhone) < 10 {
		return errors.New("phone must be 10 characters")
	}
	return nil
}

// Create - POST api/v1/businesses/
func (h *Handler) createBusiness(w http.ResponseWriter, r *http.Request) {
	var req createBusinessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := req.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	business, err := h.business.Create(r.Context(), req.Name, req.OwnerPhone)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, dto.ToBusinessResponse(business))
}

// GetByID - GET api/v1/businesses/:{id}
func (h *Handler) getBusinessByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	business, err := h.business.GetByID(r.Context(), id)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.ToBusinessResponse(business))
}

type updateBonusBalanceRequest struct {
	Amount int64 `json:"amount"`
}

func (req updateBonusBalanceRequest) Validate() error {
	if req.Amount <= 0 {
		return errors.New("amount must be positive number")
	}
	return nil
}

func (h *Handler) updateBonusBalance(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	var req updateBonusBalanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "amount is required")
		return
	}

	if err := req.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	business, err := h.business.UpdateBonusBalance(r.Context(), id, req.Amount)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.ToBusinessResponse(business))
}
