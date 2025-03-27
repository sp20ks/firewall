package entity

import (
	"net"
	"time"
)

type IPList struct {
	ID        string    `json:"id"`
	IP        net.IPNet `json:"ip"`
	ListType  string    `json:"list_type"`
	CreatorID string    `json:"creator_id"`
	CreatedAt time.Time `json:"created_at"`
}
