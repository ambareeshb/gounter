package handler

import (
	"context"
	"encoding/json"
	"gounter/internal/model"
	"net/http"

	"github.com/google/uuid"
)

type Service interface {
	CreateCounter(ctx context.Context, name string) (*model.Counter, error)
	IncrementCounter(ctx context.Context, id uuid.UUID) (*model.Counter, error)
	SoftDeleteCounter(ctx context.Context, id uuid.UUID) (int64, error)
}

type Handler struct {
	service Service
}

// NewHandler for creating new handler
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// CreateCounter handles counter creation
func (h *Handler) CreateCounter(w http.ResponseWriter, r *http.Request) {
	var counter *model.Counter
	err := json.NewDecoder(r.Body).Decode(&counter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	counter, err = h.service.CreateCounter(r.Context(), counter.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(counter)
}

// IncrementCounter handles incrementing a counter
func (h *Handler) IncrementCounter(w http.ResponseWriter, r *http.Request) {
	var counter *model.Counter
	err := json.NewDecoder(r.Body).Decode(&counter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	counter, err = h.service.IncrementCounter(r.Context(), counter.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(counter)
}

// DeleteCounter handles deleting a counter
func (h *Handler) DeleteCounter(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL (this would require a URL parameter setup)
	idString := r.URL.Query().Get("id")

	uuid, err := uuid.Parse(idString)
	if err != nil {
		http.Error(w, "Please provide valid uuid", http.StatusBadRequest)
		return
	}

	_, err = h.service.SoftDeleteCounter(r.Context(), uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
