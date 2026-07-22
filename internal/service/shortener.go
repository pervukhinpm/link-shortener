package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"strings"
	"sync"
	"time"

	"github.com/pervukhinpm/link-shortener.git/domain"
	"github.com/pervukhinpm/link-shortener.git/internal/middleware"
	"github.com/pervukhinpm/link-shortener.git/internal/model"
	"github.com/pervukhinpm/link-shortener.git/internal/repository"
)

// ShortenerServiceReaderWriter определяет интерфейс для работы с сокращенными URL.
// Предоставляет методы для создания, поиска, получения и удаления URL.
type ShortenerServiceReaderWriter interface {
	// Find ищет URL по его короткому идентификатору.
	Find(id string, ctx context.Context) (*domain.URL, error)
	// AddBatch добавляет несколько URL в хранилище.
	AddBatch(urls []domain.URL, ctx context.Context) error
	// Shorten создает сокращенный URL из оригинального.
	Shorten(original string, ctx context.Context) (*domain.URL, error)
	// GetByUserID возвращает все URL, созданные пользователем.
	GetByUserID(ctx context.Context) (*[]domain.URL, error)
	// DeleteURLBatch удаляет несколько URL пользователя.
	DeleteURLBatch(ctx context.Context, deleteBatch model.DeleteBatch)
	// GetFlagByShortURL проверяет, был ли URL удален.
	GetFlagByShortURL(ctx context.Context, shortURL string) (bool, error)
}

// ShortenerService реализует интерфейс ShortenerServiceReaderWriter.
// Предоставляет основную бизнес-логику для работы с сокращенными URL.
type ShortenerService struct {
	repo repository.Repository
}

// NewURLService создает новый экземпляр ShortenerService с указанным репозиторием.
func NewURLService(repo repository.Repository) *ShortenerService {
	return &ShortenerService{repo: repo}
}

// Find возвращает URL по его короткому идентификатору.
// Если URL не найден или помечен как удаленный, возвращает ошибку.
func (u *ShortenerService) Find(id string, ctx context.Context) (*domain.URL, error) {
	url, err := u.repo.Get(id, ctx)
	if err != nil {
		return nil, err
	}
	return url, nil
}

// Shorten создает новый сокращенный URL из оригинального.
// Генерирует уникальный короткий идентификатор и сохраняет URL в хранилище.
func (u *ShortenerService) Shorten(original string, ctx context.Context) (*domain.URL, error) {
	userID := middleware.GetUserID(ctx)
	randomBytes := make([]byte, 6)
	if _, err := rand.Read(randomBytes); err != nil {
		return nil, err
	}
	short := base64.URLEncoding.EncodeToString(randomBytes)
	short = strings.TrimRight(short, "=")
	url := domain.NewURL(short, original, userID, false)
	if err := u.repo.Add(url, ctx); err != nil {
		return nil, err
	}
	return url, nil
}

// AddBatch добавляет массив URL в хранилище.
// Каждому URL присваивается уникальный короткий идентификатор.
func (u *ShortenerService) AddBatch(urls []domain.URL, ctx context.Context) error {
	return u.repo.AddBatch(urls, ctx)
}

// GetByUserID возвращает все URL, созданные пользователем.
// Если у пользователя нет URL, возвращает пустой слайс.
func (u *ShortenerService) GetByUserID(ctx context.Context) (*[]domain.URL, error) {
	url, err := u.repo.GetByUserID(ctx)
	if err != nil {
		return nil, err
	}
	return url, nil
}

// GetFlagByShortURL проверяет, был ли URL удален.
// Использует контекст с таймаутом для предотвращения зависания.
func (u *ShortenerService) GetFlagByShortURL(ctx context.Context, shortURL string) (bool, error) {
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	isDeleted, err := u.repo.GetFlagByShortURL(ctxWithTimeout, shortURL)
	return isDeleted, err
}

