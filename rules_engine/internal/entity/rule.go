package entity

import "time"

type Action string

const (
	ActionBlock    Action = "block"
	ActionSanitize Action = "sanitize"
	ActionEscape   Action = "escape"
	ActionAllow    Action = "allow"
)

type Rule struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	AttackType string    `json:"attack_type"`
	ActionType Action    `json:"action_type"`
	IsActive   *bool     `json:"is_active"`
	CreatorID  string    `json:"creator_id"`
	CreatedAt  time.Time `json:"created_at"`
}
