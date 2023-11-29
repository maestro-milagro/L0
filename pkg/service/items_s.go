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
func (s ItemsService) Create(messageId int, item message.Item) (int, error) {
	return s.repo.Create(messageId, item)
}
func (s ItemsService) GetAll(messageId int) ([]message.Item, error) {
	return s.repo.GetAll(messageId)
}
func (s ItemsService) GetById(messageId, itemId int) (message.Item, error) {
	return s.repo.GetById(messageId, itemId)
}
