package dto

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type ErrorHandler struct{}

func (rh *ErrorHandler) CreateErrorResponce(log *slog.Logger, w http.ResponseWriter, code int, error_text string) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": error_text}); err != nil {
		log.Error("failed to write errorInfo", "error", err)
	}

}
