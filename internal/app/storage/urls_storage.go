package storage

import (
	"fmt"
)

// Глобальная переменная - хранилище сокращенных урлов
var urlsStorage map[string]string

func Init() {
	urlsStorage = make(map[string]string)
}

// Функция для сохранения результата сокращения урла
func SaveURL(shortURL, originalURL string) {
	urlsStorage[shortURL] = originalURL
}

// Функция для получения оригинального урла по сокращенному
// Если до этого мы не сокращали урл, то вернется ошибка
func GetOriginalURL(shortURL string) (string, error) {
	origURL, ok := urlsStorage[shortURL]
	if ok {
		return origURL, nil
	}

	return origURL, fmt.Errorf("failed to find original url by short url = %s", shortURL)
}

// Функция, которая проверяет есть ли уже такой сгенерированный короткий урл в нашей "БД"
func IsExists(shortURL string) bool {
	_, ok := urlsStorage[shortURL]
	return ok
}
