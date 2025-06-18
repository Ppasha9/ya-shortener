package storage

import (
	"context"

	"github.com/Ppasha9/ya-shortener/internal/app/model"
)

type Storage interface {
	// Функция для сохранения результата сокращения урла
	SaveURL(ctx context.Context, shortURL, originalURL string) (string, error)

	// Функция для сохранения результата сокращения нескольких урлов
	SaveURLs(ctx context.Context, urls []model.URLsPair) error

	// Функция для получения оригинального урла по сокращенному
	// Если до этого мы не сокращали урл, то вернется ошибка
	GetOriginalURL(ctx context.Context, shortURL string) (string, error)

	// Функция, которая проверяет есть ли уже такой сгенерированный короткий урл в нашей "БД"
	IsExists(ctx context.Context, shortURL string) (bool, error)

	Close() error
}
