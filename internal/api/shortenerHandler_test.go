package api

import (
	"github.com/pervukhinpm/link-shortener.git/domain"
	"github.com/pervukhinpm/link-shortener.git/internal/service"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateShortenerURL(t *testing.T) {
	urlService := service.NewMockService()
	baseURL := NewServerURL("http", "localhost", 8080)
	h := NewHandler(urlService, *baseURL)

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
		want              want
	}{
		{
			name:              "positive test #1",
			urlServiceShortID: "testShortID",
			contentType:       "text/plain",
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
			want: want{
				contentType: "text/plain",
				bodyURL:     "",
				statusCode:  http.StatusBadRequest,
				response:    "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := strings.NewReader(tt.want.bodyURL)
			req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/", body)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", tt.contentType)

			urlService.ShortenURL = domain.NewURL(tt.urlServiceShortID, tt.want.bodyURL, "", false)

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
	h := NewHandler(urlService, *baseURL)

	type want struct {
		statusCode int
		location   string
	}
	tests := []struct {
		name    string
		shortID string
		want    want
	}{
		{
			name:    "positive test #1",
			shortID: "shortID",
			want: want{
				statusCode: http.StatusTemporaryRedirect,
				location:   "https://practicum.yandex.ru/",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			testURL := &domain.URL{
				ID:          tt.shortID,
				OriginalURL: "https://practicum.yandex.ru/",
			}

			if tt.shortID != "" {
				urlService.ShortenURL = testURL
			}

			req, err := http.NewRequest(http.MethodGet, "/"+tt.shortID, nil)
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
	h := NewHandler(urlService, *baseURL)

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
