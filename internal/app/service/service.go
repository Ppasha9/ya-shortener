package service

import (
	"github.com/Ppasha9/ya-shortener/internal/app/storage"
	"github.com/Ppasha9/ya-shortener/internal/app/urlshortener"
)

type Service struct {
	Storage *storage.InMemoryStorage
}

func NewService(s *storage.InMemoryStorage) *Service {
	return &Service{
		Storage: s,
	}
}

func (s *Service) MakeShortURL(origURL string) (string, error) {
	var shortURL string
	for {
		shortURL = urlshortener.MakeShortURL(origURL)
		if exists := s.Storage.IsExists(shortURL); !exists {
			break
		}
	}
	err := s.Storage.SaveURL(shortURL, origURL)
	return shortURL, err
}

func (s *Service) GetOriginalURL(shortURL string) (string, error) {
	return s.Storage.GetOriginalURL(shortURL)
}
