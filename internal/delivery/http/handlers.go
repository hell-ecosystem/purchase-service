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
	mux.HandleFunc("/purchases/stats", h.getStats)

	// новые ручки
	mux.HandleFunc("/purchases/", h.routePurchaseActions) // общий маршрут для /purchases/{id}/...

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

func (h *Handler) routePurchaseActions(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	id := strings.TrimPrefix(path, "/purchases/")
	id = strings.SplitN(id, "/", 2)[0]
	if id == "" {
		h.respond(w, http.StatusBadRequest, map[string]string{"error": "missing purchase ID"})
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		h.respond(w, http.StatusUnauthorized, map[string]string{"error": "missing user ID"})
		return
	}

	switch {
	case strings.HasSuffix(path, "/purchase") && r.Method == http.MethodPatch:
		h.markAsPurchased(w, r, userID, id)
	case strings.HasSuffix(path, "/deactivate") && r.Method == http.MethodPatch:
		h.deactivate(w, r, userID, id)
	case strings.HasSuffix(path, "/reminder") && r.Method == http.MethodPost:
		h.addReminder(w, r, userID, id)
	case strings.HasSuffix(path, "/tags") && r.Method == http.MethodPost:
		h.addTags(w, r, userID, id)
	case r.Method == http.MethodPut:
		h.updatePurchase(w, r, userID, id)
	default:
		h.respond(w, http.StatusNotFound, map[string]string{"error": "unknown action"})
	}
}

func (h *Handler) markAsPurchased(w http.ResponseWriter, r *http.Request, userID, purchaseID string) {
	if err := h.service.MarkAsPurchased(r.Context(), userID, purchaseID); err != nil {
		h.respond(w, http.StatusInternalServerError, map[string]string{"error": "could not mark purchase as purchased"})
		return
	}
	h.respond(w, http.StatusOK, map[string]string{"status": "marked as purchased"})
}

func (h *Handler) deactivate(w http.ResponseWriter, r *http.Request, userID, purchaseID string) {
	if err := h.service.Deactivate(r.Context(), userID, purchaseID); err != nil {
		h.respond(w, http.StatusInternalServerError, map[string]string{"error": "could not deactivate purchase"})
		return
	}
	h.respond(w, http.StatusOK, map[string]string{"status": "purchase deactivated"})
}

func (h *Handler) updatePurchase(w http.ResponseWriter, r *http.Request, userID, purchaseID string) {
	var update domain.PurchaseUpdate
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		h.respond(w, http.StatusBadRequest, map[string]string{"error": "invalid input"})
		return
	}
	if err := h.service.Update(r.Context(), userID, purchaseID, update); err != nil {
		h.respond(w, http.StatusInternalServerError, map[string]string{"error": "could not update purchase"})
		return
	}
	h.respond(w, http.StatusOK, map[string]string{"status": "purchase updated"})
}

func (h *Handler) addReminder(w http.ResponseWriter, r *http.Request, userID, purchaseID string) {
	var reminder domain.Reminder
	if err := json.NewDecoder(r.Body).Decode(&reminder); err != nil {
		h.respond(w, http.StatusBadRequest, map[string]string{"error": "invalid input"})
		return
	}
	if err := h.service.AddReminder(r.Context(), userID, purchaseID, reminder); err != nil {
		h.respond(w, http.StatusInternalServerError, map[string]string{"error": "could not add reminder"})
		return
	}
	h.respond(w, http.StatusCreated, map[string]string{"status": "reminder added"})
}

func (h *Handler) addTags(w http.ResponseWriter, r *http.Request, userID, purchaseID string) {
	var payload struct {
		TagIDs []string `json:"tag_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.respond(w, http.StatusBadRequest, map[string]string{"error": "invalid input"})
		return
	}
	if err := h.service.AddTags(r.Context(), userID, purchaseID, payload.TagIDs); err != nil {
		h.respond(w, http.StatusInternalServerError, map[string]string{"error": "could not add tags"})
		return
	}
	h.respond(w, http.StatusCreated, map[string]string{"status": "tags added"})
}

func (h *Handler) getStats(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		h.respond(w, http.StatusUnauthorized, map[string]string{"error": "missing user ID"})
		return
	}
	stats, err := h.service.GetStatistics(r.Context(), userID)
	if err != nil {
		h.respond(w, http.StatusInternalServerError, map[string]string{"error": "could not fetch stats"})
		return
	}
	h.respond(w, http.StatusOK, stats)
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
