package entity

import "time"

type Rule struct {
	ID         string
	Name       string
	AttackType string
	ActionType string
	IsActive   bool
	CreatorID  string
	CreatedAt  time.Time
}
