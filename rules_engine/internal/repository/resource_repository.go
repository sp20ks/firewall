package repository

import "rules-engine/internal/entity"

type ResourceRepository interface {
	GetActiveResources() ([]entity.Resource, error)
	CreateResource(resource *entity.Resource) (*entity.Resource, error)
	UpdateResource(resource *entity.Resource) (*entity.Resource, error)
	GetResource(id string) (*entity.Resource, error)
}
