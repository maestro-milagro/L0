package repository

import (
	"awesomeProject"
	"github.com/jmoiron/sqlx"
)

// Т.к. в нашей структуре данные о заказе разделены на 4 таблицы
// Создадим интерфейс и определим методы для работы с каждой из них
type MessagesTB interface {
	Create(message message.Message, delId, payId int, itemsId []int) (int, error)
	GetAll() ([]message.Message, error)
	GetById(messageId int) (message.Message, error)
	Delete(messageId int) error
}

type PaymentsTB interface {
	Create(pay message.Payments) (int, error)
	GetById(messageId int) (message.Payments, error)
	GetAll() ([]message.Payments, error)
}

type DeliveriesTB interface {
	Create(delivery message.Deliveries) (int, error)
	GetById(messageId int) (message.Deliveries, error)
	GetAll() ([]message.Deliveries, error)
}
type ItemsTB interface {
	Create(item []message.Item) ([]int, error)
	GetById(messageId int) ([]message.Item, error)
	GetAll() ([]message.Item, error)
}

// Создадим структуру содержащую в себе реализации всех интерфейсов определенных выше
// чтобы получить возможность работать с каждой из таблиц из одной структуры и
// заменить или добавить новые реализации если у нас будет такая потребность
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
