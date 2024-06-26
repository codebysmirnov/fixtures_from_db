package fixture

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"fixtures_from_db/database"
)

const (
	maxRowsGenCount      = 1
	maxDefaultTextLength = 500
	openingParenthesis   = '('
	closingParenthesis   = ')'
	comma                = ','
	singleQuote          = '\''
	space                = " "
)

func addStringWithQuotes(b *strings.Builder, str string) {
	b.WriteRune(singleQuote)
	b.WriteString(str)
	b.WriteRune(singleQuote)
}

var varcharLengthConstraintRe = regexp.MustCompile("varchar\\((\\d+)\\)")
var numericAccuracyConstraintRe = regexp.MustCompile("numeric\\((\\d+),\\s*(\\d+)\\)")

func ToSQL(tableName string, columnsInfo []database.ColumnInfo) string {
	queryBuilder := strings.Builder{}
	queryBuilder.WriteString("INSERT INTO ")
	queryBuilder.WriteString(tableName)
	queryBuilder.WriteString(space)
	queryBuilder.WriteRune(openingParenthesis)
	for i := range columnsInfo {
		queryBuilder.WriteString(columnsInfo[i].ColumnName)
		if len(columnsInfo)-1 > i {
			queryBuilder.WriteRune(comma)
		}
	}
	queryBuilder.WriteRune(closingParenthesis)
	queryBuilder.WriteString(" VALUES ")
	for i := range maxRowsGenCount {
		queryBuilder.WriteRune(openingParenthesis)
		for j := range columnsInfo {
			dataType := columnsInfo[j].DataType
			switch dataType {
			case "uuid":
				addStringWithQuotes(&queryBuilder, uuid.New().String())
			case "int2":
				queryBuilder.WriteString(strconv.Itoa(int(rand.Int31n(math.MaxInt16))))
			case "int4":
				queryBuilder.WriteString(strconv.Itoa(int(rand.Int31())))
			case "int8":
				queryBuilder.WriteString(strconv.Itoa(int(rand.Int63())))
			case "bool":
				queryBuilder.WriteString(strconv.FormatBool(rand.Int31() > 0))
			case "text":
				str, err := randString(int(rand.Int63n(maxDefaultTextLength)))
				if err != nil {
					log.Fatalf("generate rand string failed: %v", err)
				}
				addStringWithQuotes(
					&queryBuilder, str)
			case "date":
				addStringWithQuotes(&queryBuilder, time.UnixMicro(rand.Int63()).Format(time.DateOnly))
			default:
				if !strings.Contains(dataType, "varchar") &&
					!strings.Contains(dataType, "timestamptz") &&
					!strings.Contains(dataType, "numeric") &&
					!strings.Contains(dataType, "jsonb") {
					fmt.Println(dataType)
				}
			}
			if strings.Contains(dataType, "jsonb") {
				addStringWithQuotes(&queryBuilder, "{}")
			}
			if strings.Contains(dataType, "numeric") {
				groups := numericAccuracyConstraintRe.FindStringSubmatch(dataType)
				if len(groups) == 0 {
					addStringWithQuotes(&queryBuilder, generateDefaultNumeric())
				} else {
					precision, err := strconv.Atoi(groups[1])
					if err != nil {
						log.Fatalf("get precision of numeric from string failed: %s", err)
					}
					scale, err := strconv.Atoi(groups[2])
					if err != nil {
						log.Fatalf("get scale of numeric from string failed: %s", err)
					}
					addStringWithQuotes(&queryBuilder, generateRandomNumeric(precision, scale))
				}
			}
			if strings.Contains(dataType, "timestamptz") {
				addStringWithQuotes(&queryBuilder, time.UnixMicro(rand.Int63()).String())
			}
			if strings.Contains(dataType, "varchar") {
				groups := varcharLengthConstraintRe.FindStringSubmatch(dataType)
				if len(groups) == 0 {
					str, err := randString(int(rand.Int63n(maxDefaultTextLength)))
					if err != nil {
						log.Fatalf("generate rand string failed: %v", err)
					}
					addStringWithQuotes(&queryBuilder, str)
				} else {
					symbolCount, err := strconv.Atoi(groups[1])
					if err != nil {
						log.Fatalf("get symbolNums failed: %s", err)
					}
					str, err := randString(int(rand.Int63n(int64(symbolCount))))
					if err != nil {
						log.Fatalf("generate rand string failed: %v", err)
					}
					addStringWithQuotes(&queryBuilder, str)
				}
			}
			if len(columnsInfo)-1 > j {
				queryBuilder.WriteRune(comma)
			}
		}
		queryBuilder.WriteRune(closingParenthesis)
		if maxRowsGenCount-1 > i {
			queryBuilder.WriteRune(comma)
		}
	}
	return queryBuilder.String()
}

const (
	maxDefaultNumericPrecision = 5
	maxDefaultNumericScale     = 2
)

// generateDefaultNumeric генерирует случайное число с максимальной точностью по умолчанию
func generateDefaultNumeric() string {
	return generateRandomNumeric(maxDefaultNumericPrecision, maxDefaultNumericScale)
}

// generateRandomNumeric генерирует случайное число в формате numeric(precision, scale)
func generateRandomNumeric(precision, scale int) string {
	if precision <= 0 || scale < 0 || scale > precision {
		log.Fatalf("Invalid precision(%d) or scale( %d)", precision, scale)
	}

	// Количество цифр до и после запятой
	intPartLen := precision - scale
	fracPartLen := scale

	// Генерация целой части
	intPart := make([]byte, intPartLen)
	for i := 0; i < intPartLen; i++ {
		intPart[i] = byte(rand.Intn(10) + '0')
	}

	// Генерация дробной части
	fracPart := make([]byte, fracPartLen)
	for i := 0; i < fracPartLen; i++ {
		fracPart[i] = byte(rand.Intn(10) + '0')
	}

	// Объединение частей в строку
	var result strings.Builder
	if intPartLen > 0 {
		result.Write(intPart)
	} else {
		result.WriteByte('0')
	}
	if fracPartLen > 0 {
		result.WriteByte('.')
		result.Write(fracPart)
	}

	return result.String()
}
