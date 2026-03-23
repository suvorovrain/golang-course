package domain

type Repo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Stars       int64  `json:"stargazers_count"`
	Forks       int64  `json:"forks_count"`
	CreatedAt   string `json:"created_at"`
}
