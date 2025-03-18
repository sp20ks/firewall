package usecase

import (
	"fmt"
	"time"

	"rules-engine/internal/entity"
	"rules-engine/internal/repository"
)

type ResorceUseCase struct {
	repo repository.ResourceRepository
}

func NewResorceUseCase(repo repository.ResourceRepository) *ResorceUseCase {
	return &ResorceUseCase{repo: repo}
}

func (r *ResorceUseCase) Create(name, method, url, host, creator_id string, is_active *bool) error {
	resource := &entity.Resource{
		Name:       name,
		HTTPMethod: method,
		URL:        url,
		Host:       host,
		CreatorID:  creator_id,
		IsActive:   is_active,
		CreatedAt:  time.Now(),
	}
	return r.repo.CreateResource(resource)
}

func (r *ResorceUseCase) Update(id, name, method, url, host string, is_active *bool) error {
	resource, err := r.repo.GetResource(id)
	if err != nil {
		return fmt.Errorf("error fetching resource: %w", err)
	}

	if name != "" {
		resource.Name = name
	}
	if method != "" {
		resource.HTTPMethod = method
	}
	if url != "" {
		resource.URL = url
	}
	if host != "" {
		resource.Host = host
	}
	if is_active != nil {
		resource.IsActive = is_active
	}
	return r.repo.UpdateResource(resource)
}

func (r *ResorceUseCase) Get() ([]entity.Resource, error) {
	return r.repo.GetActiveResources()
}
