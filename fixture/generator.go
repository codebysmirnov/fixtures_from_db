package fixture

// ColumnName представляет имя колонки в таблице базы данных.
type ColumnName string

// ColumnValue представляет значение колонки с указанием типа данных.
type ColumnValue struct {
	Type string      // Тип данных колонки (например, "string", "int").
	Data interface{} // Значение колонки.
}

// ValuesByColumn представляет карту значений, где ключом является имя колонки, а значением - ColumnValue.
type ValuesByColumn map[ColumnName]ColumnValue

// Fixture представляет структуру, содержащую данные для вставки в таблицу базы данных.
type Fixture struct {
	TableName string         // Имя таблицы, в которую будут вставлены данные.
	Data      ValuesByColumn // Данные для вставки, сгруппированные по колонкам.
}

// NewFixture создает новый объект Fixture с заданным именем таблицы и данными.
// Параметры:
// - tableName string: имя таблицы, в которую будут вставлены данные.
// - data ValuesByColumn: данные для вставки, сгруппированные по колонкам.
// Возвращаемое значение:
// - *Fixture: указатель на новый объект Fixture.
func NewFixture(tableName string, data ValuesByColumn) *Fixture {
	return &Fixture{
		TableName: tableName,
		Data:      data,
	}
}

// ColumnValueGenerator представляет интерфейс для генерации значений колонок.
type ColumnValueGenerator interface {
	Generate() ColumnValue
}
