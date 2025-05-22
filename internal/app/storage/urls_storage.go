package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
)

var StorageMutex sync.RWMutex

type FileStorageItem struct {
	ShortURL string `json:"short_url"`
	OrigURL  string `json:"orig_url"`
}

type FileStorage struct {
	Items []FileStorageItem `json:"items"`
}

type InMemoryStorage struct {
	fileStoragePath string
	urls            map[string]string
}

func NewInMemoryStorage(fileStoragePath string) (*InMemoryStorage, error) {
	file, err := os.OpenFile(fileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var urlsFromFile FileStorage
	err = decoder.Decode(&urlsFromFile)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}

	st := &InMemoryStorage{
		fileStoragePath: fileStoragePath,
		urls:            make(map[string]string),
	}

	if errors.Is(err, io.EOF) {
		return st, nil
	}

	for _, v := range urlsFromFile.Items {
		st.urls[v.ShortURL] = v.OrigURL
	}

	return st, nil
}

func (d *InMemoryStorage) Clear() {
	d.urls = make(map[string]string)
}

// Функция для сохранения результата сокращения урла
func (d *InMemoryStorage) SaveURL(shortURL, originalURL string) error {
	StorageMutex.Lock()

	// сохранили в inmemory мапку
	d.urls[shortURL] = originalURL

	// далее всю мапку сохраняем в файлик
	fs := FileStorage{
		Items: make([]FileStorageItem, 0),
	}
	for k, v := range d.urls {
		fs.Items = append(fs.Items, FileStorageItem{ShortURL: k, OrigURL: v})
	}
	data, err := json.MarshalIndent(&fs, "", "   ")
	if err != nil {
		return err
	}
	err = os.WriteFile(d.fileStoragePath, data, 0666)
	if err != nil {
		return err
	}

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
