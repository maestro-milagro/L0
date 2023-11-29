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
	Create(messageId int, pay message.Payments) (int, error)
	GetById(messageId int) (message.Payments, error)
	Delete(messageId int) error
}

type DeliveriesS interface {
	Create(messageId int, delivery message.Deliveries) (int, error)
	GetById(messageId int) (message.Deliveries, error)
	Delete(messageId int) error
}
type ItemsS interface {
	Create(messageId int, item message.Item) (int, error)
	GetAll(messageId int) ([]message.Item, error)
	GetById(messageId, itemId int) (message.Item, error)
	Delete(messageId, itemId int) error
}

type Service struct {
	MessagesS
	PaymentsS
	DeliveriesS
	ItemsS
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		MessagesS:   NewMessageService(repos.MessagesTB),
		PaymentsS:   NewPaymentsService(repos.PaymentsTB),
		DeliveriesS: NewDeliveriesService(repos.DeliveriesTB),
		ItemsS:      NewItemsService(repos.ItemsTB),
	}
}
