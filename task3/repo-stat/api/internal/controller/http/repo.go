package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
	"repo-stat/api/internal/dto"
	"repo-stat/api/internal/usecase"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewRepoHandler
// @Summary Get info about repo by url
// @Tags Repositories
// @Accept json
// @Produce json
// @Param url query string true "url of github repo"
// @Success 200 {object} dto.RepoResponse "Success"
// @Failure 400 {string} string "Incorrect request"
// @Failure 404 {string} string "Repo not found"
// @Failure 500 {string} string "Server error"
// @Failure 503 {string} string "Service is unavailable"
// @Router  /api/repositories/info [get]
func NewRepoHandler(log *slog.Logger, repo *usecase.RepoUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlParam := r.URL.Query().Get("url")
		if urlParam == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "url is required",
			})
			return
		}

		parsed, err := url.Parse(urlParam)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "url parse error",
			})
			return
		}

		fullName := strings.Trim(parsed.Path, "/")
		if fullName == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "full name is required",
			})
			return
		}

		repoInfo, err := repo.GetRepoInfo(r.Context(), urlParam)
		if err != nil {
			statusCode := http.StatusInternalServerError
			errorMsg := "Internal Server Error"

			if st, ok := status.FromError(err); ok {
				switch st.Code() {
				case codes.NotFound:
					statusCode = http.StatusNotFound
					errorMsg = "repo not found"
				case codes.InvalidArgument:
					statusCode = http.StatusBadRequest
					errorMsg = "invalid url"
				case codes.Unavailable:
					statusCode = http.StatusServiceUnavailable
					errorMsg = "service unavailable"
				}
			} else {
				statusCode = http.StatusBadRequest
				errorMsg = "Invalid url"
			}

			log.Error("get repo info failed", "error", err, "url", urlParam)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(statusCode)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": errorMsg,
			})

			return
		}

		response := dto.FromDomainRepo(repoInfo, fullName)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Error("encode response failed", "error", err)
		}

	}
}
