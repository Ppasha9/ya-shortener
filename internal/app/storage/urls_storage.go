package storage

import (
	"fmt"
)

type Database struct {
	urls map[string]string
}

func NewDatabase() *Database {
	return &Database{
		urls: make(map[string]string),
	}
}

func (d *Database) Clear() {
	d.urls = make(map[string]string)
}

// Функция для сохранения результата сокращения урла
func (d *Database) SaveURL(shortURL, originalURL string) {
	d.urls[shortURL] = originalURL
}

// Функция для получения оригинального урла по сокращенному
// Если до этого мы не сокращали урл, то вернется ошибка
func (d *Database) GetOriginalURL(shortURL string) (string, error) {
	origURL, ok := d.urls[shortURL]
	if ok {
		return origURL, nil
	}

	return origURL, fmt.Errorf("failed to find original url by short url = %s", shortURL)
}

// Функция, которая проверяет есть ли уже такой сгенерированный короткий урл в нашей "БД"
func (d *Database) IsExists(shortURL string) bool {
	_, ok := d.urls[shortURL]
	return ok
}
