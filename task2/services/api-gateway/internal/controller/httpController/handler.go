package httpController

import (
	"encoding/json"
	"golang-course/task2/services/api-gateway/internal/usecase"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	usecase *usecase.Usecase
}

func New(usecase *usecase.Usecase) *Handler {
	return &Handler{usecase: usecase}
}

// GetRepo godoc
// @Summary      Get repository info
// @Description  Get details about a GitHub repository by owner and name
// @Tags         repository
// @Accept       json
// @Produce      json
// @Param        owner   query      string  true  "Repository Owner"
// @Param        name    query      string  true  "Repository Name"
// @Success      200     {object}   domain.Repo
// @Failure      400     {string}   string  "Bad Request"
// @Failure      404     {string}   string  "Not Found"
// @Failure      500     {string}   string  "Internal Server Error"
// @Router       /repo [get]
func (h *Handler) GetRepo(w http.ResponseWriter, r *http.Request) {
	owner := r.URL.Query().Get("owner")
	name := r.URL.Query().Get("name")

	if owner == "" || name == "" {
		http.Error(w, "missing owner/name", http.StatusBadRequest)
		return
	}

	repo, err := h.usecase.GetRepoInfo(r.Context(), owner, name)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				http.Error(w, "not found", http.StatusNotFound)
			case codes.InvalidArgument:
				http.Error(w, "invalid argument", http.StatusBadRequest)
			default:
				http.Error(w, "server error", http.StatusInternalServerError)
			}
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(repo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
