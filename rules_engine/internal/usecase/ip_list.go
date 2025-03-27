package usecase

import (
	"rules-engine/internal/entity"
	"rules-engine/internal/repository"
)

type IPListUseCase struct {
	repo repository.IPListRepository
}

func NewIPListUseCase(repo repository.IPListRepository) *IPListUseCase {
	return &IPListUseCase{repo: repo}
}

func (r *IPListUseCase) Get() ([]entity.IPList, error) {
	return r.repo.GetIPLists()
}