// DeleteURLBatch асинхронно удаляет URL пользователя.
// Помечает указанные URL как удаленные в хранилище.
func (u *ShortenerService) DeleteURLBatch(ctx context.Context, deleteBatch model.DeleteBatch) {
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	formedToDelete := make([]repository.UserShortURL, 0, len(deleteBatch.ShortenedURL))

	const batchSize = 1000
	for i := 0; i < len(deleteBatch.ShortenedURL); i += batchSize {
		end := i + batchSize
		if end > len(deleteBatch.ShortenedURL) {
			end = len(deleteBatch.ShortenedURL)
		}

		batch := deleteBatch.ShortenedURL[i:end]
		batchFormed := make([]repository.UserShortURL, len(batch))

		var wg sync.WaitGroup
		wg.Add(len(batch))

		for j := range batch {
			go func(j int) {
				defer wg.Done()
				batchFormed[j] = repository.UserShortURL{
					UserID:   deleteBatch.UserID,
					ShortURL: batch[j],
				}
			}(j)
		}

		wg.Wait()
		formedToDelete = append(formedToDelete, batchFormed...)
	}

	if err := u.repo.DeleteURLBatch(ctxWithTimeout, formedToDelete); err != nil {
		return
	}
}

func generator(doneCh chan struct{}, input model.DeleteBatch) chan DeleteTask {
	inputCh := make(chan DeleteTask)

	go func() {
		defer close(inputCh)

		for _, data := range input.ShortenedURL {
			task := DeleteTask{
				UserID:       input.UserID,
				ShortenedURL: data,
			}
			select {
			case <-doneCh:
				return
			case inputCh <- task:
			}

		}
	}()

	return inputCh
}

func form(doneCh chan struct{}, inputCh chan DeleteTask) chan repository.UserShortURL {
	formRes := make(chan repository.UserShortURL)

	go func() {
		defer close(formRes)

		for data := range inputCh {

			formed := repository.UserShortURL{
				UserID:   data.UserID,
				ShortURL: data.ShortenedURL,
			}

			select {
			case <-doneCh:
				return
			case formRes <- formed:
			}
		}
	}()
	return formRes
}

func fanOut(doneCh chan struct{}, inputCh chan DeleteTask) []chan repository.UserShortURL {
	// количество горутин add
	numWorkers := 10
	// каналы, в которые отправляются результаты
	channels := make([]chan repository.UserShortURL, numWorkers)

	for i := 0; i < numWorkers; i++ {
		// получаем канал из горутины add
		formResultCh := form(doneCh, inputCh)
		// отправляем его в слайс каналов
		channels[i] = formResultCh
	}

	// возвращаем слайс каналов
	return channels
}

// fanIn объединяет несколько каналов resultChs в один.
func fanIn(doneCh chan struct{}, resultChs ...chan repository.UserShortURL) chan repository.UserShortURL {
	// конечный выходной канал в который отправляем данные из всех каналов из слайса, назовём его результирующим
	finalCh := make(chan repository.UserShortURL)

	// понадобится для ожидания всех горутин
	var wg sync.WaitGroup

	// перебираем все входящие каналы
	for _, ch := range resultChs {
		// в горутину передавать переменную цикла нельзя, поэтому делаем так
		chClosure := ch

		// инкрементируем счётчик горутин, которые нужно подождать
		wg.Add(1)

		go func() {
			// откладываем сообщение о том, что горутина завершилась
			defer wg.Done()

			// получаем данные из канала
			for data := range chClosure {
				select {
				// выходим из горутины, если канал закрылся
				case <-doneCh:
					return
				// если не закрылся, отправляем данные в конечный выходной канал
				case finalCh <- data:
				}
			}
		}()
	}

	go func() {
		// ждём завершения всех горутин
		wg.Wait()
		// когда все горутины завершились, закрываем результирующий канал
		close(finalCh)
	}()

	// возвращаем результирующий канал
	return finalCh
}

// DeleteTask представляет задачу на удаление URL.
// Используется для асинхронной обработки запросов на удаление URL.
type DeleteTask struct {
	UserID       string
	ShortenedURL string
}
