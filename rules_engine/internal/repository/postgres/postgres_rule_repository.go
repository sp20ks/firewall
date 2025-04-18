package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"rules-engine/internal/entity"

	"rules-engine/internal/repository"

	"github.com/google/uuid"
)

type PostgresRuleRepository struct {
	db *sql.DB
}

func NewPostgresRuleRepository(db *sql.DB) repository.RuleRepository {
	return &PostgresRuleRepository{db: db}
}

func (r *PostgresRuleRepository) GetActiveRules() ([]entity.Rule, error) {
	rows, err := r.db.Query("SELECT id, name, attack_type, action_type, is_active, created_at, creator_id FROM rules WHERE is_active = TRUE")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []entity.Rule
	for rows.Next() {
		var res entity.Rule
		if err := rows.Scan(&res.ID, &res.Name, &res.AttackType, &res.ActionType, &res.IsActive, &res.CreatedAt, &res.CreatorID); err != nil {
			return nil, err
		}
		rules = append(rules, res)
	}
	return rules, nil
}

func (r *PostgresRuleRepository) CreateRule(rule *entity.Rule) error {
	rule.ID = uuid.New().String()
	_, err := r.db.Exec("INSERT INTO rules (id, name, attack_type, action_type, creator_id, is_active) VALUES ($1, $2, $3, $4, $5, $6)",
		rule.ID, rule.Name, rule.AttackType, rule.ActionType, rule.CreatorID, rule.IsActive)
	return err
}

func (r *PostgresRuleRepository) UpdateRule(rule *entity.Rule) error {
	_, err := r.db.Exec("UPDATE rules SET name=$1, attack_type=$2, action_type=$3, is_active=$4 WHERE id=$5",
		rule.Name, rule.AttackType, rule.ActionType, rule.IsActive, rule.ID)
	return err
}

func (r *PostgresRuleRepository) GetRule(id string) (*entity.Rule, error) {
	query := `SELECT id, name, attack_type, action_type, created_at, creator_id, is_active FROM rules WHERE id = $1`

	rule := &entity.Rule{}
	err := r.db.QueryRow(query, id).Scan(
		&rule.ID,
		&rule.Name,
		&rule.AttackType,
		&rule.ActionType,
		&rule.CreatedAt,
		&rule.CreatorID,
		&rule.IsActive,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get rule: %w", err)
	}
	return rule, nil
}

func (r *PostgresRuleRepository) GetRulesForResource(resourceID string) ([]entity.Rule, error) {
	query := `
		SELECT t1.id, t1.name, t1.attack_type, t1.action_type, t1.is_active, t1.creator_id, t1.created_at
		FROM rules AS t1
		INNER JOIN resource_rule AS t2
		ON t1.id = t2.rule_id
		WHERE t2.resource_id = $1
	`

	rows, err := r.db.Query(query, resourceID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get rules: %w", err)
	}
	defer rows.Close()

	var rules []entity.Rule
	for rows.Next() {
		var res entity.Rule
		if err := rows.Scan(&res.ID, &res.Name, &res.AttackType, &res.ActionType, &res.IsActive, &res.CreatorID, &res.CreatedAt); err != nil {
			return nil, err
		}
		rules = append(rules, res)
	}

	return rules, nil
}

func (r *PostgresRuleRepository) GetRulesByURL(url, method string) ([]entity.Rule, error) {
	query := `
		SELECT t1.id, t1.name, t1.attack_type, t1.action_type, t1.is_active, t1.creator_id, t1.created_at
		FROM rules AS t1
		INNER JOIN resource_rule AS t2
		ON t1.id = t2.rule_id
		INNER JOIN resources AS t3
		ON t2.resource_id = t3.id
		WHERE t3.url = $1 AND t3.http_method = $2
	`

	rows, err := r.db.Query(query, url, method)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get rules: %w", err)
	}
	defer rows.Close()

	var rules []entity.Rule
	for rows.Next() {
		var res entity.Rule
		if err := rows.Scan(&res.ID, &res.Name, &res.AttackType, &res.ActionType, &res.IsActive, &res.CreatorID, &res.CreatedAt); err != nil {
			return nil, err
		}
		rules = append(rules, res)
	}

	return rules, nil
}
