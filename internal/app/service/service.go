package service

import (
	"context"

	"github.com/Ppasha9/ya-shortener/internal/app/model"
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

func (s *Service) MakeShortURL(ctx context.Context, userID uint32, origURL string) (string, error) {
	var shortURL string
	for {
		shortURL = urlshortener.MakeShortURL(origURL)
		exists, err := s.Storage.IsExists(ctx, userID, shortURL)
		if err != nil {
			return shortURL, err
		}
		if !exists {
			break
		}
	}
	shortURL, err := s.Storage.SaveURL(ctx, userID, shortURL, origURL)
	return shortURL, err
}

func (s *Service) GetOriginalURL(ctx context.Context, userID uint32, shortURL string) (string, error) {
	return s.Storage.GetOriginalURL(ctx, userID, shortURL)
}

func (s *Service) MakeShortURLsBatch(ctx context.Context, userID uint32, batch []model.ShortenBatchItemRequest) ([]model.ShortenBatchItemResponse, error) {
	urls := make([]model.URLsPair, 0)
	resp := make([]model.ShortenBatchItemResponse, 0)
	for _, v := range batch {
		var shortURL string
		for {
			shortURL = urlshortener.MakeShortURL(v.OrigURL)
			exists, err := s.Storage.IsExists(ctx, userID, shortURL)
			if err != nil {
				return nil, err
			}
			if !exists {
				break
			}
		}

		urls = append(urls, model.URLsPair{ShortURL: shortURL, OrigURL: v.OrigURL})
		resp = append(resp, model.ShortenBatchItemResponse{CorrID: v.CorrID, ShortURL: shortURL})
	}

	err := s.Storage.SaveURLs(ctx, userID, urls)
	return resp, err
}
