package repository

import "rules-engine/internal/entity"

type IPListRepository interface {
	GetIPLists() ([]entity.IPList, error)
	CreateIPList(resource *entity.IPList) error
	UpdateIPList(resource *entity.IPList) error
	GetIPList(id string) (*entity.IPList, error)
}
