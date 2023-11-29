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

func (r *DeliveriesPostgres) Create(messageId int, delivery message.Deliveries) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createListQuery := fmt.Sprintf("INSERT INTO %s VALUES (default, $1, $2, $3, $4, $5, $6, $7, $8, $10, $11, $12, $13, $14) RETURNING id", Deliveries)
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

func (r *DeliveriesPostgres) GetById(messageId int) (message.Deliveries, error) {
	var delivery message.Message

	query := fmt.Sprintf(`SELECT * FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1 AND ul.list_id = $2`,
		todoListsTable, usersListsTable)
	err := r.db.Get(&list, query, userId, listId)

	return delivery, err
}

func (r *DeliveriesPostgres) Delete(messageId int) error {
	query := fmt.Sprintf("DELETE FROM %s tl USING %s ul WHERE tl.id = ul.list_id AND ul.user_id=$1 AND ul.list_id=$2",
		todoListsTable, usersListsTable)
	_, err := r.db.Exec(query, userId, listId)

	return err
}
