package repository

import (
	"awesomeProject"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type DeliveriesPostgres struct {
	db *sqlx.DB
}

func NewDeliveriesPostgres(db *sqlx.DB) *DeliveriesPostgres {
	return &DeliveriesPostgres{db: db}
}

func (r *DeliveriesPostgres) Create(delivery message.Deliveries) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createListQuery := fmt.Sprintf("INSERT INTO %s VALUES (default, $1, $2, $3, $4, $5, $6, $7) RETURNING DeliveryId", Deliveries)
	row := tx.QueryRow(createListQuery, delivery.Name, delivery.Phone, delivery.Zip,
		delivery.City, delivery.Address, delivery.Region, delivery.Email)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}
func (r *DeliveriesPostgres) GetAll() ([]message.Deliveries, error) {
	var del []message.Deliveries
	queryDel := fmt.Sprintf("SELECT d.deliveryid, d.name, d.phone, d.zip, d.city, d.address, d.region, d.email FROM %s m INNER JOIN %s md on m.MessageId = md.MessageId INNER JOIN %s d on d.DeliveryId = md.DeliveryId", Messages, MessagesDeliveries, Deliveries)
	err := r.db.Select(&del, queryDel)
	return del, err
}
func (r *DeliveriesPostgres) GetById(messageId int) (message.Deliveries, error) {
	var del message.Deliveries
	queryDel := fmt.Sprintf("SELECT d.deliveryid, d.name, d.phone, d.zip, d.city, d.address, d.region, d.email FROM %s m INNER JOIN %s md on m.MessageId = md.MessageId INNER JOIN %s d on d.DeliveryId = md.DeliveryId WHERE m.MessageId = $1", Messages, MessagesDeliveries, Deliveries)
	err := r.db.Get(&del, queryDel, messageId)
	return del, err
}
