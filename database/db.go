package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// DB представляет структуру для работы с базой данных, используя библиотеку sqlx.
type DB struct {
	db *sqlx.DB
}

// NewDB создает новое соединение с базой данных и возвращает указатель на структуру DB.
// Параметры:
// - connectionString string: строка подключения к базе данных.
// Возвращаемые значения:
// - (*DB, error): указатель на структуру DB при успешном подключении или ошибка, если подключение не удалось.
func NewDB(connectionString string) (*DB, error) {
	db, err := sqlx.Connect("pgx", connectionString)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}
	return &DB{db: db}, nil
}

// Close закрывает соединение с базой данных.
// Возвращаемые значения:
// - error: ошибка, если закрытие соединения не удалось.
func (db *DB) Close() error {
	return db.db.Close()
}
