package service

import (
	"awesomeProject"
	"awesomeProject/pkg/repository"
)

type MessagesS interface {
	Create(message message.Message) (int, error)
	GetAll() ([]message.Message, error)
	GetById(messageId int) (message.Message, error)
	Delete(messageId int) error
}

type Service struct {
	MessagesS
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		MessagesS: NewMessageService(repos.MessagesTB, repos.DeliveriesTB, repos.PaymentsTB, repos.ItemsTB),
	}
}
