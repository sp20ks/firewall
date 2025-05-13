package usecase

import (
	"fmt"
	"time"

	"rules-engine/internal/entity"
	"rules-engine/internal/logger"
	"rules-engine/internal/repository"

	"go.uber.org/zap"
)

type ResourceUseCase struct {
	resourceRepo       repository.ResourceRepository
	iPListUseCase      *IPListUseCase
	ruleUseCase        *RuleUseCase
	resourceIPListRepo repository.ResourceIPListRepository
	resourceRuleRepo   repository.ResourceRuleRepository
}

func NewResourceUseCase(
	resourceRepo repository.ResourceRepository,
	iPListUseCase *IPListUseCase,
	ruleUseCase *RuleUseCase,
	resourceIPListRepo repository.ResourceIPListRepository,
	resourceRuleRepo repository.ResourceRuleRepository,
) *ResourceUseCase {
	return &ResourceUseCase{
		resourceRepo:       resourceRepo,
		iPListUseCase:      iPListUseCase,
		ruleUseCase:        ruleUseCase,
		resourceIPListRepo: resourceIPListRepo,
		resourceRuleRepo:   resourceRuleRepo,
	}
}

func (r *ResourceUseCase) GetResourceByID(id string) (*entity.Resource, error) {
	resource, err := r.resourceRepo.GetResource(id)
	if err != nil {
		return nil, fmt.Errorf("error fetching resource: %w", err)
	}
	if resource == nil {
		return nil, fmt.Errorf("resource not found: id=%s", id)
	}

	ipLists, err := r.iPListUseCase.GetIPListsForResource(resource.ID)
	if err != nil {
		logger.Logger().Info(
			"error fetching IP lists for resource",
			zap.String("resource_id", resource.ID),
			zap.Error(err),
		)
	} else {
		resource.IPLists = ipLists
	}
	rules, err := r.ruleUseCase.GetRulesForResource(resource.ID)
	if err != nil {
		logger.Logger().Info(
			"error fetching rules for resource",
			zap.String("resource_id", resource.ID),
			zap.Error(err),
		)
	} else {
		resource.Rules = rules
	}

	return resource, nil
}

func (r *ResourceUseCase) Create(name, method, url, host, creatorID string, isActive *bool) (*entity.Resource, error) {
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

func (r *ResourceUseCase) Update(id, name, method, url, host string, isActive *bool) (*entity.Resource, error) {
	resource, err := r.GetResourceByID(id)
	if err != nil {
		return nil, err
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

// TODO: добавить флаг, при котором отправляются вместе с ресурсами правила и списки
func (r *ResourceUseCase) Get() ([]entity.Resource, error) {
	resources, err := r.resourceRepo.GetResources()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch resources: %w", err)
	}

	for i, res := range resources {
		ipLists, err := r.iPListUseCase.GetIPListsForResource(res.ID)
		if err != nil {
			logger.Logger().Info(
				"error fetching IP lists for resource",
				zap.String("resource_id", res.ID),
				zap.Error(err),
			)
		} else {
			resources[i].IPLists = ipLists
		}
		rules, err := r.ruleUseCase.GetRulesForResource(res.ID)
		if err != nil {
			logger.Logger().Info(
				"error fetching rules for resource",
				zap.String("resource_id", res.ID),
				zap.Error(err),
			)
		} else {
			resources[i].Rules = rules
		}
	}

	return resources, nil
}

func (r *ResourceUseCase) AttachIPList(resourceID, ipListID string) error {
	if _, err := r.GetResourceByID(resourceID); err != nil {
		return err
	}
	if _, err := r.iPListUseCase.getIPListByID(ipListID); err != nil {
		return err
	}
	return r.resourceIPListRepo.AttachIPList(resourceID, ipListID)
}

func (r *ResourceUseCase) DetachIPList(resourceID, ipListID string) error {
	if _, err := r.GetResourceByID(resourceID); err != nil {
		return err
	}
	if _, err := r.iPListUseCase.getIPListByID(ipListID); err != nil {
		return err
	}
	return r.resourceIPListRepo.DetachIPList(resourceID, ipListID)
}

func (r *ResourceUseCase) AttachRule(resourceID, ruleID string) error {
	if _, err := r.GetResourceByID(resourceID); err != nil {
		return err
	}
	if _, err := r.ruleUseCase.GetRuleByID(ruleID); err != nil {
		return err
	}
	return r.resourceRuleRepo.AttachRule(resourceID, ruleID)
}

func (r *ResourceUseCase) DetachRule(resourceID, ruleID string) error {
	if _, err := r.GetResourceByID(resourceID); err != nil {
		return err
	}
	if _, err := r.ruleUseCase.GetRuleByID(ruleID); err != nil {
		return err
	}
	return r.resourceRuleRepo.DetachRule(resourceID, ruleID)
}
