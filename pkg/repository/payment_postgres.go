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
	createListQuery := fmt.Sprintf("INSERT INTO %s VALUES (default, $1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING PaymentId", Payments)
	row := tx.QueryRow(createListQuery, pay.Transaction, pay.RequestId, pay.Currency, pay.Provider, pay.Amount, pay.PaymentDt, pay.Bank, pay.DeliveryCost, pay.GoodsTotal, pay.CustomFee)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}
	return id, tx.Commit()
}
func (r *PaymentsPostgres) GetAll() ([]message.Payments, error) {
	var pay []message.Payments
	query := fmt.Sprintf("SELECT d.paymentid, d.transaction, d.requestId, d.currency, d.provider, d.amount, d.paymentDt, d.bank, d.deliveryCost, d.goodsTotal, d.customFee FROM %s m INNER JOIN %s md on m.MessageId = md.MessageId INNER JOIN %s d on d.PaymentId = md.PaymentId", Messages, MessagesPayments, Payments)
	err := r.db.Select(&pay, query)
	return pay, err
}
func (r *PaymentsPostgres) GetById(messageId int) (message.Payments, error) {
	var pay message.Payments
	queryDel := fmt.Sprintf("SELECT d.paymentid, d.transaction, d.requestId, d.currency, d.provider, d.amount, d.paymentDt, d.bank, d.deliveryCost, d.goodsTotal, d.customFee FROM %s m INNER JOIN %s md on m.MessageId = md.MessageId INNER JOIN %s d on d.PaymentId = md.PaymentId WHERE m.MessageId = $1", Messages, MessagesPayments, Payments)
	err := r.db.Get(&pay, queryDel, messageId)
	return pay, err
}
