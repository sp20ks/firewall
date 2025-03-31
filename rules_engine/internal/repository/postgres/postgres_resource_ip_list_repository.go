package postgres

import (
	"database/sql"

	"rules-engine/internal/repository"

	"github.com/google/uuid"
)

type PostgresResourceIPListRepository struct {
	db *sql.DB
}

func NewPostgresResourceIPListRepository(db *sql.DB) repository.ResourceIPListRepository {
	return &PostgresResourceIPListRepository{db: db}
}

func (r *PostgresResourceIPListRepository) AttachIPList(resourceID, ipListID string) error {
	_, err := r.db.Exec("INSERT INTO resource_ip_list (id, resource_id, ip_list_id) VALUES ($1, $2, $3)",
		uuid.New().String(), resourceID, ipListID)
	return err
}

func (r *PostgresResourceIPListRepository) DetachIPList(resourceID, ipListID string) error {
	_, err := r.db.Exec("DELETE FROM resource_ip_list WHERE resource_id = $1 AND ip_list_id = $2",
		resourceID, ipListID)
	return err
}
