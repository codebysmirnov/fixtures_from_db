package fixture

import (
	"crypto/rand"
	"math/big"
)

// letterBytes содержит символы, из которых будет состоять случайная строка.
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const letterBytesLength = int64(len(letterBytes))

// randString генерирует случайную строку заданной длины.
// Параметры:
// - n int: длина генерируемой строки
// Возвращает:
// - string: случайно сгенерированная строка
func randString(n int) (string, error) {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		idx, err := rand.Int(rand.Reader, big.NewInt(letterBytesLength))
		if err != nil {
			return "", err
		}
		b[i] = letterBytes[idx.Int64()]
	}
	return string(b), nil
}
