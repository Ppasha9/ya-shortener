package urlshortener

import (
	"math/rand"
	"time"

	"github.com/Ppasha9/ya-shortener/internal/app/storage"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const lenOfGeneratedShortURLs = 8

func generateRandomString() string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	var result []byte
	for i := 0; i < lenOfGeneratedShortURLs; i++ {
		index := seededRand.Intn(len(charset))
		result = append(result, charset[index])
	}

	return string(result)
}

func MakeShortURL(origURL string) string {
	var shortURL string
	for {
		shortURL = generateRandomString()
		if exists := storage.IsExists(shortURL); !exists {
			break
		}
	}
	return shortURL
}
