package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pervukhinpm/link-shortener.git/domain"
	"github.com/pervukhinpm/link-shortener.git/internal/middleware"
	"github.com/pervukhinpm/link-shortener.git/internal/service"
	"go.uber.org/zap"
)

func TestCreateShortenerURL(t *testing.T) {
	urlService := service.NewMockService()
	baseURL := NewServerURL("http", "localhost", 8080)
	h := NewShortenerHandler(urlService, *baseURL)

	// Создаем мок логгера
	logger := zap.NewNop()
	middleware.Log = logger.Sugar()

	type want struct {
		contentType string
		statusCode  int
		bodyURL     string
		response    string
	}
	tests := []struct {
		name              string
		urlServiceShortID string
		contentType       string
		method            string
		setNilURL         bool
		want              want
	}{
		{
			name:              "positive test #1",
			urlServiceShortID: "testShortID",
			contentType:       "text/plain",
			method:            http.MethodPost,
			setNilURL:         false,
			want: want{
				contentType: "text/plain",
				bodyURL:     "https://practicum.yandex.ru/",
				statusCode:  http.StatusCreated,
				response:    "http://localhost:8080/testShortID",
			},
		},
		{
			name:              "empty body test #2",
			urlServiceShortID: "",
			contentType:       "text/plain",
			method:            http.MethodPost,
			setNilURL:         false,
			want: want{
				contentType: "text/plain",
				bodyURL:     "",
				statusCode:  http.StatusBadRequest,
				response:    "",
			},
		},
		{
			name:              "wrong HTTP method",
			urlServiceShortID: "testShortID",
			contentType:       "text/plain",
			method:            http.MethodGet,
			setNilURL:         false,
			want: want{
				contentType: "text/plain",
				bodyURL:     "https://practicum.yandex.ru/",
				statusCode:  http.StatusBadRequest,
				response:    "",
			},
		},
		{
			name:              "service error",
			urlServiceShortID: "testShortID",
			contentType:       "text/plain",
			method:            http.MethodPost,
			setNilURL:         true,
			want: want{
				contentType: "text/plain",
				bodyURL:     "https://practicum.yandex.ru/",
				statusCode:  http.StatusBadRequest,
				response:    "",
			},
		},
		{
			name:              "invalid content type",
			urlServiceShortID: "testShortID",
			contentType:       "application/json",
			method:            http.MethodPost,
			setNilURL:         false,
			want: want{
				contentType: "text/plain",
				bodyURL:     "https://practicum.yandex.ru/",
				statusCode:  http.StatusBadRequest,
				response:    "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := strings.NewReader(tt.want.bodyURL)
			req, err := http.NewRequest(tt.method, "http://localhost:8080/", body)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", tt.contentType)

			if tt.setNilURL {
				urlService.ShortenURL = nil
			} else {
				urlService.ShortenURL = domain.NewURL(tt.urlServiceShortID, tt.want.bodyURL, "", false)
			}

			buf, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatal(err)
			}
			req.Body = io.NopCloser(strings.NewReader(string(buf)))

			rr := httptest.NewRecorder()
			h.CreateShortenerURL(rr, req)

			if status := rr.Code; status != tt.want.statusCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.want.statusCode)
			}

			if string(buf) != tt.want.bodyURL {
				t.Errorf("handler returned wrong request body: got %v want %v",
					string(buf), tt.want.bodyURL)
			}

			if contentType := rr.Header().Get("Content-Type"); !strings.HasPrefix(contentType, tt.want.contentType) {
				t.Errorf("handler returned wrong content type: got %v want %v",
					contentType, tt.want.contentType)
			}

			if tt.want.response != "" {
				response := rr.Body.String()
				if response != tt.want.response {
					t.Errorf("handler returned unexpected body: got %v want %v",
						rr.Body.String(), tt.want.response)
				}
			}
		})
	}
}

