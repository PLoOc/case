package api

import (
	"encoding/json"
	"net/http"

	"github.com/PLoOc/case/aggregator/internal/infra/repository"
)

type Handler struct {
	repo *repository.DynamoDBRepository
}

func NewHandler(repo *repository.DynamoDBRepository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handler) GetEvents(w http.ResponseWriter, r *http.Request) {
	developerID := r.PathValue("developer_id")
	events, err := h.repo.GetEvents(r.Context(), developerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func (h *Handler) GetSummary(w http.ResponseWriter, r *http.Request) {
	developerID := r.PathValue("developer_id")
	summary, err := h.repo.GetSummary(r.Context(), developerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}