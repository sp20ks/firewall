package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"rules-engine/internal/entity"

	"rules-engine/internal/repository"
)

type PostgresResourceRepository struct {
	db *sql.DB
}

func NewPostgresResourceRepository(db *sql.DB) repository.ResourceRepository {
	return &PostgresResourceRepository{db: db}
}

func (r *PostgresResourceRepository) GetActiveResources() ([]entity.Resource, error) {
	rows, err := r.db.Query("SELECT id, name, http_method, url, created_at, creator_id FROM resources WHERE is_active = TRUE")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resources []entity.Resource
	for rows.Next() {
		var res entity.Resource
		if err := rows.Scan(&res.ID, &res.Name, &res.HTTPMethod, &res.URL, &res.CreatedAt, &res.CreatorID); err != nil {
			return nil, err
		}
		resources = append(resources, res)
	}
	return resources, nil
}

func (r *PostgresResourceRepository) CreateResource(resource *entity.Resource) error {
	_, err := r.db.Exec("INSERT INTO resources (id, name, http_method, url, creator_id, is_active) VALUES ($1, $2, $3, $4, $5, $6)",
		resource.ID, resource.Name, resource.HTTPMethod, resource.URL, resource.CreatorID, resource.IsActive)
	return err
}

func (r *PostgresResourceRepository) UpdateResource(resource *entity.Resource) error {
	_, err := r.db.Exec("UPDATE resources SET name=$1, http_method=$2, url=$3, is_active=$4 WHERE id=$5",
		resource.Name, resource.HTTPMethod, resource.URL, resource.IsActive, resource.ID)
	return err
}

func (r *PostgresResourceRepository) GetResource(id string) (*entity.Resource, error) {
	query := `SELECT id, name, http_method, url, created_at, creator_id, is_active FROM resources WHERE id = $1`

	resource := &entity.Resource{}
	err := r.db.QueryRow(query, id).Scan(
		&resource.ID,
		&resource.Name,
		&resource.HTTPMethod,
		&resource.URL,
		&resource.CreatedAt,
		&resource.CreatorID,
		&resource.IsActive,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get resource: %w", err)
	}
	return resource, nil
}
