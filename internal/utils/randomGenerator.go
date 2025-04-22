// Package utils предоставляет вспомогательные функции для работы с данными.
// Включает в себя:
//   - Генерацию UUID
//   - Работу со случайными числами
package utils

import (
	"crypto/rand"
	"fmt"
	"io"
)

// GenerateUUID создает новый UUID версии 4.
// Использует криптографически безопасный генератор случайных чисел.
// В случае ошибки возвращает пустую строку и описание проблемы.
func GenerateUUID() (string, error) {
	uuid := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, uuid)
	if err != nil {
		return "", err
	}

	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%12x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}
