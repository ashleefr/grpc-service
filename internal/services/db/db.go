package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // импорт драйвера PostgreSQL
)

type DB struct {
	*sql.DB
}

func NewDB(dataSourceName string) (*DB, error) {
	// Создание подключения к базе данных
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Пингуем базу данных, чтобы убедиться, что подключение установлено
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return &DB{DB: db}, nil
}
