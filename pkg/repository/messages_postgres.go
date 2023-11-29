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

func (r *MessagesPostgres) Create(message message.Message) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createListQuery := fmt.Sprintf("INSERT INTO %s VALUES (default, $1, $2, $3, $4, $5, $6, $7, $8, $10, $11, $12, $13, $14) RETURNING id", Messages)
	row := tx.QueryRow(createListQuery, message.OrderUid, message.TrackNumber, message.Entry,
		message.Delivery, message.Payment, message.Items, message.Locale,
		message.InternalSignature, message.CustomerId, message.DeliveryService,
		message.Shardkey, message.SmId, message.DateCreated, message.OofShard)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (r *MessagesPostgres) GetAll() ([]message.Message, error) {
	var lists []message.Message

	query := fmt.Sprintf("SELECT * FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1",
		todoListsTable, usersListsTable)
	err := r.db.Select(&lists, query, userId)

	return lists, err
}

func (r *MessagesPostgres) GetById(messageId int) (message.Message, error) {
	var message1 message.Message

	query := fmt.Sprintf(`SELECT tl.id, tl.title, tl.description FROM %s tl
								INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1 AND ul.list_id = $2`,
		todoListsTable, usersListsTable)
	err := r.db.Get(&list, query, userId, listId)

	return message1, err
}

func (r *MessagesPostgres) Delete(messageId int) error {
	query := fmt.Sprintf("DELETE FROM %s tl USING %s ul WHERE tl.id = ul.list_id AND ul.user_id=$1 AND ul.list_id=$2",
		todoListsTable, usersListsTable)
	_, err := r.db.Exec(query, userId, listId)

	return err
}
