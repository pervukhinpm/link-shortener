package api_test

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"

	"github.com/pervukhinpm/link-shortener.git/domain"
	"github.com/pervukhinpm/link-shortener.git/internal/api"
	"github.com/pervukhinpm/link-shortener.git/internal/service"
)

func ExampleShortenerHandler_CreateShortenerURL() {
	// Создаем тестовый сервис
	urlService := service.NewMockService()
	baseURL := api.NewServerURL("http", "localhost", 8080)
	handler := api.NewShortenerHandler(urlService, *baseURL)

	// Настраиваем мок-сервис
	urlService.ShortenURL = &domain.URL{
		ID:          "testShortID",
		OriginalURL: "https://example.com",
	}

	// Создаем тестовый запрос
	req := httptest.NewRequest("POST", "/", strings.NewReader("https://example.com"))
	req.Header.Set("Content-Type", "text/plain")
	w := httptest.NewRecorder()

	// Вызываем обработчик
	handler.CreateShortenerURL(w, req)

	// Проверяем ответ
	resp := w.Result()
	fmt.Printf("Status: %d\n", resp.StatusCode)
	fmt.Printf("Body: %s\n", w.Body.String())

	// Output:
	// Status: 201
	// Body: http://localhost:8080/testShortID
}

func ExampleShortenerHandler_CreateShortenerURL_gzip() {
	// Создаем тестовый сервис
	urlService := service.NewMockService()
	baseURL := api.NewServerURL("http", "localhost", 8080)
	handler := api.NewShortenerHandler(urlService, *baseURL)

	// Настраиваем мок-сервис
	urlService.ShortenURL = &domain.URL{
		ID:          "testShortID",
		OriginalURL: "https://example.com",
	}

	// Создаем gzip-сжатый запрос
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	zw.Write([]byte("https://example.com"))
	zw.Close()

	req := httptest.NewRequest("POST", "/", &buf)
	req.Header.Set("Content-Type", "application/x-gzip")
	req.Header.Set("Content-Encoding", "gzip")
	w := httptest.NewRecorder()

	// Вызываем обработчик
	handler.CreateShortenerURL(w, req)

	// Проверяем ответ
	resp := w.Result()
	fmt.Printf("Status: %d\n", resp.StatusCode)
	fmt.Printf("Body: %s\n", w.Body.String())

	// Output:
	// Status: 201
	// Body: http://localhost:8080/testShortID
}

func ExampleShortenerHandler_CreateJSONShortenerURL() {
	// Создаем тестовый сервис
	urlService := service.NewMockService()
	baseURL := api.NewServerURL("http", "localhost", 8080)
	handler := api.NewShortenerHandler(urlService, *baseURL)

	// Настраиваем мок-сервис
	urlService.ShortenURL = &domain.URL{
		ID:          "testShortID",
		OriginalURL: "https://example.com",
	}

	// Создаем JSON запрос
	body := map[string]string{"url": "https://example.com"}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/api/shorten", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Вызываем обработчик
	handler.CreateJSONShortenerURL(w, req)

	// Проверяем ответ
	resp := w.Result()
	fmt.Printf("Status: %d\n", resp.StatusCode)
	fmt.Printf("Body: %s\n", w.Body.String())

	// Output:
	// Status: 201
	// Body: {"result":"http://localhost:8080/testShortID"}
}

func ExampleShortenerHandler_GetShortenerURL() {
	// Создаем тестовый сервис
	urlService := service.NewMockService()
	baseURL := api.NewServerURL("http", "localhost", 8080)
	handler := api.NewShortenerHandler(urlService, *baseURL)

	// Настраиваем мок-сервис
	urlService.ShortenURL = &domain.URL{
		ID:          "testShortID",
		OriginalURL: "https://example.com",
	}

	// Создаем тестовый запрос
	req := httptest.NewRequest("GET", "/testShortID", nil)
	w := httptest.NewRecorder()

	// Вызываем обработчик
	handler.GetShortenerURL(w, req)

	// Проверяем ответ
	resp := w.Result()
	fmt.Printf("Status: %d\n", resp.StatusCode)
	fmt.Printf("Location: %s\n", resp.Header.Get("Location"))

	// Output:
	// Status: 307
	// Location: https://example.com
}

func ExampleShortenerHandler_GetURLsByUser() {
	// Создаем тестовый сервис
	urlService := service.NewMockService()
	baseURL := api.NewServerURL("http", "localhost", 8080)
	handler := api.NewShortenerHandler(urlService, *baseURL)

	// Настраиваем мок-сервис
	urlService.ShortenURL = &domain.URL{
		ID:          "testShortID",
		OriginalURL: "https://example.com",
	}

	// Создаем тестовый запрос
	req := httptest.NewRequest("GET", "/api/user/urls", nil)
	w := httptest.NewRecorder()

	// Вызываем обработчик
	handler.GetURLsByUser(w, req)

	// Проверяем ответ
	resp := w.Result()
	fmt.Printf("Status: %d\n", resp.StatusCode)
	fmt.Printf("Body: %s\n", w.Body.String())

	// Output:
	// Status: 200
	// Body: [{"short_url":"http://localhost:8080/testShortID","original_url":"https://example.com"}]
}

func ExampleShortenerHandler_DeleteURLBatchByUser() {
	// Создаем тестовый сервис
	urlService := service.NewMockService()
	baseURL := api.NewServerURL("http", "localhost", 8080)
	handler := api.NewShortenerHandler(urlService, *baseURL)

	// Создаем JSON запрос
	body := []string{"http://localhost:8080/testShortID"}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("DELETE", "/api/user/urls", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Вызываем обработчик
	handler.DeleteURLBatchByUser(w, req)

	// Проверяем ответ
	resp := w.Result()
	fmt.Printf("Status: %d\n", resp.StatusCode)

	// Output:
	// Status: 202
}
