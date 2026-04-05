package rest

import (
    "encoding/json"
    "net/http"
    "strings"

    "github-info-system/gateway/internal/adapter/grpc"
    "github-info-system/gateway/internal/domain"
)

type Handler struct{
	collectorClient *grpc.CollectorClient
}


func NewHandler(collectorAddr string) (*Handler, error) {
    client, err := grpc.NewCollectorClient(collectorAddr)
	
	
    if err != nil {
        return nil, err
    }

    return &Handler{
        collectorClient: client,
    }, nil
}



func (h *Handler) GetRepoInfo(w http.ResponseWriter, r *http.Request){
	parts := strings.Split(r.URL.Path, "/") 
	
	if len(parts) < 6 {
			http.Error(w, "Неверный формат URL. Ожидается: /api/v1/repos/{owner}/{repo}", http.StatusBadRequest)
			return
		}
	owner := parts[4]
    repo := parts[5]	
	resp, err := h.collectorClient.GetRepoInfo(r.Context(), owner, repo)
	

	if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    
    if resp.Error != "" {
        http.Error(w, resp.Error, http.StatusNotFound)
        return
    }

	result := domain.RepoInfo{
		Name:        resp.Name,
		FullName:    resp.FullName,
		Description: resp.Description,
		URL:         resp.Url,
		Stars:       resp.Stars,
		Forks:       resp.Forks,
		Watchers:    resp.Watchers,
		Language:    resp.Language,
		CreatedAt:   resp.CreatedAt,
		UpdatedAt:   resp.UpdatedAt,
	}
	w.Header().Set("Content-Type", "application/json")
	
	encoder:=json.NewEncoder(w)
	
	encoder.SetIndent("", "  ")
	encoder.Encode(result)
	
}


func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK) 
    w.Write([]byte(`{"status": "ok"}`)) 
}