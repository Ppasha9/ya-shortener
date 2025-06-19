package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"

	serviceerrors "github.com/Ppasha9/ya-shortener/internal/app/errors"
	"github.com/Ppasha9/ya-shortener/internal/app/model"
)

type FileStorageItem struct {
	ShortURL string `json:"short_url"`
	OrigURL  string `json:"orig_url"`
	UserID   uint32 `json:"user_id"`
}

type FileStorage struct {
	Items []FileStorageItem `json:"items"`
}

type InMemoryStorageItem struct {
	OrigURL string
	UserID  uint32
}

type InMemoryStorage struct {
	fileStoragePath string
	urls            map[string]InMemoryStorageItem

	mutex sync.RWMutex
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
		urls:            make(map[string]InMemoryStorageItem),
	}

	if errors.Is(err, io.EOF) {
		return st, nil
	}

	for _, v := range urlsFromFile.Items {
		st.urls[v.ShortURL] = InMemoryStorageItem{OrigURL: v.OrigURL, UserID: v.UserID}
	}

	return st, nil
}

func (d *InMemoryStorage) Clear() {
	d.urls = make(map[string]InMemoryStorageItem)
}

func (d *InMemoryStorage) writeStorageToFile() error {
	fs := FileStorage{
		Items: make([]FileStorageItem, 0),
	}
	for shortURL := range d.urls {
		fs.Items = append(fs.Items, FileStorageItem{ShortURL: shortURL, OrigURL: d.urls[shortURL].OrigURL, UserID: d.urls[shortURL].UserID})
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

func (d *InMemoryStorage) SaveURL(ctx context.Context, userID uint32, shortURL, originalURL string) (string, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// Проверяем, что пытаемся получить укороченный урл того урла, который уже сокращали до этого
	for s, v := range d.urls {
		if v.OrigURL == originalURL {
			return s, serviceerrors.ErrOrigURLDuplicate
		}
	}

	// сохранили в inmemory мапку
	d.urls[shortURL] = InMemoryStorageItem{OrigURL: originalURL, UserID: userID}

	// далее всю мапку сохраняем в файлик
	return shortURL, d.writeStorageToFile()
}

func (d *InMemoryStorage) SaveURLs(ctx context.Context, userID uint32, urls []model.URLsPair) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// сохранили в inmemory мапку
	for _, u := range urls {
		d.urls[u.ShortURL] = InMemoryStorageItem{OrigURL: u.OrigURL, UserID: userID}
	}

	// далее всю мапку сохраняем в файлик
	return d.writeStorageToFile()
}

func (d *InMemoryStorage) GetOriginalURL(ctx context.Context, shortURL string) (string, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	v, ok := d.urls[shortURL]
	if ok {
		return v.OrigURL, nil
	}
	return "", fmt.Errorf("failed to find original url by short url = %s", shortURL)
}

func (d *InMemoryStorage) GetUserURLs(ctx context.Context, userID uint32) ([]model.URLsPair, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	res := make([]model.URLsPair, 0)
	for s, v := range d.urls {
		if v.UserID != userID {
			continue
		}

		res = append(res, model.URLsPair{ShortURL: s, OrigURL: v.OrigURL})
	}

	return res, nil
}

func (d *InMemoryStorage) IsExists(ctx context.Context, shortURL string) (bool, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	_, ok := d.urls[shortURL]
	return ok, nil
}

func (d *InMemoryStorage) Close() error {
	return nil
}
