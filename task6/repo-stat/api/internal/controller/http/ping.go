package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"repo-stat/api/internal/domain"
	"repo-stat/api/internal/dto"
	"repo-stat/api/internal/usecase"
)

// NewPingHandler godoc
// @Summary      Проверка здоровья сервисов
// @Description  Отправляет ping-запросы в Processor и Subscriber и возвращает статус каждого сервиса
// @Tags         Ping
// @Produce      json
// @Success      200  {object}  dto.PingResponse  "Все сервисы работают"
// @Failure      503  {object}  dto.PingResponse  "Один или несколько сервисов недоступны"
// @Router       /api/ping [get]
func NewPingHandler(log *slog.Logger, pingUC *usecase.Ping) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		processorStatus := pingUC.ProcessorPing(r.Context())
		subscriberStatus := pingUC.SubscriberPing(r.Context())

		services := []dto.ServiceStatus{
			{Name: "processor", Status: string(processorStatus)},
			{Name: "subscriber", Status: string(subscriberStatus)},
		}

		overallStatus := "ok"
		httpCode := http.StatusOK

		if processorStatus == domain.PingStatusDown || subscriberStatus == domain.PingStatusDown {
			overallStatus = "degraded"
			httpCode = http.StatusServiceUnavailable
		}

		response := dto.PingResponse{
			Status:   overallStatus,
			Services: services,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(httpCode)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Error("failed to encode ping response", "error", err)
			http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
			return
		}
	}
}
