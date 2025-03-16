package entity

import "time"

type Resource struct {
	ID         string
	Name       string
	HTTPMethod string
	URL        string
	CreatorID  string
	IsActive   bool
	CreatedAt  time.Time
}