func TestGetShortenerURL(t *testing.T) {
	urlService := service.NewMockService()
	baseURL := NewServerURL("http", "localhost", 8080)
	h := NewShortenerHandler(urlService, *baseURL)

	// Создаем мок логгера
	logger := zap.NewNop()
	middleware.Log = logger.Sugar()

	type want struct {
		statusCode int
		location   string
	}
	tests := []struct {
		name    string
		shortID string
		method  string
		want    want
	}{
		{
			name:    "positive test #1",
			shortID: "shortID",
			method:  http.MethodGet,
			want: want{
				statusCode: http.StatusTemporaryRedirect,
				location:   "https://practicum.yandex.ru/",
			},
		},
		{
			name:    "not found URL",
			shortID: "nonexistent",
			method:  http.MethodGet,
			want: want{
				statusCode: http.StatusBadRequest,
				location:   "",
			},
		},
		{
			name:    "wrong HTTP method",
			shortID: "shortID",
			method:  http.MethodPost,
			want: want{
				statusCode: http.StatusBadRequest,
				location:   "",
			},
		},
		{
			name:    "empty shortID",
			shortID: "",
			method:  http.MethodGet,
			want: want{
				statusCode: http.StatusBadRequest,
				location:   "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testURL := &domain.URL{
				ID:          tt.shortID,
				OriginalURL: "https://practicum.yandex.ru/",
			}

			if tt.shortID == "shortID" {
				urlService.ShortenURL = testURL
			} else {
				urlService.ShortenURL = nil
			}

			req, err := http.NewRequest(tt.method, "/"+tt.shortID, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			h.GetShortenerURL(rr, req)

			if status := rr.Code; status != tt.want.statusCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.want.statusCode)
			}
			if location := rr.Header().Get("Location"); location != tt.want.location {
				t.Errorf("handler returned wrong location header: got %v want %v",
					location, tt.want.location)
			}
		})
	}
}

func TestCreateJSONShortenerURL(t *testing.T) {
	urlService := service.NewMockService()
	baseURL := NewServerURL("http", "localhost", 8080)
	h := NewShortenerHandler(urlService, *baseURL)

	type want struct {
		contentType string
		statusCode  int
		response    string
	}
	tests := []struct {
		name        string
		requestBody string
		shortURL    string
		contentType string
		want        want
	}{
		{
			name:        "valid JSON request",
			requestBody: `{"url": "https://practicum.yandex.ru/"}`,
			shortURL:    "shortURL",
			contentType: "application/json",
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusCreated,
				response:    `{"result":"http://localhost:8080/shortURL"}`,
			},
		},
		{
			name:        "invalid content type",
			requestBody: `{"url": "https://practicum.yandex.ru/"}`,
			contentType: "text/plain",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
				response:    "Only application/json supported Media Type!\n",
			},
		},
		{
			name:        "empty URL field",
			requestBody: `{"url": ""}`,
			contentType: "application/json",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
				response:    "Empty URL!\n",
			},
		},
		{
			name:        "invalid JSON format",
			requestBody: `{"url": "https://practicum.yandex.ru/"`,
			contentType: "application/json",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
				response:    "unexpected end of JSON input\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testURL := &domain.URL{
				ID:          tt.shortURL,
				OriginalURL: "https://practicum.yandex.ru/",
			}
			if tt.shortURL != "" {
				urlService.ShortenURL = testURL
			}

			req, err := http.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(tt.requestBody))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", tt.contentType)

			rr := httptest.NewRecorder()
			h.CreateJSONShortenerURL(rr, req)

			if status := rr.Code; status != tt.want.statusCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.want.statusCode)
			}

			if contentType := rr.Header().Get("Content-Type"); contentType != tt.want.contentType {
				t.Errorf("handler returned wrong content type: got %v want %v",
					contentType, tt.want.contentType)
			}

			if rr.Body.String() != tt.want.response {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tt.want.response)
			}
		})
	}
}

