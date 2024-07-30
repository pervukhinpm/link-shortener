package repository

import (
	"bufio"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pervukhinpm/link-shortener.git/domain"
	"io"
	"os"
)

type Repository interface {
	Add(url *domain.URL) error
	Get(id string) (*domain.URL, error)
}

type FileRepository struct {
	fileName string
	storage  map[string]string
	writer   URLFileWriter
	reader   URLFileReader
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
		make(map[string]string),
		*writer,
		*reader,
	}

	for _, v := range reader.URLFileModels {
		repository.storage[v.ShortUrl] = v.OriginalUrl
	}

	reader.Close()

	return repository, nil
}

func (r *FileRepository) Add(url *domain.URL) error {
	uuid, err := GenerateUUID()
	if err != nil {
		return err
	}
	urlFileModel := NewURLFileModel(uuid, url.ID, url.OriginalURL)
	err = r.writer.WriteURL(urlFileModel)
	if err != nil {
		return err
	}
	r.storage[url.ID] = url.OriginalURL
	return nil
}

func (r *FileRepository) Get(id string) (*domain.URL, error) {
	url, exists := r.storage[id]
	if !exists {
		return nil, errors.New("URL not found")
	}
	return domain.NewURL(id, url), nil
}

type URLFileModel struct {
	Uuid        string `json:"uuid"`
	ShortUrl    string `json:"short_url"`
	OriginalUrl string `json:"original_url"`
}

func NewURLFileModel(uuid, shortUrl, originalUrl string) *URLFileModel {
	return &URLFileModel{
		Uuid:        uuid,
		ShortUrl:    shortUrl,
		OriginalUrl: originalUrl,
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

func GenerateUUID() (string, error) {
	uuid := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, uuid)
	if err != nil {
		return "", err
	}

	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%12x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}
