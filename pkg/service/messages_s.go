package service

import (
	message "awesomeProject"
	"awesomeProject/pkg/repository"
)

type MessageService struct {
	repo  repository.MessagesTB
	repoD repository.DeliveriesTB
	repoP repository.PaymentsTB
	repoI repository.ItemsTB
}

func NewMessageService(repo repository.MessagesTB, repoD repository.DeliveriesTB, repoP repository.PaymentsTB, repoI repository.ItemsTB) *MessageService {
	return &MessageService{
		repo:  repo,
		repoD: repoD,
		repoP: repoP,
		repoI: repoI,
	}
}
func (s *MessageService) Create(message message.Message) (int, error) {
	delId, err := s.repoD.Create(message.Delivery)
	if err != nil {
		return 0, err
	}
	payId, err := s.repoP.Create(message.Payment)
	if err != nil {
		return 0, err
	}
	itemsId, err := s.repoI.Create(message.Items)
	if err != nil {
		return 0, err
	}
	return s.repo.Create(message, delId, payId, itemsId)
}
func (s *MessageService) GetAll() ([]message.Message, error) {
	return s.repo.GetAll()
}
func (s *MessageService) GetById(messageId int) (message.Message, error) {
	return s.repo.GetById(messageId)
}
func (s *MessageService) Delete(messageId int) error {
	return s.repo.Delete(messageId)
}
