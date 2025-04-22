// Package middleware предоставляет middleware-компоненты для обработки HTTP-запросов.
// Включает в себя:
//   - JWT-аутентификацию
//   - Логирование
//   - Сжатие gzip
package middleware

import (
	"context"
	"net/http"

	"github.com/pervukhinpm/link-shortener.git/internal/jwt"
)

// CookieName представляет имя cookie для хранения JWT-токена.
// Используется для аутентификации пользователя.
const CookieName = "jwt"

// Auth является middleware для проверки аутентификации пользователя.
// Проверяет наличие и валидность JWT-токена в cookie.
// В случае успеха добавляет идентификатор пользователя в контекст запроса.
func Auth(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			var tokenString string

			cookie, err := r.Cookie(CookieName)

			if err != nil {
				tokenString, err = jwt.BuildJWTString()

				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				http.SetCookie(w, &http.Cookie{
					Name:  CookieName,
					Value: tokenString,
				})
			} else {
				tokenString = cookie.Value
			}

			userID, err := jwt.GetUserID(tokenString)

			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := setUserID(r.Context(), userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}(next)
}

// UserID представляет тип для хранения идентификатора пользователя в контексте запроса.
// Используется для передачи идентификатора пользователя между middleware и обработчиками.
type UserID struct {
	// ID - идентификатор пользователя
	ID string
}

// setUserID добавляет идентификатор пользователя в контекст запроса.
// Используется для передачи идентификатора между middleware и обработчиками.
func setUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserID{}, userID)
}

// GetUserID возвращает идентификатор пользователя из контекста запроса.
// Если идентификатор не найден, возвращает пустую строку.
func GetUserID(ctx context.Context) string {
	userID, ok := ctx.Value(UserID{}).(string)
	if !ok {
		return ""
	}

	return userID
}
