package repository

import (
	"awesomeProject"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type PaymentsPostgres struct {
	db *sqlx.DB
}

func NewPaymentsPostgres(db *sqlx.DB) *PaymentsPostgres {
	return &PaymentsPostgres{db: db}
}

func (r *PaymentsPostgres) Create(pay message.Payments) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createListQuery := fmt.Sprintf("INSERT INTO %s VALUES (default, $1, $2, $3, $4, $5, $6, $7, $8, $10) RETURNING id", Payments)
	row := tx.QueryRow(createListQuery, pay.Transaction, pay.RequestId, pay.Currency, pay.Provider, pay.Amount, pay.PaymentDt, pay.Bank, pay.DeliveryCost, pay.GoodsTotal, pay.CustomFee)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}
	return id, tx.Commit()
}
