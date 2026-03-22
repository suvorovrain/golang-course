package domain

import (
	"fmt"
	"time"
)

type Repository struct {
	Name        string    `json:"name" example:"go-github"`
	Description string    `json:"description" example:"Go library for accessing the GitHub API"`
	Stars       int       `json:"stars" example:"12345"`
	Forks       int       `json:"forks" example:"2345"`
	CreatedAt   time.Time `json:"created_at" example:"2013-05-24T16:22:20Z"`
}

func (r Repository) String() string {
	return fmt.Sprintf(
		"Name: %s\nDescription: %s\nStars: %d\nForks: %d\nCreated: %s",
		r.Name, r.Description, r.Stars, r.Forks, r.CreatedAt,
	)
}
