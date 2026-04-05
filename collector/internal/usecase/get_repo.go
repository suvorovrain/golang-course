package usecase
import (
    "fmt"
    "github-info-system/collector/internal/adapter/github"
    "github-info-system/collector/internal/domain"
)
type GetRepoUseCase struct{
	githubClient *github.GitHubClient 
}

func NewGetRepoUseCase() *GetRepoUseCase{
	return &GetRepoUseCase{
		githubClient: github.NewGitHubClient(),
		
	}
}
func (uc *GetRepoUseCase) Execute(owner, repo string) (*domain.RepoInfo, error){
	 if owner == "" || repo == "" {
        return nil, fmt.Errorf("owner и repo обязательны")
    }
	ghRepo, err := uc.githubClient.GetRepo(owner, repo) 
	if err!= nil{
		return nil, err
	}

	return &domain.RepoInfo{
		Name:        ghRepo.Name,
        FullName:    ghRepo.FullName,
        Description: ghRepo.Description,
        URL:         ghRepo.HTMLURL,
        Stars:       int32(ghRepo.Stargazers),
        Forks:       int32(ghRepo.Forks),
        Watchers:    int32(ghRepo.Watchers),
        Language:    ghRepo.Language,
        CreatedAt:   ghRepo.CreatedAt,
        UpdatedAt:   ghRepo.UpdatedAt,

	}, nil
}