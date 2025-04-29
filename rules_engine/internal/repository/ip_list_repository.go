package repository

import "rules-engine/internal/entity"

type IPListRepository interface {
	GetIPLists() ([]entity.IPList, error)
	CreateIPList(ipList *entity.IPList) (*entity.IPList, error)
	UpdateIPList(ipList *entity.IPList) (*entity.IPList, error)
	GetIPList(id string) (*entity.IPList, error)
	GetIPListsForResource(resourceID string) ([]entity.IPList, error)
	GetIPListsByURL(url, method string) ([]entity.IPList, error)
}
