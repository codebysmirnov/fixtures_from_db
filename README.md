# Генерация fixtures для базы данных

Утилита позволяет сгенерировать произвольные данные основываясь на схеме базы данных.

Основной принцип работы:
1. При старте утилиты нужно передать строку для подключения к базе данных. 
2. Утилита проанализирует схему: таблицы, их связи с другими таблицами, колонки, их типы, ограничения и точность.
3. Вернет INSERT statements в STDOUT.

Поддерживаемые типы данных

| Тип данных                            | Поддержка |
|---------------------------------------|-----------|
| TEXT, VARCHAR                         | +         |
| uuid                                  | +         |
| smallint, integer, bigint             | +         |
| bool                                  | +         |
| date                                  | +         |
| timestamp, timestamptz                | +         |
| jsonb                                 | +         |
| numeric                               | +         |
| jsonb                                 | +         |
| serial, bigserial                     | -         |
| point, box, path, line, lseg, polygon | -         |
| bytea, array                          | -         |
| cidr, inet                            | -         |
| interval                              | -         |
| macaddr, macaddr8                     | -         |
| money                                 | -         |
| real                                  | -         |
| time [with time zone]                 | -         |
| tsquery, tsvector                     | -         |
| xml                                   | -         |



_Вывод можно перенаправить как в файл, так и сразу в бд._

_На текущий момент поддерживается PostgreSQL. Тестировалось на версиях 14.11, 16.2_

# Запуск

Для генерации конфигуарции фикстур необходимо воспользоваться командой introspect, 
указав строку для подключения к базе данных.

```shell
go build .
./fixtures_from_db introspect -c 'host=localhost
port=5432
user=username
dbname=database_name
password=passwod
sslmode=disable'
```

