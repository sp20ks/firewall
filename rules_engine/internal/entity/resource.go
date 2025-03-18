package entity

import "time"

type Resource struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	HTTPMethod string    `json:"http_method"`
	URL        string    `json:"url"`
	Host       string    `json:"host"`
	CreatorID  string    `json:"creator_id"`
	IsActive   *bool     `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
}
