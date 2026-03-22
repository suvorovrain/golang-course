package domain

type RepoInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	StarsCount  int    `json:"stargazers_count"`
	ForksCount  int    `json:"forks_count"`
	CreatedAt   string `json:"created_at"`
}
