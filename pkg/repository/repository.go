package repository

import (
	"awesomeProject"
	"github.com/jmoiron/sqlx"
)

type MessagesTB interface {
	Create(message message.Message, delId, payId, itemsId int) (int, error)
	GetAll() ([]message.Message, error)
	GetById(messageId int) (message.Message, error)
	Delete(messageId int) error
}

type PaymentsTB interface {
	Create(pay message.Payments) (int, error)
}

type DeliveriesTB interface {
	Create(delivery message.Deliveries) (int, error)
}
type ItemsTB interface {
	Create(item message.Item) (int, error)
	GetAll(messageId, itemsId int) ([]message.Item, error)
	GetById(messageId, itemId int) (message.Item, error)
}

type Repository struct {
	MessagesTB
	PaymentsTB
	DeliveriesTB
	ItemsTB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		MessagesTB:   NewMessagesPostgres(db),
		PaymentsTB:   NewPaymentsPostgres(db),
		DeliveriesTB: NewDeliveriesPostgres(db),
		ItemsTB:      NewItemsPostgres(db),
	}
}
