package entity

import "time"

type IPList struct {
	ID        string
	IP        string
	ListType  string 
	CreatorID string
	CreatedAt time.Time
}
