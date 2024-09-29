package repository

import (
	"context"
	"database/sql"
	"github.com/pervukhinpm/link-shortener.git/domain"
	"github.com/pervukhinpm/link-shortener.git/internal/errs"
	"github.com/pervukhinpm/link-shortener.git/internal/middleware"
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
	uuid, err := GenerateUUID()
	if err != nil {
		middleware.Log.Error("Error generating uuid", zap.Error(err))
		return err
	}

	query := `
	INSERT INTO urls
	VALUES ($1, $2, $3)
    ON CONFLICT (original_url) DO NOTHING;
	`

	result, err := dr.db.ExecContext(ctx, query, uuid, url.ID, url.OriginalURL)

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

		return errs.NewOriginalURLAlreadyExists(domain.NewURL(existingShortURL, url.OriginalURL))
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
	return domain.NewURL(id, originalURL), nil
}

func (dr *DatabaseRepository) createDB() error {
	query := `
	CREATE TABLE IF NOT EXISTS urls (
		uuid varchar NOT NULL PRIMARY KEY,
		short_url varchar NOT NULL,
		original_url varchar NOT NULL UNIQUE
	);`
	_, err := dr.db.ExecContext(context.Background(), query)
	return err
}

func (dr *DatabaseRepository) AddBatch(urls []domain.URL, ctx context.Context) error {
	tx, err := dr.db.Begin()
	if err != nil {
		return err
	}
	for _, v := range urls {
		uuid, err := GenerateUUID()
		if err != nil {
			middleware.Log.Error("Error generating uuid", zap.Error(err))
			return err
		}

		_, err = tx.ExecContext(ctx, "INSERT INTO urls (uuid, short_url, original_url)"+
			" VALUES ($1, $2, $3) ON CONFLICT (uuid) DO NOTHING", uuid, v.ID, v.OriginalURL)
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
