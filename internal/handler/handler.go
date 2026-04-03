package handler

import (
	"log/slog"
	"net/http"

	"github.com/GulzhanKarakul/payment-service/internal/service"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	client        service.ClientService
	business      service.BusinessService
	transaction   service.TransactionService
	bonusSettings service.BonusSettingsService
	log           *slog.Logger
}

func NewHandler(
	client service.ClientService,
	business service.BusinessService,
	transaction service.TransactionService,
	bonusSettings service.BonusSettingsService,
	log *slog.Logger,
) *Handler {
	return &Handler{
		client:        client,
		business:      business,
		transaction:   transaction,
		bonusSettings: bonusSettings,
		log:           log,
	}
}

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/clients", func(r chi.Router) {
			r.Post("/", h.createClient)
			r.Get("/phone/{phone}", h.getClientByPhone)

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", h.getClientByID)
			})
		})

		r.Route("/businesses", func(r chi.Router) {
			r.Post("/", h.createBusiness)

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", h.getBusinessByID)
				r.Post("/balance", h.updateBonusBalance)
				r.Post("/settings", h.upsertBonusSettings)
				r.Get("/settings", h.getBonusSettingsByBusinessID)
			})
		})

		r.Route("/transactions", func(r chi.Router) {
			r.Post("/", h.createTransaction)
			r.Get("/", h.getTransactionsByClientID)

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", h.getTransactionByID)
				r.Patch("/cancel", h.cancelTransaction)
			})
		})
	})

	return r
}
