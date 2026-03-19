package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"task3/api/internal/dto"
	"task3/api/internal/usecase"
)

func NewPingHandler(log *slog.Logger, ping *usecase.Ping) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := ping.Execute(r.Context())

		response := dto.PingResponse{
			Reply: string(status),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Error("failed to write ping response", "error", err)
		}
	}
}
