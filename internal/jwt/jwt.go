// Package jwt предоставляет функциональность для работы с JWT-токенами.
// Включает в себя:
//   - Создание токенов
//   - Проверку токенов
//   - Управление временем жизни токенов
package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// Claims представляет набор утверждений (claims) JWT токена.
// Используется для хранения информации о пользователе в токене.
type Claims struct {
	// UserID - идентификатор пользователя
	UserID string `json:"user_id"`
	// StandardClaims - стандартные утверждения JWT
	jwt.RegisteredClaims
}

// Константы для JWT
const (
	// TokenExp - время жизни токена
	TokenExp = time.Hour * 3
	// SecretKey - секретный ключ для JWT
	SecretKey = "secret_key"
)

// BuildJWTString создает новый JWT-токен для указанного пользователя.
// В случае ошибки возвращает пустую строку и описание ошибки.
func BuildJWTString() (string, error) {
	userID := uuid.NewString()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GetUserID извлекает идентификатор пользователя из JWT-токена.
// В случае ошибки возвращает пустую строку и описание ошибки.
func GetUserID(tokenString string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signed method: %v", t.Header["alg"])
		}
		return []byte(SecretKey), nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", fmt.Errorf("token is not valid")
	}
	return claims.UserID, nil
}
