package service

import (
	"context"
	"log"
	"os"
	"runtime/pprof"
	"testing"

	"github.com/pervukhinpm/link-shortener.git/domain"
	"github.com/pervukhinpm/link-shortener.git/internal/model"
	"github.com/pervukhinpm/link-shortener.git/internal/repository"
)

func BenchmarkShortenerService(b *testing.B) {
	if err := os.MkdirAll("profiles", 0755); err != nil {
		log.Fatal(err)
	}

	baseProfile, err := os.Create("profiles/base.pprof")
	if err != nil {
		log.Fatal(err)
	}
	defer baseProfile.Close()

	if err := pprof.WriteHeapProfile(baseProfile); err != nil {
		log.Fatal(err)
	}

	mockRepo := repository.NewMockRepository()
	service := NewURLService(mockRepo)
	ctx := context.Background()

	b.Run("Shorten", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := service.Shorten("https://example.com", ctx)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Find", func(b *testing.B) {
		// First create a URL to find
		url, err := service.Shorten("https://example.com", ctx)
		if err != nil {
			b.Fatal(err)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := service.Find(url.ID, ctx)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("AddBatch", func(b *testing.B) {
		urls := make([]domain.URL, 100)
		for i := range urls {
			urls[i] = domain.URL{
				ID:          "test",
				OriginalURL: "https://example.com",
				UserID:      "user1",
				IsDeleted:   false,
			}
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			err := service.AddBatch(urls, ctx)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("DeleteURLBatch", func(b *testing.B) {
		deleteBatch := model.DeleteBatch{
			UserID:       "user1",
			ShortenedURL: []string{"test1", "test2", "test3"},
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			service.DeleteURLBatch(ctx, deleteBatch)
		}
	})

	resultProfile, err := os.Create("profiles/result.pprof")
	if err != nil {
		log.Fatal(err)
	}
	defer resultProfile.Close()

	if err := pprof.WriteHeapProfile(resultProfile); err != nil {
		log.Fatal(err)
	}
}
