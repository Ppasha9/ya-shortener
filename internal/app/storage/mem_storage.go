package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/Ppasha9/ya-shortener/internal/app/model"
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

func (d *InMemoryStorage) writeStorageToFile() error {
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
	return nil
}

func (d *InMemoryStorage) SaveURL(ctx context.Context, shortURL, originalURL string) error {
	StorageMutex.Lock()
	defer StorageMutex.Unlock()

	// сохранили в inmemory мапку
	d.urls[shortURL] = originalURL

	// далее всю мапку сохраняем в файлик
	return d.writeStorageToFile()
}

func (d *InMemoryStorage) SaveURLs(ctx context.Context, urls []model.URLsPair) error {
	StorageMutex.Lock()
	defer StorageMutex.Unlock()

	// сохранили в inmemory мапку
	for _, u := range urls {
		d.urls[u.ShortURL] = u.OrigURL
	}

	// далее всю мапку сохраняем в файлик
	return d.writeStorageToFile()
}

func (d *InMemoryStorage) GetOriginalURL(ctx context.Context, shortURL string) (string, error) {
	StorageMutex.Lock()
	origURL, ok := d.urls[shortURL]
	StorageMutex.Unlock()
	if ok {
		return origURL, nil
	}

	return origURL, fmt.Errorf("failed to find original url by short url = %s", shortURL)
}

func (d *InMemoryStorage) IsExists(ctx context.Context, shortURL string) (bool, error) {
	StorageMutex.Lock()
	_, ok := d.urls[shortURL]
	StorageMutex.Unlock()
	return ok, nil
}

func (d *InMemoryStorage) Close() error {
	return nil
}
