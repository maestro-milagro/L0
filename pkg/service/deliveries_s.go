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
func (s DeliveriesService) Create(delivery message.Deliveries) (int, error) {
	return s.repo.Create(delivery)
}
