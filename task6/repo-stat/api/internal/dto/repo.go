package dto

type RepoResponse struct {
	Owner       string `json:"owner"`
	Repo        string `json:"repo"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Stars       int32  `json:"stars"`
	Forks       int32  `json:"forks"`
	CreatedAt   string `json:"created_at"`
	Visibility  string `json:"visibility"`
}

type SubscriptionInfoResponse struct {
	Repositories []RepoResponse `json:"repositories"`
}
