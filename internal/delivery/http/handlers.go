package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/hell-ecosystem/purchase-service/internal/app"
	"github.com/hell-ecosystem/purchase-service/internal/domain"
)

type Handler struct {
	service app.PurchaseService
	logger  *slog.Logger
}

func NewHandler(service app.PurchaseService, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/purchases", h.routePurchases)
	return h.loggingMiddleware(mux)
}

func (h *Handler) routePurchases(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.create(w, r)
	case http.MethodGet:
		h.getAll(w, r)
	default:
		h.logger.Warn("method not allowed", slog.String("method", r.Method), slog.String("path", r.URL.Path))
		h.respond(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		h.logger.Warn("unsupported content type", slog.String("content_type", r.Header.Get("Content-Type")))
		h.respond(w, http.StatusUnsupportedMediaType, map[string]string{"error": "Content-Type must be application/json"})
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		h.logger.Warn("missing user ID")
		h.respond(w, http.StatusUnauthorized, map[string]string{"error": "missing user ID"})
		return
	}

	var items []domain.Purchase
	err := json.NewDecoder(r.Body).Decode(&items)
	if err != nil {
		h.logger.Error("invalid JSON body", slog.String("error", err.Error()))
		h.respond(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}

	if len(items) == 0 {
		h.logger.Warn("empty purchase list")
		h.respond(w, http.StatusBadRequest, map[string]string{"error": "no items provided"})
		return
	}

	for _, item := range items {
		if strings.TrimSpace(item.Name) == "" || strings.TrimSpace(item.Category) == "" {
			h.logger.Warn("validation failed", slog.Any("item", item))
			h.respond(w, http.StatusBadRequest, map[string]string{"error": "name and category are required"})
			return
		}
	}

	err = h.service.Create(r.Context(), userID, items)
	if err != nil {
		h.logger.Error("failed to create purchases", slog.String("error", err.Error()))
		h.respond(w, http.StatusInternalServerError, map[string]string{"error": "failed to create purchases"})
		return
	}

	h.logger.Info("purchases created", slog.Int("count", len(items)), slog.String("user_id", userID))
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) getAll(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		h.logger.Warn("missing user ID")
		h.respond(w, http.StatusUnauthorized, map[string]string{"error": "missing user ID"})
		return
	}

	items, err := h.service.GetAll(r.Context(), userID)
	if err != nil {
		h.logger.Error("failed to retrieve purchases", slog.String("error", err.Error()))
		h.respond(w, http.StatusInternalServerError, map[string]string{"error": "failed to retrieve purchases"})
		return
	}

	h.logger.Info("purchases fetched", slog.Int("count", len(items)), slog.String("user_id", userID))
	h.respond(w, http.StatusOK, items)
}

func (h *Handler) respond(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.logger.Info("incoming request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("remote", r.RemoteAddr),
		)
		next.ServeHTTP(w, r)
	})
}
