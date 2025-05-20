package storage

import (
	"fmt"
	"sync"
)

var StorageMutex sync.RWMutex

type InMemoryStorage struct {
	urls map[string]string
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		urls: make(map[string]string),
	}
}

func (d *InMemoryStorage) Clear() {
	d.urls = make(map[string]string)
}

// Функция для сохранения результата сокращения урла
func (d *InMemoryStorage) SaveURL(shortURL, originalURL string) error {
	StorageMutex.Lock()
	d.urls[shortURL] = originalURL
	StorageMutex.Unlock()
	return nil
}

// Функция для получения оригинального урла по сокращенному
// Если до этого мы не сокращали урл, то вернется ошибка
func (d *InMemoryStorage) GetOriginalURL(shortURL string) (string, error) {
	origURL, ok := d.urls[shortURL]
	if ok {
		return origURL, nil
	}

	return origURL, fmt.Errorf("failed to find original url by short url = %s", shortURL)
}

// Функция, которая проверяет есть ли уже такой сгенерированный короткий урл в нашей "БД"
func (d *InMemoryStorage) IsExists(shortURL string) bool {
	_, ok := d.urls[shortURL]
	return ok
}
