package repository

import "rules-engine/internal/entity"

type RuleRepository interface {
	GetActiveRules() ([]entity.Rule, error)
	CreateRule(rule *entity.Rule) error
	UpdateRule(rule *entity.Rule) error
	GetRule(id string) (*entity.Rule, error)
}
