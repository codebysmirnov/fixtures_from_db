package database

import (
	"context"
	"fmt"
)

// getAllTableNamesQuery представляет SQL-запрос для получения всех имен таблиц в базе данных,
// исключая системные схемы pg_catalog и information_schema.
const getAllTableNamesQuery = `SELECT tablename
				FROM pg_tables
				WHERE schemaname NOT IN ('pg_catalog', 'information_schema')
				`

// GetAllTableNames возвращает все имена таблиц в базе данных, исключая системные схемы.
// Параметры:
// - ctx context.Context: контекст для управления временем выполнения запроса и его отменой.
// Возвращаемые значения:
// - ([]string, error): срез с именами таблиц или ошибка, если запрос не удался.
func (db *DB) GetAllTableNames(ctx context.Context) ([]string, error) {
	var allTables []string
	err := db.db.SelectContext(ctx, &allTables, getAllTableNamesQuery)
	if err != nil {
		return nil, fmt.Errorf("error getting all table names: %w", err)
	}
	return allTables, nil
}
