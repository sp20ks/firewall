package repository

import "rules-engine/internal/entity"

type IPListRepository interface {
	GetIPLists() ([]entity.IPList, error)
	CreateIPList(ipList *entity.IPList) error
	UpdateIPList(ipList *entity.IPList) error
	GetIPList(id string) (*entity.IPList, error)
}
