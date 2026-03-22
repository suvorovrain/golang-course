package domain

import "time"

type RepoInfo struct {
	Name        string
	Description string
	StarsCount  int
	ForksCount  int
	CreatedAt   time.Time
}
