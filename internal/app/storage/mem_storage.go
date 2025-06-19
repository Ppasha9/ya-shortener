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

type InMemoryStorage struct {
	fileStoragePath string
	urls            map[uint32]map[string]string

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
		urls:            make(map[uint32]map[string]string),
	}

	if errors.Is(err, io.EOF) {
		return st, nil
	}

	for _, v := range urlsFromFile.Items {
		_, ok := st.urls[v.UserID]
		if !ok {
			st.urls[v.UserID] = make(map[string]string)
		}

		st.urls[v.UserID][v.ShortURL] = v.OrigURL
	}

	return st, nil
}

func (d *InMemoryStorage) Clear() {
	d.urls = make(map[uint32]map[string]string)
}

func (d *InMemoryStorage) writeStorageToFile() error {
	fs := FileStorage{
		Items: make([]FileStorageItem, 0),
	}
	for userID := range d.urls {
		for short, orig := range d.urls[userID] {
			fs.Items = append(fs.Items, FileStorageItem{ShortURL: short, OrigURL: orig, UserID: userID})
		}
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
	for s, o := range d.urls[userID] {
		if o == originalURL {
			return s, serviceerrors.ErrOrigURLDuplicate
		}
	}

	// сохранили в inmemory мапку
	_, ok := d.urls[userID]
	if !ok {
		d.urls[userID] = make(map[string]string)
	}
	d.urls[userID][shortURL] = originalURL

	// далее всю мапку сохраняем в файлик
	return shortURL, d.writeStorageToFile()
}

func (d *InMemoryStorage) SaveURLs(ctx context.Context, userID uint32, urls []model.URLsPair) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// сохранили в inmemory мапку
	_, ok := d.urls[userID]
	if !ok {
		d.urls[userID] = make(map[string]string)
	}
	for _, u := range urls {
		d.urls[userID][u.ShortURL] = u.OrigURL
	}

	// далее всю мапку сохраняем в файлик
	return d.writeStorageToFile()
}

func (d *InMemoryStorage) GetOriginalURL(ctx context.Context, userID uint32, shortURL string) (string, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	_, ok := d.urls[userID]
	if !ok {
		return "", fmt.Errorf("failed to find original url by short url = %s", shortURL)
	}
	origURL, ok := d.urls[userID][shortURL]
	if ok {
		return origURL, nil
	}
	return origURL, fmt.Errorf("failed to find original url by short url = %s", shortURL)
}

func (d *InMemoryStorage) IsExists(ctx context.Context, userID uint32, shortURL string) (bool, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	_, ok := d.urls[userID]
	if !ok {
		return false, nil
	}
	_, ok = d.urls[userID][shortURL]
	return ok, nil
}

func (d *InMemoryStorage) Close() error {
	return nil
}
