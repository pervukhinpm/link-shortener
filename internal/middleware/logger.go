// Package middleware предоставляет middleware-компоненты для обработки HTTP-запросов.
// Включает в себя:
//   - JWT-аутентификацию
//   - Логирование
//   - Сжатие gzip
package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Log представляет глобальный логгер для всего приложения.
// Используется для записи логов в формате JSON.
var Log *zap.SugaredLogger

// Initialize инициализирует глобальный логгер.
// Настраивает логгер для записи в формате JSON.
func Initialize() {
	zl, err := zap.NewProduction()
	if err != nil {
		panic("failed to initialize zap logger")
	}
	Log = zl.Sugar()
}

// loggingResponseWriter реализует http.ResponseWriter с поддержкой логирования.
// Используется для отслеживания статуса ответа и размера тела.
type loggingResponseWriter struct {
	// http.ResponseWriter - встроенный ResponseWriter
	http.ResponseWriter
	// statusCode - код статуса ответа
	statusCode int
	// bodySize - размер тела ответа
	bodySize int
}

// Write записывает данные в ответ и обновляет размер тела.
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.bodySize += size
	return size, err
}

// WriteHeader устанавливает код статуса ответа.
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.statusCode = statusCode
}

// Logger является middleware для логирования HTTP-запросов и ответов.
// Записывает информацию о:
//   - URI запроса
//   - Методе запроса
//   - Размере ответа
//   - Коде статуса
//   - Времени обработки
func Logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Log.Infow(
			"request",
			"uri", r.RequestURI,
			"method", r.Method,
		)

		lrw := loggingResponseWriter{ResponseWriter: w}

		start := time.Now()
		h.ServeHTTP(&lrw, r) // Вызов метода ServeHTTP для хендлера
		duration := time.Since(start)

		Log.Infow(
			"response",
			"size", lrw.bodySize,
			"status", lrw.statusCode,
			"duration", duration,
		)
	})
}
