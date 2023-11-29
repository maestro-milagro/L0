package repository

import (
	"awesomeProject"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type ItemsPostgres struct {
	db *sqlx.DB
}

func NewItemsPostgres(db *sqlx.DB) *ItemsPostgres {
	return &ItemsPostgres{db: db}
}

func (r *ItemsPostgres) Create(messageId int, item message.Item) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createListQuery := fmt.Sprintf("INSERT INTO %s VALUES (default, $1, $2, $3, $4, $5, $6, $7, $8, $10, $11) RETURNING id", Items)
	row := tx.QueryRow(createListQuery, item.ChrtId, item.TrackNumber, item.Price,
		item.Rid, item.Name, item.Sale, item.Size, item.TotalPrice,
		item.NmId, item.Brand, item.Status)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (r *ItemsPostgres) GetAll(messageId int) ([]message.Item, error) {
	var lists []message.Message

	query := fmt.Sprintf("SELECT * FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1",
		todoListsTable, usersListsTable)
	err := r.db.Select(&lists, query, userId)

	return lists, err
}

func (r *ItemsPostgres) GetById(messageId, itemId int) (message.Item, error) {
	var message1 message.Message

	query := fmt.Sprintf(`SELECT tl.id, tl.title, tl.description FROM %s tl
								INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1 AND ul.list_id = $2`,
		todoListsTable, usersListsTable)
	err := r.db.Get(&list, query, userId, listId)

	return message1, err
}

func (r *ItemsPostgres) Delete(messageId, itemId int) error {
	query := fmt.Sprintf("DELETE FROM %s tl USING %s ul WHERE tl.id = ul.list_id AND ul.user_id=$1 AND ul.list_id=$2",
		todoListsTable, usersListsTable)
	_, err := r.db.Exec(query, userId, listId)

	return err
}
