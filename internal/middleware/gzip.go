// Package middleware предоставляет middleware-компоненты для обработки HTTP-запросов.
// Включает в себя:
//   - JWT-аутентификацию
//   - Логирование
//   - Сжатие gzip
package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// compressWriter реализует http.ResponseWriter с поддержкой сжатия gzip.
// Используется для сжатия ответов сервера.
type compressWriter struct {
	// w - оригинальный ResponseWriter
	w http.ResponseWriter
	// zw - gzip.Writer для сжатия данных
	zw *gzip.Writer
}

// newCompressWriter создает новый экземпляр compressWriter.
// Инициализирует gzip.Writer для сжатия ответов.
func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

// Write записывает сжатые данные в ответ.
func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

// WriteHeader устанавливает HTTP-заголовки ответа.
func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 || statusCode == 409 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

// Header возвращает HTTP-заголовки ответа.
func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

// Close закрывает gzip.Writer.
func (c *compressWriter) Close() error {
	return c.zw.Close()
}

// compressReader реализует io.ReadCloser с поддержкой распаковки gzip.
// Используется для распаковки входящих запросов.
type compressReader struct {
	// r - оригинальный ReadCloser
	r io.ReadCloser
	// zr - gzip.Reader для распаковки данных
	zr *gzip.Reader
}

// newCompressReader создает новый экземпляр compressReader.
// Инициализирует gzip.Reader для распаковки входящих запросов.
func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

// Read читает и распаковывает данные из запроса.
func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

// Close закрывает gzip.Reader и оригинальный ReadCloser.
func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

// Gzip является middleware для обработки сжатия gzip.
// Обрабатывает:
//   - Сжатие ответов, если клиент поддерживает gzip (Accept-Encoding: gzip)
//   - Распаковку запросов, если они сжаты gzip (Content-Encoding: gzip или Content-Type: application/x-gzip)
func Gzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w

		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			cw := newCompressWriter(w)
			ow = cw
			defer cw.Close()
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		contentType := r.Header.Get("Content-Type")
		sendsGzip := strings.Contains(contentEncoding, "gzip") || contentType == "application/x-gzip"
		if sendsGzip {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				http.Error(w, "Invalid gzip data", http.StatusBadRequest)
				return
			}
			r.Body = cr
			defer cr.Close()
		}

		next.ServeHTTP(ow, r)
	})
}
