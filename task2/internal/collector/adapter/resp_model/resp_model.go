package collectorrespmodel

type RepoInfo struct {
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Stargazers  uint64 `json:"stargazers"`
	Forks       uint64 `json:"forks"`
	CreatedAt   string `json:"created_at"`
}
