package http

import (
	"APIGatway/internal/client/grpc"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type Handler struct {
	collectorClient *grpc.CollectorClient
}

func NewHandler(collectorClient *grpc.CollectorClient) *Handler {
	return &Handler{
		collectorClient: collectorClient,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/repositories", h.GetRepository)
}

func (h *Handler) GetRepository(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	repoURL := r.URL.Query().Get("url")
	if repoURL == "" {
		h.sendError(w, "URL parameter is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	resp, err := h.collectorClient.GetRepository(ctx, repoURL)
	if err != nil {
		h.sendError(w, "Failed to call collector service: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if resp.ErrorCode != 0 {
		httpCode := http.StatusInternalServerError
		if resp.ErrorCode == 404 {
			httpCode = http.StatusNotFound
		} else if resp.ErrorCode == 400 {
			httpCode = http.StatusBadRequest
		}
		h.sendError(w, resp.ErrorMessage, httpCode)
		return
	}

	result := map[string]interface{}{
		"full_name":   resp.Repository.FullName,
		"description": resp.Repository.Description,
		"stars":       resp.Repository.Stars,
		"forks":       resp.Repository.Forks,
		"created_at":  resp.Repository.CreatedAt,
		"language":    resp.Repository.Language,
	}

	h.sendJSON(w, http.StatusOK, result)
}

func (h *Handler) sendJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) sendError(w http.ResponseWriter, message string, code int) {
	h.sendJSON(w, code, map[string]string{"error": message})
}
