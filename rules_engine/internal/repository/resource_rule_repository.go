package repository

type ResourceRuleRepository interface {
	AttachRule(resourceID, ruleID string) error
	DetachRule(resourceID, ruleID string) error
}
