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

func (r *ItemsPostgres) Create(item []message.Item) ([]int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return []int{}, err
	}

	ids := make([]int, len(item))
	var id int
	for i, v := range item {
		createListQuery := fmt.Sprintf("INSERT INTO %s VALUES (default, $1, $2, $3, $4, $5, $6, $7, $8, $10, $11) RETURNING id", Items)
		row := tx.QueryRow(createListQuery, v.ChrtId, v.TrackNumber, v.Price,
			v.Rid, v.Name, v.Sale, v.Size, v.TotalPrice,
			v.NmId, v.Brand, v.Status)
		if err := row.Scan(&id); err != nil {
			tx.Rollback()
			return []int{}, err
		}
		ids[i] = id

	}

	return ids, tx.Commit()
}

//func (r *ItemsPostgres) GetAll(messageId int) ([]message.Item, error) {
//	var lists []message.Message
//
//	query := fmt.Sprintf("SELECT * FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1",
//		todoListsTable, usersListsTable)
//	err := r.db.Select(&lists, query, userId)
//
//	return lists, err
//}
//
//func (r *ItemsPostgres) GetById(messageId, itemId int) (message.Item, error) {
//	var message1 message.Message
//
//	query := fmt.Sprintf(`SELECT tl.id, tl.title, tl.description FROM %s tl
//								INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1 AND ul.list_id = $2`,
//		todoListsTable, usersListsTable)
//	err := r.db.Get(&list, query, userId, listId)
//
//	return message1, err
//}
