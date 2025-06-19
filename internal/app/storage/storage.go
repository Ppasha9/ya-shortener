package storage

import (
	"context"

	"github.com/Ppasha9/ya-shortener/internal/app/model"
)

type Storage interface {
	// Функция для сохранения результата сокращения урла
	SaveURL(ctx context.Context, userID uint32, shortURL, originalURL string) (string, error)

	// Функция для сохранения результата сокращения нескольких урлов
	SaveURLs(ctx context.Context, userID uint32, urls []model.URLsPair) error

	// Функция для получения оригинального урла по сокращенному
	// Если до этого мы не сокращали урл, то вернется ошибка
	GetOriginalURL(ctx context.Context, shortURL string) (string, error)

	// Функция, которая возвращает все пары урлов, которые определенный юзера сокращал до этого
	// В случае, когда юзер еще никакие урыл не сокращал, должен вернуться пустой слайс без ошибки
	GetUserURLs(ctx context.Context, userID uint32) ([]model.URLsPair, error)

	// Функция, которая проверяет есть ли уже такой сгенерированный короткий урл в нашей "БД"
	IsExists(ctx context.Context, shortURL string) (bool, error)

	Close() error
}
