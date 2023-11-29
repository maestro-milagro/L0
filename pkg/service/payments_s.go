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
func (s *PaymentsService) Create(messageId int, pay message.Payments) (int, error) {
	return s.repo.Create(messageId, pay)
}
func (s *PaymentsService) GetById(messageId int) (message.Payments, error) {
	return s.repo.GetById(messageId)
}
func (s *PaymentsService) Delete(messageId int) error {
	return s.repo.Delete(messageId)
}
