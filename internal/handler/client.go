package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/GulzhanKarakul/payment-service/internal/dto"
	"github.com/go-chi/chi/v5"
)

type createClientRequest struct {
	Phone string `json:"phone"`
	Name  string `json:"name"`
}

func (req createClientRequest) Validate() error {
	if req.Phone == "" {
		return errors.New("phone is required")
	}
	if len(req.Phone) < 10 || len(req.Phone) > 16 {
		return errors.New("phone must be at least 10 characters")
	}
	if req.Name == "" {
		return errors.New("name is required")
	}
	return nil
}

// createClient - POST api/v1/clients/
func (h *Handler) createClient(w http.ResponseWriter, r *http.Request) {
	var req createClientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := req.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	client, err := h.client.Create(r.Context(), req.Phone, req.Name)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, dto.ToClientResponse(client))
}

// getClientByID - GET api/v1/clients/:{id}
func (h *Handler) getClientByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	client, err := h.client.GetByID(r.Context(), id)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.ToClientResponse(client))
}

// getByPhone - GET api/v1/clients/phone/:{phone}
func (h *Handler) getClientByPhone(w http.ResponseWriter, r *http.Request) {
	phone := chi.URLParam(r, "phone")
	if phone == "" || len(phone) < 10 {
		writeError(w, http.StatusBadRequest, "phone is required")
		return
	}

	client, err := h.client.GetByPhone(r.Context(), phone)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.ToClientResponse(client))
}
