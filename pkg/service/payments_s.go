package service

import (
	message "awesomeProject"
	"awesomeProject/pkg/repository"
)

type PaymentsService struct {
	repo repository.PaymentsTB
}

func NewPaymentsService(repo repository.PaymentsTB) *PaymentsService {
	return &PaymentsService{repo: repo}
}
func (s *PaymentsService) Create(pay message.Payments) (int, error) {
	return s.repo.Create(pay)
}
