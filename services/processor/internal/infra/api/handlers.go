package api

import (
	"encoding/json"
	"net/http"

	"github.com/PLoOc/case/aggregator/internal/infra/repository"
)

type APIHandler struct {
	repo *repository.DynamoDBRepository
}

func NewAPIHandler(repo *repository.DynamoDBRepository) *APIHandler {
	return &APIHandler{repo: repo}
}

func (h *APIHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *APIHandler) GetSummary(w http.ResponseWriter, r *http.Request) {
	developerID := r.PathValue("developer_id")

	summary, err := h.repo.GetSummary(r.Context(), developerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}