package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"repo-stat/api/internal/domain"
	"repo-stat/api/internal/dto"
	"repo-stat/api/internal/usecase"
)

func NewPingHandler(log *slog.Logger, ping *usecase.Ping, agu *usecase.ApiGatewayUsecase, eh *dto.ErrorHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		collector_status := ping.Execute(r.Context())
		processor_status := agu.Ping(r.Context())
		w.Header().Set("Content-Type", "application/json")
		response := &dto.PingResponse{}
		if collector_status == domain.PingStatusUp && processor_status == domain.PingStatusUp {
			response = dto.CreatePingResponce("ok", processor_status, collector_status)
			w.WriteHeader(http.StatusOK)
		} else {
			response = dto.CreatePingResponce("degraded", processor_status, collector_status)
			w.WriteHeader(http.StatusServiceUnavailable)
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Error("failed to write ping response", "error", err)
		}

	}
}
