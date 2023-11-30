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

type PaymentsS interface {
	Create(pay message.Payments) (int, error)
}

type DeliveriesS interface {
	Create(delivery message.Deliveries) (int, error)
}
type ItemsS interface {
	Create(item []message.Item) ([]int, error)
}

type Service struct {
	MessagesS
	PaymentsS
	DeliveriesS
	ItemsS
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		MessagesS:   NewMessageService(repos.MessagesTB, repos.DeliveriesTB, repos.PaymentsTB, repos.ItemsTB),
		PaymentsS:   NewPaymentsService(repos.PaymentsTB),
		DeliveriesS: NewDeliveriesService(repos.DeliveriesTB),
		ItemsS:      NewItemsService(repos.ItemsTB),
	}
}
