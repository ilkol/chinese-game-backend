package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewDBConnection(user, pass, host, port, db string) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, db)

	connection, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("Ошибка подключения к БД: %w", err)
	}

	err = connection.Ping()
	if err != nil {
		return nil, fmt.Errorf("БД недоступна: %w", err)
	}

	return connection, nil
}
