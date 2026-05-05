package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/GulzhanKarakul/payment-service/internal/domain"
	"github.com/GulzhanKarakul/payment-service/internal/dto"
	"github.com/go-chi/chi/v5"
)

type createTransactionRequest struct {
	BusinessID  string  `json:"business_id"`
	ClientID    string  `json:"client_id"`
	Amount      int64   `json:"amount"`
	Description *string `json:"description"`
}

func (req createTransactionRequest) Validate() error {
	if req.BusinessID == "" {
		return errors.New("business id is required")
	}
	if req.ClientID == "" {
		return errors.New("client id is required")
	}
	if req.Amount <= 0 {
		return errors.New("amount must be positive number")
	}
	return nil
}


// createTransaction - POST api/v1/transactions/
func (h *Handler) createTransaction(w http.ResponseWriter, r *http.Request) {
	var req createTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := req.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	transaction, err := h.transaction.Create(r.Context(), req.ClientID, req.BusinessID, req.Amount, req.Description)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, dto.ToTransactionResponse(transaction))
}

// getTransactionByID - GET api/v1/transactions/:{id}
func (h *Handler) getTransactionByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	transaction, err := h.transaction.GetByID(r.Context(), id)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.ToTransactionResponse(transaction))
}

func toTransactionListResponse(tl []domain.Transaction) []dto.TransactionResponse {
	result := make([]dto.TransactionResponse, 0, len(tl))
	for _, t := range tl {
		result = append(result, dto.ToTransactionResponse(t))
	}
	return result
}

// getTransactionsByClientID - GET api/v1/transactions?client_id=X
func (h *Handler) getTransactionsByClientID(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("client_id")
	if clientID == "" {
		writeError(w, http.StatusBadRequest, "client_id is required")
		return
	}

	limit := 20
	offset := 0
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	transactions, err := h.transaction.GetByClientID(r.Context(), clientID, limit, offset)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toTransactionListResponse(transactions))
}

// cancelTransaction - PATCH api/v1/transactions/:{id}/cancel
func (h *Handler) cancelTransaction(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	err := h.transaction.Cancel(r.Context(), id)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "transaction cancelled"})
}
