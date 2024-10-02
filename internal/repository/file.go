package repository

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"github.com/pervukhinpm/link-shortener.git/domain"
	"github.com/pervukhinpm/link-shortener.git/internal/errs"
	"github.com/pervukhinpm/link-shortener.git/internal/middleware"
	"github.com/pervukhinpm/link-shortener.git/internal/utils"
	"os"
)

type FileRepository struct {
	fileName string
	storage  map[string]URLFileModel
	writer   URLFileWriter
	reader   URLFileReader
}

func (r *FileRepository) Close() error {
	return r.reader.Close()
}

func NewFileRepository(fileName string) (*FileRepository, error) {
	writer, err := NewURLFileWriter(fileName)
	if err != nil {
		return nil, err
	}

	reader, err := NewURLFileReader(fileName)
	if err != nil {
		return nil, err
	}

	err = reader.ReadURL()
	if err != nil {
		return nil, err
	}

	repository := &FileRepository{
		fileName,
		make(map[string]URLFileModel),
		*writer,
		*reader,
	}

	for _, v := range reader.URLFileModels {
		repository.storage[v.ShortURL] = v
	}

	reader.Close()

	return repository, nil
}

func (r *FileRepository) Add(url *domain.URL, ctx context.Context) error {
	existingURL, _ := r.Get(url.ID, ctx)
	if existingURL != nil {
		return errs.NewOriginalURLAlreadyExists(existingURL)
	}
	uuid, err := utils.GenerateUUID()
	if err != nil {
		return err
	}
	urlFileModel := NewURLFileModel(uuid, url.ID, url.OriginalURL)
	err = r.writer.WriteURL(urlFileModel)
	if err != nil {
		return err
	}
	r.storage[url.ID] = *urlFileModel
	return nil
}

func (r *FileRepository) AddBatch(urls []domain.URL, ctx context.Context) error {
	for _, url := range urls {
		if err := r.Add(&url, ctx); err != nil {
			return err
		}
	}
	return nil
}

func (r *FileRepository) Get(id string, ctx context.Context) (*domain.URL, error) {
	userID := middleware.GetUserID(ctx)
	url, exists := r.storage[id]
	if !exists {
		return nil, errors.New("URL not found")
	}
	return domain.NewURL(id, url.OriginalURL, userID), nil
}

func (r *FileRepository) GetByUserID(ctx context.Context) (*[]domain.URL, error) {
	var urls []domain.URL

	userID := middleware.GetUserID(ctx)

	for _, record := range r.storage {
		if record.UserID == userID {
			url := domain.NewURL(record.ShortURL, record.OriginalURL, record.UserID)
			urls = append(urls, *url)
		}
	}

	// Если ничего не найдено, возвращаем пустой список
	return &urls, nil
}

type URLFileModel struct {
	UUID        string `json:"uuid"`
	UserID      string `json:"user_uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func NewURLFileModel(uuid, shortURL, originalURL string) *URLFileModel {
	return &URLFileModel{
		UUID:        uuid,
		ShortURL:    shortURL,
		OriginalURL: originalURL,
	}
}

type URLFileWriter struct {
	file   *os.File
	writer *bufio.Writer
}

func NewURLFileWriter(filename string) (*URLFileWriter, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return &URLFileWriter{
		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
}

func (u *URLFileWriter) WriteURL(fu *URLFileModel) error {
	data, err := json.Marshal(&fu)
	if err != nil {
		return err
	}
	if _, err := u.writer.Write(data); err != nil {
		return err
	}
	if err := u.writer.WriteByte('\n'); err != nil {
		return err
	}
	return u.writer.Flush()
}

type URLFileReader struct {
	file          *os.File
	scanner       *bufio.Scanner
	URLFileModels []URLFileModel
}

func NewURLFileReader(filename string) (*URLFileReader, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return &URLFileReader{
		file:          file,
		scanner:       bufio.NewScanner(file),
		URLFileModels: nil,
	}, nil
}

func (u *URLFileReader) ReadURL() error {
	u.URLFileModels = []URLFileModel{}
	for u.scanner.Scan() {
		data := u.scanner.Bytes()
		tempFormed := URLFileModel{}
		err := json.Unmarshal(data, &tempFormed)
		if err != nil {
			return err
		}
		u.URLFileModels = append(u.URLFileModels, tempFormed)
	}
	if err := u.scanner.Err(); err != nil {
		return u.scanner.Err()
	}
	return nil
}

func (u *URLFileReader) Close() error {
	return u.file.Close()
}
