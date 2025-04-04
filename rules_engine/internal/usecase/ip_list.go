package usecase

import (
	"fmt"
	"net"
	"rules-engine/internal/entity"
	"rules-engine/internal/repository"
)

type IPListUseCase struct {
	repo repository.IPListRepository
}

func NewIPListUseCase(repo repository.IPListRepository) *IPListUseCase {
	return &IPListUseCase{repo: repo}
}

func (i *IPListUseCase) Get() ([]entity.IPList, error) {
	return i.repo.GetIPLists()
}

func (i *IPListUseCase) Create(cidrStr, listType, CreatorID string) error {
	list := &entity.IPList{
		ListType:  listType,
		CreatorID: CreatorID,
	}

	ip, ipNet, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return fmt.Errorf("failed to parse CIDR: %w", err)
	}
	ipNet.IP = ip
	list.IP = *ipNet

	return i.repo.CreateIPList(list)
}

func (i *IPListUseCase) Update(id, cidrStr, listType string) error {
	list, err := i.repo.GetIPList(id)
	if err != nil {
		return fmt.Errorf("error fetching ip list: %w", err)
	}

	if list == nil {
		return fmt.Errorf("ip list with id=%s not found", id)
	}

	if cidrStr != "" {
		ip, ipNet, err := net.ParseCIDR(cidrStr)
		if err != nil {
			return fmt.Errorf("failed to parse CIDR: %w", err)
		}
		ipNet.IP = ip
		list.IP = *ipNet
	}
	if listType != "" {
		list.ListType = listType
	}
	return i.repo.UpdateIPList(list)
}
