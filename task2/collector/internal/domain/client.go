package domain

import (
	"time"
)

type Repository struct {
	FullName    string
	Description string
	Stars       int
	Forks       int
	CreatedAt   time.Time
	Language    string
	HTMLURL     string
}
