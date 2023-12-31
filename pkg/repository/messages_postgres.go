package repository

import (
	"awesomeProject"
	"fmt"
	"github.com/jmoiron/sqlx"
)

// Создаем структуру для реализации интерфейса для работы с таблицей message
type MessagesPostgres struct {
	db *sqlx.DB
}

func NewMessagesPostgres(db *sqlx.DB) *MessagesPostgres {
	return &MessagesPostgres{db: db}
}

// Функция для сохранения переданного объекта в бд
func (r *MessagesPostgres) Create(message message.Message, delId, payId int, itemsId []int) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	// Создаем запрос и сохраняем данные переданные нам для таблицы message в саму таблицу,
	//а так же сохраняем id в таблицы связывающие её с остальными
	createMessageQuery := fmt.Sprintf("INSERT INTO %s VALUES (default, $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING MessageId", Messages)
	row := tx.QueryRow(createMessageQuery, message.OrderUid, message.TrackNumber, message.Entry, message.Locale,
		message.InternalSignature, message.CustomerId, message.DeliveryService,
		message.Shardkey, message.SmId, message.OofShard, message.DateCreated)
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

// Функция для получения всех объектов из бд
func (r *MessagesPostgres) GetAll() ([]message.Message, error) {
	var lists []message.Message

	query := fmt.Sprintf("SELECT * FROM %s", Messages)
	err := r.db.Select(&lists, query)

	return lists, err
}

// Функция для получения объекта по id
func (r *MessagesPostgres) GetById(messageId int) (message.Message, error) {
	var message1 message.Message
	query := fmt.Sprintf("SELECT * FROM %s WHERE MessageId = $1", Messages)
	err := r.db.Get(&message1, query, messageId)
	return message1, err
}

// Функция для удаления объекта по id
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
