package grpc

import (
    "context"
    "github-info-system/api/generated"
    "github-info-system/collector/internal/usecase"
)
type Server struct{ 
	generated.UnimplementedCollectorServiceServer
	
	
	
	useCase *usecase.GetRepoUseCase
}
func NewServer() *Server{ 
	return &Server{
		useCase: usecase.NewGetRepoUseCase(),
		
	}
}
func (s *Server) GetRepoInfo(ctx context.Context, req *generated.RepoRequest) (*generated.RepoResponse, error){
	repo, err := s.useCase.Execute(req.Owner, req.Repo)
	if err!= nil{
		return &generated.RepoResponse{
			Error : err.Error(),
		}, nil
	}
	    return &generated.RepoResponse{
        Name:        repo.Name,
        FullName:    repo.FullName,
        Description: repo.Description,
        Url:         repo.URL,
        Stars:       repo.Stars,
        Forks:       repo.Forks,
        Watchers:    repo.Watchers,
        Language:    repo.Language,
        CreatedAt:   repo.CreatedAt,
        UpdatedAt:   repo.UpdatedAt,
    }, nil
}
