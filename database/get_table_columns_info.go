package database

import (
	"context"
	"fmt"

	"github.com/lib/pq"
)

// ColumnInfo представляет структуру, содержащую информацию о колонке таблицы в базе данных.
type ColumnInfo struct {
	ColumnName        string         `db:"column_name"`
	DataType          string         `db:"data_type"`
	IsNotNull         bool           `db:"is_not_null"`
	DefaultValue      *string        `db:"default_value"`
	IsPrimaryKey      bool           `db:"is_primary_key"`
	IsForeignKey      bool           `db:"is_foreign_key"`
	ForeignTable      *string        `db:"foreign_table"`
	ForeignColumnName *string        `db:"foreign_column_name"`
	Values            pq.StringArray `db:"values"`
}

// Константа getTableColumnsInfoQuery представляет SQL-запрос для получения информации о колонках таблицы.
const getTableColumnsInfoQuery = `SELECT a.attname AS column_name,
       CASE
           WHEN t.typname = 'varchar' THEN t.typname || '(' || (a.atttypmod - 4) || ')'
           WHEN t.typname = 'bpchar' THEN 'char(' || (a.atttypmod - 4) || ')'
           WHEN t.typname = 'numeric' THEN t.typname || '(' || ((a.atttypmod - 4) >> 16) || ',' || ((a.atttypmod - 4) & 65535) || ')'
           WHEN t.typname = 'timestamp' THEN t.typname || '(' || (a.atttypmod - 4) || ')'
           WHEN t.typname = 'timestamptz' THEN t.typname || '(' || (a.atttypmod - 4) || ')'
           ELSE t.typname
           END AS data_type,
       a.attnotnull AS is_not_null,
       (SELECT pg_get_expr(d.adbin, d.adrelid)
        FROM pg_attrdef d
        WHERE d.adrelid = a.attrelid AND d.adnum = a.attnum) AS default_value,
       i.indisprimary IS NOT NULL AS is_primary_key,
       fk.conname IS NOT NULL AS is_foreign_key,
       fk.confrelid::regclass AS foreign_table,
       af.attname AS foreign_column_name,
       (SELECT array_agg(enumlabel)
        FROM pg_enum e
        WHERE e.enumtypid = t.oid) AS values
FROM pg_attribute a
         JOIN pg_class c ON a.attrelid = c.oid
         JOIN pg_type t ON a.atttypid = t.oid
         LEFT JOIN pg_index i ON c.oid = i.indrelid AND a.attnum = ANY(i.indkey)
         LEFT JOIN pg_constraint fk ON a.attnum = ANY(fk.conkey) AND fk.conrelid = c.oid AND fk.contype = 'f'
         LEFT JOIN pg_attribute af ON af.attnum = ANY(fk.confkey) AND af.attrelid = fk.confrelid
WHERE c.relname = $1
  AND a.attnum > 0
  AND NOT a.attisdropped
ORDER BY a.attnum;`

// GetTableColumnsInfo возвращает информацию о колонках указанной таблицы.
// Параметры:
// - ctx context.Context: контекст для управления временем выполнения запроса и его отменой.
// - table string: имя таблицы, для которой нужно получить информацию о колонках.
// Возвращаемые значения:
// - ([]ColumnInfo, error): срез структур ColumnInfo с информацией о колонках или ошибка, если запрос не удался.
func (db *DB) GetTableColumnsInfo(ctx context.Context, table string) ([]ColumnInfo, error) {
	var info []ColumnInfo
	err := db.db.SelectContext(ctx, &info, getTableColumnsInfoQuery, table)
	if err != nil {
		return nil, fmt.Errorf("failed to get table columns info: %w", err)
	}

	return info, nil
}
