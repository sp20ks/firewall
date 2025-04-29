package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"rules-engine/internal/entity"

	"rules-engine/internal/repository"

	"github.com/google/uuid"
)

type PostgresResourceRepository struct {
	db *sql.DB
}

func NewPostgresResourceRepository(db *sql.DB) repository.ResourceRepository {
	return &PostgresResourceRepository{db: db}
}

func (r *PostgresResourceRepository) GetActiveResources() ([]entity.Resource, error) {
	rows, err := r.db.Query("SELECT id, name, http_method, url, host, is_active, created_at, creator_id FROM resources WHERE is_active = TRUE")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resources []entity.Resource
	for rows.Next() {
		var res entity.Resource
		if err := rows.Scan(&res.ID, &res.Name, &res.HTTPMethod, &res.URL, &res.Host, &res.IsActive, &res.CreatedAt, &res.CreatorID); err != nil {
			return nil, err
		}
		resources = append(resources, res)
	}
	return resources, nil
}

func (r *PostgresResourceRepository) CreateResource(resource *entity.Resource) (*entity.Resource, error) {
	resource.ID = uuid.New().String()

	var createdResource entity.Resource
	err := r.db.QueryRow(`
		INSERT INTO resources (id, name, http_method, url, host, creator_id, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, name, http_method, url, host, creator_id, is_active, created_at
	`, resource.ID, resource.Name, resource.HTTPMethod, resource.URL, resource.Host, resource.CreatorID, resource.IsActive).Scan(
		&createdResource.ID,
		&createdResource.Name,
		&createdResource.HTTPMethod,
		&createdResource.URL,
		&createdResource.Host,
		&createdResource.CreatorID,
		&createdResource.IsActive,
		&createdResource.CreatedAt,
	)
	return &createdResource, err
}

func (r *PostgresResourceRepository) UpdateResource(resource *entity.Resource) (*entity.Resource, error) {
	var updatedResource entity.Resource

	err := r.db.QueryRow(`
		UPDATE resources
		SET name=$1, http_method=$2, url=$3, host=$4, is_active=$5
		WHERE id=$6
		RETURNING id, name, http_method, url, host, creator_id, is_active, created_at
	`, resource.Name, resource.HTTPMethod, resource.URL, resource.Host, resource.IsActive, resource.ID).Scan(
		&updatedResource.ID,
		&updatedResource.Name,
		&updatedResource.HTTPMethod,
		&updatedResource.URL,
		&updatedResource.Host,
		&updatedResource.CreatorID,
		&updatedResource.IsActive,
		&updatedResource.CreatedAt,
	)

	return &updatedResource, err
}

func (r *PostgresResourceRepository) GetResource(id string) (*entity.Resource, error) {
	query := `SELECT id, name, http_method, url, host, created_at, creator_id, is_active FROM resources WHERE id = $1`

	resource := &entity.Resource{}
	err := r.db.QueryRow(query, id).Scan(
		&resource.ID,
		&resource.Name,
		&resource.HTTPMethod,
		&resource.URL,
		&resource.Host,
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
