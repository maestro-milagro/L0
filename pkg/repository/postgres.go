package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

// Переменные хранящие имена таблиц
const (
	Messages   = "message"
	Deliveries = "delivery"
	Payments   = "payment"
	Items      = "item"
	//Таблицы связывающие message с delivery, payment и item соответсвенно
	MessagesDeliveries = "messagedeliveries"
	MessagesPayments   = "messagepayments"
	MessagesItems      = "messageitem"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

// Функция осуществляющая подключение к бд
func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
