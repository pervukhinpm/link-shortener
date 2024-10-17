package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pervukhinpm/link-shortener.git/domain"
	"github.com/pervukhinpm/link-shortener.git/internal/errs"
	"github.com/pervukhinpm/link-shortener.git/internal/middleware"
	"github.com/pervukhinpm/link-shortener.git/internal/utils"
	"go.uber.org/zap"
)

type DatabaseRepository struct {
	db *pgxpool.Pool
}

func (dr *DatabaseRepository) Close() error {
	dr.db.Close()
	return nil
}

func NewDatabaseRepository(db *pgxpool.Pool) (*DatabaseRepository, error) {
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
	result, err := dr.db.Exec(ctx, query, uuid, url.ID, url.OriginalURL, userID, url.IsDeleted)

	if err != nil {
		middleware.Log.Error("Error inserting url", zap.Error(err))
		return err
	}

	rowsAffected := result.RowsAffected()

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
	err := dr.db.QueryRow(ctx, query, originalURL).Scan(&shortURL)
	if err != nil {
		return "", err
	}
	return shortURL, nil
}

func (dr *DatabaseRepository) Get(id string, ctx context.Context) (*domain.URL, error) {
	query := `
	SELECT original_url from urls WHERE short_url = $1;
	`
	originalURLRow := dr.db.QueryRow(ctx, query, id)

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
	_, err := dr.db.Exec(context.Background(), query)
	return err
}

func (dr *DatabaseRepository) AddBatch(urls []domain.URL, ctx context.Context) error {
	tx, err := dr.db.Begin(ctx)
	if err != nil {
		return err
	}

	userID := middleware.GetUserID(ctx)

	defer tx.Rollback(ctx)

	batch := &pgx.Batch{}
	query := "INSERT INTO urls (uuid, short_url, original_url, user_id, is_deleted)" +
		" VALUES ($1, $2, $3, $4, $5) ON CONFLICT (uuid) DO NOTHING"

	for _, v := range urls {
		uuid, err := utils.GenerateUUID()
		if err != nil {
			middleware.Log.Error("Error generating uuid", zap.Error(err))
			return err
		}
		batch.Queue(query, uuid, v.ID, v.OriginalURL, userID, v.IsDeleted)
	}
	dr.db.SendBatch(ctx, batch)
	return tx.Commit(ctx)
}

func (dr *DatabaseRepository) GetByUserID(ctx context.Context) (*[]domain.URL, error) {
	userID := middleware.GetUserID(ctx)

	query := `
    SELECT short_url, original_url, is_deleted FROM urls WHERE user_id = $1;
    `
	rows, err := dr.db.Query(ctx, query, userID)
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
        WHERE short_url = $1;
    `

	var isDeleted bool
	err := dr.db.QueryRow(ctx, query, shortenedURL).Scan(&isDeleted)
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
	query := `UPDATE urls SET is_deleted = $1 WHERE user_id = $2 AND short_url = $3;`
	batch := &pgx.Batch{}
	for _, v := range urls {
		batch.Queue(query, true, v.UserID, v.ShortURL)
	}

	br := dr.db.SendBatch(ctx, batch)
	defer br.Close()

	for _, v := range urls {
		_, err := br.Exec()
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				middleware.Log.Error("Error deleting short URL %v, %v\n", v.ShortURL, zap.Error(err))
			}
		}
	}
	return br.Close()
}
