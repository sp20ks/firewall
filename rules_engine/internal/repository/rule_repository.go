package repository

import "rules-engine/internal/entity"

type RuleRepository interface {
	GetActiveRules() ([]entity.Rule, error)
	CreateRule(rule *entity.Rule) (*entity.Rule, error)
	UpdateRule(rule *entity.Rule) (*entity.Rule, error)
	GetRule(id string) (*entity.Rule, error)
	GetRulesForResource(id string) ([]entity.Rule, error)
	GetRulesByURL(url, method string) ([]entity.Rule, error)
}
