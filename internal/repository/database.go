package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/pervukhinpm/link-shortener.git/domain"
	"github.com/pervukhinpm/link-shortener.git/internal/errs"
	"github.com/pervukhinpm/link-shortener.git/internal/middleware"
	"github.com/pervukhinpm/link-shortener.git/internal/utils"
	"go.uber.org/zap"
)

type DatabaseRepository struct {
	db *sql.DB
}

func (dr *DatabaseRepository) Close() error {
	return dr.db.Close()
}

func NewDatabaseRepository(db *sql.DB) (*DatabaseRepository, error) {
	dbRepository := DatabaseRepository{
		db: db,
	}
	err := dbRepository.createDB()
	if err != nil {
		return nil, err
	}
	return &dbRepository, nil
}

func (dr *DatabaseRepository) Add(url *domain.URL, ctx context.Context) error {
	uuid, err := utils.GenerateUUID()
	if err != nil {
		middleware.Log.Error("Error generating uuid", zap.Error(err))
		return err
	}

	query := `
	INSERT INTO urls
	VALUES ($1, $2, $3, $4, $5)
    ON CONFLICT (original_url) DO NOTHING;
	`

	userID := middleware.GetUserID(ctx)
	result, err := dr.db.ExecContext(ctx, query, uuid, url.ID, url.OriginalURL, userID, url.IsDeleted)

	if err != nil {
		middleware.Log.Error("Error inserting url", zap.Error(err))
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		middleware.Log.Error("Error getting rows affected", zap.Error(err))
		return err
	}

	if rowsAffected == 0 {
		middleware.Log.Info("URL already exists, fetching existing short URL", zap.String("original_url", url.OriginalURL))

		existingShortURL, err := dr.getShortURLByOriginal(url.OriginalURL, ctx)
		if err != nil {
			middleware.Log.Error("Error getting existing short URL", zap.Error(err))
			return err
		}

		return errs.NewOriginalURLAlreadyExists(domain.NewURL(existingShortURL, url.OriginalURL, userID, false))
	}

	return nil
}

func (dr *DatabaseRepository) getShortURLByOriginal(originalURL string, ctx context.Context) (string, error) {
	query := `
    SELECT short_url FROM urls WHERE original_url = $1;
    `
	var shortURL string
	err := dr.db.QueryRowContext(ctx, query, originalURL).Scan(&shortURL)
	if err != nil {
		return "", err
	}
	return shortURL, nil
}

func (dr *DatabaseRepository) Get(id string, ctx context.Context) (*domain.URL, error) {
	query := `
	SELECT original_url from urls WHERE short_url = $1;
	`
	originalURLRow := dr.db.QueryRowContext(ctx, query, id)

	var originalURL string
	err := originalURLRow.Scan(&originalURL)
	if err != nil {
		return nil, err
	}

	userID := middleware.GetUserID(ctx)
	return domain.NewURL(id, originalURL, userID, false), nil
}

func (dr *DatabaseRepository) createDB() error {
	query := `
	CREATE TABLE IF NOT EXISTS urls (
		uuid varchar NOT NULL PRIMARY KEY,
		short_url varchar NOT NULL UNIQUE,
		original_url varchar NOT NULL UNIQUE,
		user_id varchar NOT NULL,
		is_deleted BOOLEAN NOT NULL DEFAULT FALSE
	);`
	_, err := dr.db.ExecContext(context.Background(), query)
	return err
}

func (dr *DatabaseRepository) AddBatch(urls []domain.URL, ctx context.Context) error {
	tx, err := dr.db.Begin()
	if err != nil {
		return err
	}

	userID := middleware.GetUserID(ctx)

	for _, v := range urls {
		uuid, err := utils.GenerateUUID()
		if err != nil {
			middleware.Log.Error("Error generating uuid", zap.Error(err))
			return err
		}

		_, err = tx.ExecContext(ctx, "INSERT INTO urls (uuid, short_url, original_url, user_id, is_deleted)"+
			" VALUES ($1, $2, $3, $4, $5) ON CONFLICT (uuid) DO NOTHING", uuid, v.ID, v.OriginalURL, userID, v.IsDeleted)
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				middleware.Log.Error("rollback transaction", zap.Error(err))
				return err
			}
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		middleware.Log.Error("commit transaction", zap.Error(err))
		return err
	}
	return err
}

func (dr *DatabaseRepository) GetByUserID(ctx context.Context) (*[]domain.URL, error) {
	userID := middleware.GetUserID(ctx)

	query := `
    SELECT short_url, original_url, is_deleted FROM urls WHERE user_id = $1;
    `
	rows, err := dr.db.QueryContext(ctx, query, userID)
	if err != nil {
		middleware.Log.Error("Error querying URLs by user ID", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var urls []domain.URL

	for rows.Next() {
		var shortURL, originalURL string
		var isDeleted bool
		err := rows.Scan(&shortURL, &originalURL, &isDeleted)
		if err != nil {
			middleware.Log.Error("Error scanning row", zap.Error(err))
			return nil, err
		}

		urls = append(urls, domain.URL{
			ID:          shortURL,
			OriginalURL: originalURL,
			UserID:      userID,
			IsDeleted:   isDeleted,
		})
	}

	if err := rows.Err(); err != nil {
		middleware.Log.Error("Error iterating over rows", zap.Error(err))
		return nil, err
	}

	return &urls, nil
}

func (dr *DatabaseRepository) GetFlagByShortURL(ctx context.Context, shortenedURL string) (bool, error) {
	query := `
        SELECT is_deleted
        FROM urls
        WHERE short_url = $1
    `

	var isDeleted bool
	err := dr.db.QueryRowContext(ctx, query, shortenedURL).Scan(&isDeleted)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, errs.ErrURLNotFound
		}
		middleware.Log.Error("Error querying short URL", zap.Error(err))
		return false, err
	}

	return isDeleted, nil
}

func (dr *DatabaseRepository) DeleteURLBatch(ctx context.Context, urls []UserShortURL) error {
	if len(urls) == 0 {
		return nil
	}

	// Создаём динамический запрос для каждого URL
	query := `
        UPDATE urls 
        SET is_deleted = TRUE 
        WHERE short_url IN (`

	args := make([]interface{}, len(urls)+1)
	for i, url := range urls {
		query += fmt.Sprintf("$%d,", i+1)
		args[i] = url.ShortURL
	}
	query = query[:len(query)-1]
	query += ") AND user_id = $" + fmt.Sprintf("%d", len(urls)+1)
	args[len(urls)] = urls[0].UserID

	_, err := dr.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
