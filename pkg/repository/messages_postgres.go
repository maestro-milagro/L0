package repository

import (
	"awesomeProject"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type MessagesPostgres struct {
	db *sqlx.DB
}

func NewMessagesPostgres(db *sqlx.DB) *MessagesPostgres {
	return &MessagesPostgres{db: db}
}

func (r *MessagesPostgres) Create(message message.Message, delId, payId int, itemsId []int) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createMessageQuery := fmt.Sprintf("INSERT INTO %s VALUES (default, $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING MessageId", Messages)
	row := tx.QueryRow(createMessageQuery, message.OrderUid, message.TrackNumber, message.Entry, message.Locale,
		message.InternalSignature, message.CustomerId, message.DeliveryService,
		message.Shardkey, message.SmId, message.DateCreated, message.OofShard)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}
	createMessagesDeliveriesQuery := fmt.Sprintf("INSERT INTO %s VALUES (default, $1, $2)", MessagesDeliveries)
	_, err = tx.Exec(createMessagesDeliveriesQuery, id, delId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	createMessagesPaymentsQuery := fmt.Sprintf("INSERT INTO %s VALUES (default, $1, $2)", MessagesPayments)
	_, err = tx.Exec(createMessagesPaymentsQuery, id, payId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	for _, v := range itemsId {
		createMessagesItemsQuery := fmt.Sprintf("INSERT INTO %s VALUES (default, $1, $2)", MessagesItems)
		_, err = tx.Exec(createMessagesItemsQuery, id, v)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	return id, tx.Commit()
}

func (r *MessagesPostgres) GetAll() ([]message.Message, error) {
	var lists []message.Message

	query := fmt.Sprintf("SELECT * FROM %s m INNER JOIN %s md on m.MessageId = md.MessageId INNER JOIN %s d on d.DeliveryId = md.DeliveryId "+
		"INNER JOIN %s mp on m.MessageId = mp.MessageId INNER JOIN %s p on p.PaymentId = mp.PaymentId "+
		"INNER JOIN %s mi on m.MessageId = mi.MessageId INNER JOIN %s i on i.ItemId = mi.ItemId",
		Messages, MessagesDeliveries, Deliveries,
		MessagesPayments, Payments, MessagesItems, Items)
	err := r.db.Select(&lists, query)

	return lists, err
}

func (r *MessagesPostgres) GetById(messageId int) (message.Message, error) {
	var message1 message.Message
	var del message.Deliveries
	//var pay message.Payments
	//var items message.Item

	//query := fmt.Sprintf("SELECT * FROM %s m INNER JOIN %s md on m.MessageId = md.MessageId INNER JOIN %s d on d.DeliveryId = md.DeliveryId INNER JOIN %s mp on m.MessageId = mp.MessageId INNER JOIN %s p on p.PaymentId = mp.PaymentId INNER JOIN %s mi on m.MessageId = mi.MessageId INNER JOIN %s i on i.ItemId = mi.ItemId WHERE m.MessageId = $1",
	//	Messages, MessagesDeliveries, Deliveries,
	//	MessagesPayments, Payments, MessagesItems, Items)
	//err := r.db.Get(&message1, query, messageId)
	queryDel := fmt.Sprintf("SELECT d.deliveryid, d.name, d.phone, d.zip, d.city, d.address, d.region, d.email FROM %s m INNER JOIN %s md on m.MessageId = md.MessageId INNER JOIN %s d on d.DeliveryId = md.DeliveryId WHERE m.MessageId = $1", Messages, MessagesDeliveries, Deliveries)
	err := r.db.Get(&del, queryDel, messageId)
	message1.Delivery = del
	return message1, err
}

func (r *MessagesPostgres) Delete(messageId int) error {
	query := fmt.Sprintf("DELETE FROM %s md USING %s d WHERE md.MessageId = d.MessageId AND d.MessageId=$1",
		MessagesDeliveries, Deliveries)
	_, err := r.db.Exec(query, messageId)
	query = fmt.Sprintf("DELETE FROM %s m USING %s mp WHERE m.MessageId = mp.MessageId AND mp.MessageId=$1",
		MessagesPayments, Payments)
	_, err = r.db.Exec(query, messageId)
	query = fmt.Sprintf("DELETE FROM %s m USING %s mi WHERE m.MessageId = mi.MessageId AND mi.MessageId=$1",
		MessagesItems, Items)
	_, err = r.db.Exec(query, messageId)
	return err
}
