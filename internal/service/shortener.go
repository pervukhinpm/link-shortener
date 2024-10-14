package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"github.com/pervukhinpm/link-shortener.git/domain"
	"github.com/pervukhinpm/link-shortener.git/internal/middleware"
	"github.com/pervukhinpm/link-shortener.git/internal/model"
	"github.com/pervukhinpm/link-shortener.git/internal/repository"
	"strings"
	"sync"
)

type ShortenerServiceReaderWriter interface {
	Find(id string, ctx context.Context) (*domain.URL, error)
	AddBatch(urls []domain.URL, ctx context.Context) error
	Shorten(original string, ctx context.Context) (*domain.URL, error)
	GetByUserID(ctx context.Context) (*[]domain.URL, error)
	DeleteURLBatch(ctx context.Context, deleteBatch model.DeleteBatch)
	GetFlagByShortURL(ctx context.Context, shortURL string) (bool, error)
}

type ShortenerService struct {
	repo repository.Repository
}

func NewURLService(repo repository.Repository) *ShortenerService {
	return &ShortenerService{repo: repo}
}

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

func (u *ShortenerService) AddBatch(urls []domain.URL, ctx context.Context) error {
	if err := u.repo.AddBatch(urls, ctx); err != nil {
		return err
	}
	return nil
}

func (u *ShortenerService) Find(id string, ctx context.Context) (*domain.URL, error) {
	url, err := u.repo.Get(id, ctx)
	if err != nil {
		return nil, err
	}
	return url, nil
}

func (u *ShortenerService) GetByUserID(ctx context.Context) (*[]domain.URL, error) {
	url, err := u.repo.GetByUserID(ctx)
	if err != nil {
		return nil, err
	}
	return url, nil
}

func (u *ShortenerService) GetFlagByShortURL(ctx context.Context, shortURL string) (bool, error) {
	isDeleted, err := u.repo.GetFlagByShortURL(ctx, shortURL)
	return isDeleted, err
}

func (u *ShortenerService) DeleteURLBatch(ctx context.Context, deleteBatch model.DeleteBatch) {
	doneCh := make(chan struct{})
	defer close(doneCh)

	inputCh := generator(doneCh, deleteBatch)
	channels := fanOut(doneCh, inputCh)
	formResultCh := fanIn(doneCh, channels...)

	var formedToDelete []repository.UserShortURL
	for form := range formResultCh {
		formedToDelete = append(formedToDelete, form)
	}

	err := u.repo.DeleteURLBatch(ctx, formedToDelete)
	if err != nil {
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

type DeleteTask struct {
	ShortenedURL string
	UserID       string
}
