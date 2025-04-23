package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/hell-ecosystem/purchase-service/internal/app"
	"github.com/hell-ecosystem/purchase-service/internal/domain"
)

type Handler struct {
	service app.PurchaseService
}

func NewHandler(service app.PurchaseService) *Handler {
	return &Handler{service: service}
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
		h.respond(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		h.respond(w, http.StatusUnsupportedMediaType, map[string]string{"error": "Content-Type must be application/json"})
		return
	}
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		h.respond(w, http.StatusUnauthorized, map[string]string{"error": "missing user ID"})
		return
	}
	var items []domain.Purchase
	err := json.NewDecoder(r.Body).Decode(&items)
	if err != nil {
		h.respond(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	if len(items) == 0 {
		h.respond(w, http.StatusBadRequest, map[string]string{"error": "no items provided"})
		return
	}
	for _, item := range items {
		if strings.TrimSpace(item.Name) == "" || strings.TrimSpace(item.Category) == "" {
			h.respond(w, http.StatusBadRequest, map[string]string{"error": "name and category are required"})
			return
		}
	}
	err = h.service.Create(r.Context(), userID, items)
	if err != nil {
		h.respond(w, http.StatusInternalServerError, map[string]string{"error": "failed to create purchases"})
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) getAll(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		h.respond(w, http.StatusUnauthorized, map[string]string{"error": "missing user ID"})
		return
	}
	items, err := h.service.GetAll(r.Context(), userID)
	if err != nil {
		h.respond(w, http.StatusInternalServerError, map[string]string{"error": "failed to retrieve purchases"})
		return
	}
	h.respond(w, http.StatusOK, items)
}

func (h *Handler) respond(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}
