package service

import (
	message "awesomeProject"
	"awesomeProject/pkg/repository"
)

type DeliveriesService struct {
	repo repository.DeliveriesTB
}

func NewDeliveriesService(repo repository.DeliveriesTB) *DeliveriesService {
	return &DeliveriesService{repo: repo}
}
func (s DeliveriesService) Create(messageId int, delivery message.Deliveries) (int, error) {
	return s.repo.Create(messageId, delivery)
}
func (s DeliveriesService) GetById(messageId int) (message.Deliveries, error) {
	return s.repo.GetById(messageId)
}
func (s DeliveriesService) Delete(messageId int) error {
	return s.repo.Delete(messageId)
}
