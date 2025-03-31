package usecase

import (
	"errors"
	"fmt"
	"log"
	"time"

	"rules-engine/internal/entity"
	"rules-engine/internal/repository"
)

// Определяем ошибки для читаемости
var (
	ErrResourceNotFound = errors.New("resource not found")
	ErrIPListNotFound   = errors.New("IP list not found")
)

type ResourceUseCase struct {
	resourceRepo       repository.ResourceRepository
	ipListRepo         repository.IPListRepository
	resourceIPListRepo repository.ResourceIPListRepository
}

func NewResourceUseCase(
	resourceRepo repository.ResourceRepository,
	ipListRepo repository.IPListRepository,
	resourceIPListRepo repository.ResourceIPListRepository,
) *ResourceUseCase {
	return &ResourceUseCase{
		resourceRepo:       resourceRepo,
		ipListRepo:         ipListRepo,
		resourceIPListRepo: resourceIPListRepo,
	}
}

func (r *ResourceUseCase) getResourceByID(id string) (*entity.Resource, error) {
	resource, err := r.resourceRepo.GetResource(id)
	if err != nil {
		return nil, fmt.Errorf("error fetching resource: %w", err)
	}
	if resource == nil {
		return nil, fmt.Errorf("%w: id=%s", ErrResourceNotFound, id)
	}
	return resource, nil
}

func (r *ResourceUseCase) getIPListByID(id string) (*entity.IPList, error) {
	list, err := r.ipListRepo.GetIPList(id)
	if err != nil {
		return nil, fmt.Errorf("error fetching IP list: %w", err)
	}
	if list == nil {
		return nil, fmt.Errorf("%w: id=%s", ErrIPListNotFound, id)
	}
	return list, nil
}

func (r *ResourceUseCase) Create(name, method, url, host, creatorID string, isActive *bool) error {
	resource := &entity.Resource{
		Name:       name,
		HTTPMethod: method,
		URL:        url,
		Host:       host,
		CreatorID:  creatorID,
		IsActive:   isActive,
		CreatedAt:  time.Now(),
	}
	return r.resourceRepo.CreateResource(resource)
}

func (r *ResourceUseCase) Update(id, name, method, url, host string, isActive *bool) error {
	resource, err := r.getResourceByID(id)
	if err != nil {
		return err
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
	if isActive != nil {
		resource.IsActive = isActive
	}

	return r.resourceRepo.UpdateResource(resource)
}

func (r *ResourceUseCase) Get() ([]entity.Resource, error) {
	resources, err := r.resourceRepo.GetActiveResources()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch resources: %w", err)
	}

	for i, res := range resources {
		ipLists, err := r.ipListRepo.GetIPListsForResource(res.ID)
		if err != nil {
			log.Printf("error fetching IP lists for resource %s: %v", res.ID, err)
			continue
		}
		resources[i].IPLists = ipLists
	}

	return resources, nil
}

func (r *ResourceUseCase) AttachIPList(resourceID, ipListID string) error {
	if _, err := r.getResourceByID(resourceID); err != nil {
		return err
	}
	if _, err := r.getIPListByID(ipListID); err != nil {
		return err
	}
	return r.resourceIPListRepo.AttachIPList(resourceID, ipListID)
}

func (r *ResourceUseCase) DetachIPList(resourceID, ipListID string) error {
	if _, err := r.getResourceByID(resourceID); err != nil {
		return err
	}
	if _, err := r.getIPListByID(ipListID); err != nil {
		return err
	}
	return r.resourceIPListRepo.DetachIPList(resourceID, ipListID)
}
