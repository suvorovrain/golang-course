package handler

import (
	"net/http"
	"task2/gateway/internal/adapter"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HTTPHandler struct {
	adapter *adapter.GRPCCollectorAdapter
}

func NewHTTPHandler(a *adapter.GRPCCollectorAdapter) *HTTPHandler {
	return &HTTPHandler{
		adapter: a,
	}
}

// GetRepo godoc
// @Summary      Получить информацию о репозитории
// @Tags         repositories
// @Produce      json
// @Param        owner   path      string  true  "Имя владельца"
// @Param        repo    path      string  true  "Название репозитория"
// @Success      200  {object}  domain.Repository
// @Failure      404  {object}  map[string]string
// @Router       /api/repo/{owner}/{repo} [get]
func (h *HTTPHandler) GetRepo(c *gin.Context) {
	owner := c.Param("owner")
	repo := c.Param("repo")

	res, err := h.adapter.GetRepository(c.Request.Context(), owner, repo)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{"error": st.Message()})
				return
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid arguments"})
				return
			}
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
