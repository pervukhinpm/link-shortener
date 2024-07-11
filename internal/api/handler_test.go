package api

import (
	"github.com/pervukhinpm/link-shortener.git/domain"
	"github.com/pervukhinpm/link-shortener.git/internal/url"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateShortenerURL(t *testing.T) {
	urlService := url.NewMockService()
	h := NewHandler(urlService)

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

			urlService.ShortenUrl = domain.NewURL(tt.urlServiceShortID, tt.want.bodyURL)

			buf, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatal(err)
			}
			req.Body = io.NopCloser(strings.NewReader(string(buf)))

			rr := httptest.NewRecorder()
			h.ServeHTTP(rr, req)

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
	urlService := url.NewMockService()
	h := NewHandler(urlService)

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
				urlService.ShortenUrl = testURL
			}

			req, err := http.NewRequest(http.MethodGet, "/"+tt.shortID, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			h.ServeHTTP(rr, req)

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
