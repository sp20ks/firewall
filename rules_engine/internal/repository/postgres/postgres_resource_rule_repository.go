package postgres

import (
	"database/sql"

	"rules-engine/internal/repository"

	"github.com/google/uuid"
)

type PostgresResourceRuleRepository struct {
	db *sql.DB
}

func NewPostgresResourceRuleRepository(db *sql.DB) repository.ResourceRuleRepository {
	return &PostgresResourceRuleRepository{db: db}
}

func (r *PostgresResourceRuleRepository) AttachRule(resourceID, ruleID string) error {
	_, err := r.db.Exec("INSERT INTO resource_rule (id, resource_id, rule_id) VALUES ($1, $2, $3)",
		uuid.New().String(), resourceID, ruleID)
	return err
}

func (r *PostgresResourceRuleRepository) DetachRule(resourceID, ruleID string) error {
	_, err := r.db.Exec("DELETE FROM resource_rule WHERE resource_id = $1 AND rule_id = $2",
		resourceID, ruleID)
	return err
}
