package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/pervukhinpm/link-shortener.git/domain"
	"github.com/pervukhinpm/link-shortener.git/internal/middleware"
	"github.com/pervukhinpm/link-shortener.git/internal/model"
	"github.com/pervukhinpm/link-shortener.git/internal/service"
	"io"
	"net/http"
	"strings"
)

type ShortenerHandler struct {
	urlService service.ShortenerServiceReaderWriter
	baseURL    ServerURL
}

func NewHandler(urlService service.ShortenerServiceReaderWriter, baseURL ServerURL) *ShortenerHandler {
	return &ShortenerHandler{
		urlService: urlService,
		baseURL:    baseURL,
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

	shortURL, err := h.urlService.Shorten(string(body), r.Context())
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

	origURL, err := h.urlService.Find(shortID, r.Context())
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

	shortURL, err := h.urlService.Shorten(createShortenerBody.URL, r.Context())
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

func (h *ShortenerHandler) BatchCreateJSONShortenerURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusBadRequest)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		http.Error(w, "Only application/json supported Media Type!", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Empty Body!", http.StatusBadRequest)
		return
	}
	if len(body) == 0 {
		http.Error(w, "Empty body!", http.StatusBadRequest)
		return
	}

	var batchRequestBody model.BatchRequestBody
	err = json.Unmarshal(body, &batchRequestBody.BatchList)
	if err != nil {
		http.Error(w, "Invalid JSON!", http.StatusBadRequest)
		return
	}

	batchRequestCount := len(batchRequestBody.BatchList)
	urls := make([]domain.URL, batchRequestCount)

	for i, v := range batchRequestBody.BatchList {
		urls[i] = *domain.NewURL(v.CorrelationID, v.OriginalURL)
	}

	err = h.urlService.AddBatch(urls, r.Context())
	if err != nil {
		middleware.Log.Error("Shortener.BatchCreateJSONShortenerURL: " + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	respData := make([]model.BatchResponseItem, batchRequestCount)
	for i, v := range urls {
		respData[i] = model.BatchResponseItem{
			CorrelationID: v.ID,
			ShortURL:      fmt.Sprintf("%s/%s", h.baseURL.String(), v.ID),
		}
	}
	jsonResp, err := json.Marshal(respData)
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
