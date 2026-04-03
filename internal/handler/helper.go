package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/GulzhanKarakul/payment-service/internal/domain"
)

type errorResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, errorResponse{Error: msg})
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrClientNotFound),
		errors.Is(err, domain.ErrBusinessNotFound),
		errors.Is(err, domain.ErrBonusSettingsNotFound):
		writeError(w, http.StatusNotFound, "not found")

	case errors.Is(err, domain.ErrClientAlreadyExist),
		errors.Is(err, domain.ErrBusinessAlreadyExist),
		errors.Is(err, domain.ErrTransactionCancelled):
		writeError(w, http.StatusConflict, err.Error())

	case errors.Is(err, domain.ErrClientIsNotActive),
		errors.Is(err, domain.ErrBusinessIsNotActive),
		errors.Is(err, domain.ErrInsufficientBonusBalance):
		writeError(w, http.StatusUnprocessableEntity, err.Error())

	case errors.Is(err, domain.ErrInvalidInput):
		writeError(w, http.StatusBadRequest, err.Error())

	default:
		h.log.Error("internal server error", "error", err.Error())
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}