func TestGetURLsByUser(t *testing.T) {
	urlService := service.NewMockService()
	baseURL := NewServerURL("http", "localhost", 8080)
	h := NewShortenerHandler(urlService, *baseURL)

	// Создаем мок логгера
	logger := zap.NewNop()
	middleware.Log = logger.Sugar()

	type want struct {
		statusCode  int
		contentType string
		response    string
	}

	tests := []struct {
		name   string
		url    *domain.URL
		method string
		want   want
	}{
		{
			name: "positive test with URL",
			url: &domain.URL{
				ID:          "short1",
				OriginalURL: "https://example.com/1",
			},
			method: http.MethodGet,
			want: want{
				statusCode:  http.StatusOK,
				contentType: "application/json",
				response:    `[{"short_url":"http://localhost:8080/short1","original_url":"https://example.com/1"}]`,
			},
		},
		{
			name:   "no URLs found",
			url:    nil,
			method: http.MethodGet,
			want: want{
				statusCode: http.StatusNoContent,
			},
		},
		{
			name: "wrong HTTP method",
			url: &domain.URL{
				ID:          "short1",
				OriginalURL: "https://example.com/1",
			},
			method: http.MethodPost,
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "malformed URL in response",
			url: &domain.URL{
				ID:          "short1",
				OriginalURL: "not a valid url",
			},
			method: http.MethodGet,
			want: want{
				statusCode:  http.StatusOK,
				contentType: "application/json",
				response:    `[{"short_url":"http://localhost:8080/short1","original_url":"not a valid url"}]`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urlService.ShortenURL = tt.url

			req, err := http.NewRequest(tt.method, "/api/user/urls", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			h.GetURLsByUser(rr, req)

			if status := rr.Code; status != tt.want.statusCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.want.statusCode)
			}

			if tt.want.contentType != "" {
				if contentType := rr.Header().Get("Content-Type"); contentType != tt.want.contentType {
					t.Errorf("handler returned wrong content type: got %v want %v",
						contentType, tt.want.contentType)
				}
			}

			if tt.want.response != "" {
				got := strings.TrimSpace(rr.Body.String())
				want := strings.TrimSpace(tt.want.response)
				if got != want {
					t.Errorf("handler returned unexpected body: got %v want %v",
						got, want)
				}
			}
		})
	}
}

func TestDeleteURLBatchByUser(t *testing.T) {
	urlService := service.NewMockService()
	baseURL := NewServerURL("http", "localhost", 8080)
	h := NewShortenerHandler(urlService, *baseURL)

	// Создаем мок логгера
	logger := zap.NewNop()
	middleware.Log = logger.Sugar()

	type want struct {
		statusCode int
	}

	tests := []struct {
		name        string
		requestBody string
		contentType string
		method      string
		want        want
	}{
		{
			name:        "valid delete request",
			requestBody: `["short1", "short2"]`,
			contentType: "application/json",
			method:      http.MethodDelete,
			want: want{
				statusCode: http.StatusAccepted,
			},
		},
		{
			name:        "invalid content type",
			requestBody: `["short1", "short2"]`,
			contentType: "text/plain",
			method:      http.MethodDelete,
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name:        "invalid JSON",
			requestBody: `["short1", "short2"`,
			contentType: "application/json",
			method:      http.MethodDelete,
			want: want{
				statusCode: http.StatusInternalServerError,
			},
		},
		{
			name:        "wrong HTTP method",
			requestBody: `["short1", "short2"]`,
			contentType: "application/json",
			method:      http.MethodPost,
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name:        "empty request body",
			requestBody: ``,
			contentType: "application/json",
			method:      http.MethodDelete,
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name:        "empty array",
			requestBody: `[]`,
			contentType: "application/json",
			method:      http.MethodDelete,
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name:        "invalid URL format",
			requestBody: `[123, 456]`,
			contentType: "application/json",
			method:      http.MethodDelete,
			want: want{
				statusCode: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, "/api/user/urls", strings.NewReader(tt.requestBody))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", tt.contentType)

			rr := httptest.NewRecorder()
			h.DeleteURLBatchByUser(rr, req)

			if status := rr.Code; status != tt.want.statusCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.want.statusCode)
			}
		})
	}
}
