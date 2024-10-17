package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/pervukhinpm/link-shortener.git/domain"
	"github.com/pervukhinpm/link-shortener.git/internal/errs"
	"github.com/pervukhinpm/link-shortener.git/internal/middleware"
	"github.com/pervukhinpm/link-shortener.git/internal/model"
	"github.com/pervukhinpm/link-shortener.git/internal/service"
	"go.uber.org/zap"
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
		if existingErr := new(errs.OriginalURLAlreadyExists); errors.As(err, &existingErr) {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusConflict)
			_, err = fmt.Fprintf(w, "%s/%s", h.baseURL.String(), existingErr.URL.ID)
			if err != nil {
				return
			}
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "text/plain")
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

	deleted, err := h.urlService.GetFlagByShortURL(r.Context(), shortID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if deleted {
		w.WriteHeader(http.StatusGone)
		return
	}

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
		if existingErr := new(errs.OriginalURLAlreadyExists); errors.As(err, &existingErr) {
			middleware.Log.Info("original url already exists", zap.String("url", existingErr.URL.OriginalURL))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			result := fmt.Sprintf("%s/%s", h.baseURL.String(), existingErr.URL.ID)
			response := model.CreateShortenerResponse{Result: result}
			jsonResp, err := json.Marshal(response)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_, err = w.Write(jsonResp)
			if err != nil {
				return
			}
			return
		}
		middleware.Log.Error("create shortener failed", zap.Error(err))
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

	userID := middleware.GetUserID(r.Context())

	for i, v := range batchRequestBody.BatchList {
		urls[i] = *domain.NewURL(v.CorrelationID, v.OriginalURL, userID, false)
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

func (h *ShortenerHandler) getURLsByUser(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie(middleware.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	urls, err := h.urlService.GetByUserID(r.Context())
	if err != nil {
		middleware.Log.Error("error to get url")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var shortURLBatch []model.URLByUserBatchResponseItem
	for _, url := range *urls {
		shortURLBatch = append(shortURLBatch, model.URLByUserBatchResponseItem{
			ShortURL:    fmt.Sprintf("%s/%s", h.baseURL.String(), url.ID),
			OriginalURL: url.OriginalURL,
		})
	}

	if len(shortURLBatch) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(shortURLBatch)
	if err != nil {
		middleware.Log.Error("error to create response", zap.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (h *ShortenerHandler) DeleteURLBatchByUser(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		http.Error(w, "Invalid Content-Type", http.StatusBadRequest)
		return
	}

	userID := middleware.GetUserID(r.Context())

	var deleteBatch model.DeleteBatch

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(bodyBytes, &deleteBatch.ShortenedURL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	deleteBatch.UserID = userID

	go h.urlService.DeleteURLBatch(r.Context(), deleteBatch)

	w.WriteHeader(http.StatusAccepted)
}
