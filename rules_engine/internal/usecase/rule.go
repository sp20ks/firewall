package usecase

import (
	"fmt"
	"rules-engine/internal/entity"
	"rules-engine/internal/repository"
	"time"
)

type RuleUseCase struct {
	repo repository.RuleRepository
}

func NewRuleUseCase(repo repository.RuleRepository) *RuleUseCase {
	return &RuleUseCase{repo: repo}
}

func (r *RuleUseCase) Get() ([]entity.Rule, error) {
	return r.repo.GetActiveRules()
}

func (r *RuleUseCase) Create(name, attackType, actionType, creatorID string, isActive *bool) error {
	rule := &entity.Rule{
		Name:       name,
		AttackType: attackType,
		ActionType: actionType,
		CreatorID:  creatorID,
		IsActive:   isActive,
		CreatedAt:  time.Now(),
	}

	return r.repo.CreateRule(rule)
}

func (r *RuleUseCase) Update(id, name, attackType, actionType string, isActive *bool) error {
	rule, err := r.repo.GetRule(id)
	if err != nil {
		return fmt.Errorf("error fetching rule: %w", err)
	}

	if rule == nil {
		return fmt.Errorf("rule with id=%s not found", id)
	}

	if name != "" {
		rule.Name = name
	}
	if attackType != "" {
		rule.AttackType = attackType
	}
	if actionType != "" {
		rule.ActionType = actionType
	}
	if isActive != nil {
		rule.IsActive = isActive
	}
	return r.repo.UpdateRule(rule)
}

func (r *RuleUseCase) GetRuleByID(id string) (*entity.Rule, error) {
	rule, err := r.repo.GetRule(id)
	if err != nil {
		return nil, fmt.Errorf("error fetching rule: %w", err)
	}

	if rule == nil {
		return nil, fmt.Errorf("rule not found: id=%s", id)
	}

	return rule, nil
}

func (r *RuleUseCase) GetRulesForResource(id string) ([]entity.Rule, error) {
	rules, err := r.repo.GetRulesForResource(id)
	if err != nil {
		return nil, fmt.Errorf("error fetching rules for resource %s: %w", id, err)
	}

	return rules, nil
}
