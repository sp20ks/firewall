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
		rule.ID, rule.Name, rule.AttackType, rule.ActionType, rule.IsActive, rule.CreatorID)
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
