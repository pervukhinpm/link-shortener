package repository

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/pervukhinpm/link-shortener.git/domain"
	"github.com/pervukhinpm/link-shortener.git/internal/errs"
	"github.com/pervukhinpm/link-shortener.git/internal/middleware"
	"github.com/pervukhinpm/link-shortener.git/internal/utils"
)

// FileRepository реализует интерфейс Repository для работы с файловым хранилищем.
// Использует JSON для сериализации данных.
type FileRepository struct {
	fileName string
	storage  map[string]URLFileModel
	writer   URLFileWriter
	reader   URLFileReader
}

// Close закрывает файловое хранилище.
func (r *FileRepository) Close() error {
	return r.reader.Close()
}

// NewFileRepository создает новый экземпляр FileRepository.
// Инициализирует хранилище и загружает существующие данные.
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

// Add добавляет новый URL в файловое хранилище.
// Проверяет уникальность URL перед добавлением.
func (r *FileRepository) Add(url *domain.URL, ctx context.Context) error {
	existingURL, _ := r.Get(url.ID, ctx)
	if existingURL != nil {
		return errs.NewOriginalURLAlreadyExists(existingURL)
	}
	uuid, err := utils.GenerateUUID()
	if err != nil {
		return err
	}
	urlFileModel := NewURLFileModel(uuid, url.ID, url.OriginalURL, false)
	err = r.writer.WriteURL(urlFileModel)
	if err != nil {
		return err
	}
	r.storage[url.ID] = *urlFileModel
	return nil
}

// AddBatch добавляет несколько URL в файловое хранилище.
// Вызывает Add для каждого URL в пакете.
func (r *FileRepository) AddBatch(urls []domain.URL, ctx context.Context) error {
	for _, url := range urls {
		if err := r.Add(&url, ctx); err != nil {
			return err
		}
	}
	return nil
}

// Get возвращает URL по его короткому идентификатору.
// Если URL не найден, возвращает ошибку.
func (r *FileRepository) Get(id string, ctx context.Context) (*domain.URL, error) {
	userID := middleware.GetUserID(ctx)
	url, exists := r.storage[id]
	if !exists {
		return nil, errors.New("URL not found")
	}
	return domain.NewURL(id, url.OriginalURL, userID, url.IsDeleted), nil
}

// GetByUserID возвращает все URL, созданные пользователем.
// Если URL не найдены, возвращает пустой слайс.
func (r *FileRepository) GetByUserID(ctx context.Context) (*[]domain.URL, error) {
	var urls []domain.URL

	userID := middleware.GetUserID(ctx)

	for _, record := range r.storage {
		if record.UserID == userID {
			url := domain.NewURL(record.ShortURL, record.OriginalURL, record.UserID, record.IsDeleted)
			urls = append(urls, *url)
		}
	}

	// Если ничего не найдено, возвращаем пустой список
	return &urls, nil
}

// DeleteURLBatch помечает несколько URL как удаленные.
// Перезаписывает файл с обновленными данными.
func (r *FileRepository) DeleteURLBatch(ctx context.Context, urls []UserShortURL) error {
	userID := middleware.GetUserID(ctx)

	for _, url := range urls {
		storedURL, exists := r.storage[url.ShortURL]
		if exists && storedURL.UserID == userID {
			storedURL.IsDeleted = true
			r.storage[url.ShortURL] = storedURL
		}
	}

	return r.rewriteFile()
}

// GetFlagByShortURL проверяет, был ли URL удален.
// Если URL не найден, возвращает ошибку.
func (r *FileRepository) GetFlagByShortURL(_ context.Context, shortenedURL string) (bool, error) {
	return r.storage[shortenedURL].IsDeleted, nil
}

// rewriteFile перезаписывает файл с текущими данными хранилища.
// Используется после операций изменения данных.
func (r *FileRepository) rewriteFile() error {
	fileWriter, err := NewURLFileWriter(r.fileName)
	if err != nil {
		return err
	}
	defer fileWriter.file.Close()

	for _, urlModel := range r.storage {
		err := fileWriter.WriteURL(&urlModel)
		if err != nil {
			return err
		}
	}

	return nil
}

// URLFileModel представляет структуру для хранения URL в файле.
// Содержит все необходимые поля для работы с URL.
type URLFileModel struct {
	UUID        string `json:"uuid"`
	UserID      string `json:"user_uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	IsDeleted   bool   `json:"is_deleted"`
}

// NewURLFileModel создает новый экземпляр URLFileModel.
func NewURLFileModel(uuid, shortURL, originalURL string, isDeleted bool) *URLFileModel {
	return &URLFileModel{
		UUID:        uuid,
		ShortURL:    shortURL,
		OriginalURL: originalURL,
		IsDeleted:   isDeleted,
	}
}

// URLFileWriter реализует запись URL в файл.
// Использует буферизованную запись для оптимизации производительности.
type URLFileWriter struct {
	file   *os.File
	writer *bufio.Writer
}

// NewURLFileWriter создает новый экземпляр URLFileWriter.
// Открывает файл для записи в режиме добавления.
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

// WriteURL записывает URL в файл в формате JSON.
// Каждая запись добавляется в новую строку.
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

// URLFileReader реализует чтение URL из файла.
// Использует буферизованное чтение для оптимизации производительности.
type URLFileReader struct {
	file          *os.File
	scanner       *bufio.Scanner
	URLFileModels []URLFileModel
}

// NewURLFileReader создает новый экземпляр URLFileReader.
// Открывает файл для чтения.
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

// ReadURL читает все URL из файла.
// Парсит JSON и сохраняет данные в память.
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

// Close закрывает файловый поток.
func (u *URLFileReader) Close() error {
	return u.file.Close()
}
