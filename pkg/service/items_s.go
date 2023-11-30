package service

import (
	message "awesomeProject"
	"awesomeProject/pkg/repository"
)

type ItemsService struct {
	repo repository.ItemsTB
}

func NewItemsService(repo repository.ItemsTB) *ItemsService {
	return &ItemsService{repo: repo}
}
func (s ItemsService) Create(item []message.Item) ([]int, error) {
	return s.repo.Create(item)
}
