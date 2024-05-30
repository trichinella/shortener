package random

import (
	"math/rand"
	"time"
)

// Честно взятая из интернета функция для получения рандомной строки
// Seed зафиксировал для упрощения работы с локальным API
func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)

	result := make([]byte, length)
	for i := range result {
		result[i] = charset[random.Intn(len(charset))]
	}

	return string(result)
}
