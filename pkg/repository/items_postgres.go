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
		createListQuery := fmt.Sprintf("INSERT INTO %s VALUES (default, $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING ItemId", Items)
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

func (r *ItemsPostgres) GetAll() ([]message.Item, error) {
	var it []message.Item
	query := fmt.Sprintf("SELECT d.itemid, d.chrtId, d.trackNumber, d.price, d.rid, d.name, d.sale, d.size, d.totalPrice, d.nmid, d.nrand, d.status FROM %s m INNER JOIN %s md on m.MessageId = md.MessageId INNER JOIN %s d on d.ItemId = md.ItemId", Messages, MessagesItems, Items)
	err := r.db.Select(&it, query)
	return it, err
}
func (r *ItemsPostgres) GetById(messageId int) ([]message.Item, error) {
	var it []message.Item
	query := fmt.Sprintf("SELECT d.itemid, d.chrtId, d.trackNumber, d.price, d.rid, d.name, d.sale, d.size, d.totalPrice, d.nmid, d.brand, d.status FROM %s m INNER JOIN %s md on m.MessageId = md.MessageId INNER JOIN %s d on d.ItemId = md.ItemId WHERE m.MessageId = $1", Messages, MessagesItems, Items)
	err := r.db.Select(&it, query, messageId)
	return it, err
}
