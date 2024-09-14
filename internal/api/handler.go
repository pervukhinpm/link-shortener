package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/pervukhinpm/link-shortener.git/internal/model"
	"github.com/pervukhinpm/link-shortener.git/internal/repository"
	"github.com/pervukhinpm/link-shortener.git/internal/url"
	"io"
	"net/http"
	"strings"
)

type ShortenerHandler struct {
	urlService url.ShortenerServiceReaderWriter
	baseURL    ServerURL
	dsn        string
}

func NewHandler(urlService url.ShortenerServiceReaderWriter, baseURL ServerURL, dsn string) *ShortenerHandler {
	return &ShortenerHandler{
		urlService: urlService,
		baseURL:    baseURL,
		dsn:        dsn,
	}
}

func (h *ShortenerHandler) CreateShortenerURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(body) == 0 {
		http.Error(w, "Empty body!", http.StatusBadRequest)
		return
	}

	shortURL, err := h.urlService.Shorten(string(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, err = fmt.Fprintf(w, "%s/%s", h.baseURL.String(), shortURL.ID)
	if err != nil {
		return
	}
}

func (h *ShortenerHandler) GetShortenerURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusBadRequest)
		return
	}
	shortID := chi.URLParam(r, "id")

	origURL, err := h.urlService.Find(shortID)
	if err != nil {
		http.Error(w, "URL not found!", http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", origURL.OriginalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *ShortenerHandler) CreateJSONShortenerURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusBadRequest)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		http.Error(w, "Only application/json supported Media Type!", http.StatusBadRequest)
		return
	}

	var createShortenerBody model.CreateShortenerBody
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &createShortenerBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(createShortenerBody.URL) == 0 {
		http.Error(w, "Empty URL!", http.StatusBadRequest)
		return
	}

	shortURL, err := h.urlService.Shorten(createShortenerBody.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := fmt.Sprintf("%s/%s", h.baseURL.String(), shortURL.ID)
	response := model.CreateShortenerResponse{Result: result}

	jsonResp, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(jsonResp)
	if err != nil {
		return
	}
}

func (h *ShortenerHandler) PingDatabase(w http.ResponseWriter, r *http.Request) {
	if h.dsn == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err := repository.PingDatabase(h.dsn)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
