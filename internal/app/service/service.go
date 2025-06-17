package service

import (
	"context"

	"github.com/Ppasha9/ya-shortener/internal/app/storage"
	"github.com/Ppasha9/ya-shortener/internal/app/urlshortener"
)

type Service struct {
	Storage storage.Storage
}

func NewService(s storage.Storage) *Service {
	return &Service{
		Storage: s,
	}
}

func (s *Service) MakeShortURL(ctx context.Context, origURL string) (string, error) {
	var shortURL string
	for {
		shortURL = urlshortener.MakeShortURL(origURL)
		exists, err := s.Storage.IsExists(ctx, shortURL)
		if err != nil {
			return shortURL, err
		}
		if !exists {
			break
		}
	}
	err := s.Storage.SaveURL(ctx, shortURL, origURL)
	return shortURL, err
}

func (s *Service) GetOriginalURL(ctx context.Context, shortURL string) (string, error) {
	return s.Storage.GetOriginalURL(ctx, shortURL)
}
