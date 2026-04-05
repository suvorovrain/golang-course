package domain


type RepoInfo struct {
    Name        string `json:"name"`
    FullName    string `json:"full_name"`
    Description string `json:"description"`
    URL         string `json:"url"`
    Stars       int32  `json:"stars"`
    Forks       int32  `json:"forks"`
    Watchers    int32  `json:"watchers"`
    Language    string `json:"language"`
    CreatedAt   string `json:"created_at"`
    UpdatedAt   string `json:"updated_at"`
}