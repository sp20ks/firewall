package usecase

import (
	"fmt"
	"log"
	"time"

	"rules-engine/internal/entity"
	"rules-engine/internal/repository"
)

type ResourceUseCase struct {
	resourceRepo repository.ResourceRepository
	ipListRepo   repository.IPListRepository
}

func NewResourceUseCase(resourceRepo repository.ResourceRepository, ipListRepo repository.IPListRepository) *ResourceUseCase {
	return &ResourceUseCase{
		resourceRepo: resourceRepo,
		ipListRepo:   ipListRepo,
	}
}

func (r *ResourceUseCase) Create(name, method, url, host, creator_id string, is_active *bool) error {
	resource := &entity.Resource{
		Name:       name,
		HTTPMethod: method,
		URL:        url,
		Host:       host,
		CreatorID:  creator_id,
		IsActive:   is_active,
		CreatedAt:  time.Now(),
	}
	return r.resourceRepo.CreateResource(resource)
}

func (r *ResourceUseCase) Update(id, name, method, url, host string, is_active *bool) error {
	resource, err := r.resourceRepo.GetResource(id)
	if err != nil {
		return fmt.Errorf("error fetching resource: %w", err)
	}

	if resource == nil {
		return fmt.Errorf("resource with id=%s not found", id)
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
	return r.resourceRepo.UpdateResource(resource)
}

func (r *ResourceUseCase) Get() ([]entity.Resource, error) {
	resources, err := r.resourceRepo.GetActiveResources()
	if err != nil {
		return nil, err
	}

	for i, res := range resources {
		ipLists, err := r.ipListRepo.GetIPListsForResource(res.ID)
		if err != nil {
			log.Println(err)
			return nil, fmt.Errorf("failed to fetch IP lists for resource %s: %w", res.ID, err)
		}
		resources[i].IPLists = ipLists
	}

	return resources, nil
}
