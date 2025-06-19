package auth

import (
	"crypto/rand"
	"encoding/binary"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
)

var keyLen = 8
var tokenTTL = time.Hour

type claims struct {
	jwt.RegisteredClaims
	UserID uint32
}

type Auth struct {
	// ключ для симметричного шифрования куки с айдишкой пользователя
	// будет генерится случайно при каждом запуске сервиса
	// в реальном приложении надо хранить в надежном месте, здесь генерим
	key []byte

	// тоже обычно надо брать из конфига, например
	// в данном приложении будет захардкожено
	TokenTTL time.Duration
}

func generateRandomByteSlice(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func NewAuth() (*Auth, error) {
	var err error
	res := &Auth{}

	res.key, err = generateRandomByteSlice(keyLen)
	if err != nil {
		return nil, err
	}
	res.TokenTTL = tokenTTL

	return res, nil
}

func GenerateUserID() (uint32, error) {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		return 0, errors.Wrap(err, "failed to generate random token id")
	}
	return binary.BigEndian.Uint32(b), nil
}

func (a *Auth) NewJWT(userID uint32) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.TokenTTL)),
		},
		UserID: userID,
	})
	return token.SignedString(a.key)
}

func (a *Auth) ParseUserIDFromJWT(tokenStr string) (uint32, error) {
	claims := &claims{}
	// парсим из строки токена tokenString в структуру claims
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return a.key, nil
	})

	if err != nil || !token.Valid {
		return 0, errors.Wrap(err, "cannot parse JWT token")
	}

	return claims.UserID, nil
}
