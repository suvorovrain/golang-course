package domain

type Repository struct {
	Owner       string
	Repo        string
	FullName    string
	Description string
	Stars       int32
	Forks       int32
	CreatedAt   string
	Visibility  string
}

type Subscription struct {
	Owner string
	Repo  string
}
