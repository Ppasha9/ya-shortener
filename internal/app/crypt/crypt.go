package crypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"

	"github.com/pkg/errors"

	serviceerrors "github.com/Ppasha9/ya-shortener/internal/app/errors"
)

var saltSize = 12

type Crypt struct {
	// ключ для симметричного шифрования куки с айдишкой пользователя
	// будет генерится случайно при каждом запуске сервиса
	key []byte

	// это подпись, которая будет добавляться в данные, необходимые для шифровки
	// будет генерится рандомно при каждом запуске сервиса
	// она нужна для того, чтобы проверять подлинность симметрично подписанной куки
	salt []byte

	// объект для шифрования и дешифрования данных
	cryptCipher cipher.Block
}

func generateRandomByteSlice(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func NewCrypt() (*Crypt, error) {
	var err error
	res := &Crypt{}

	res.key, err = generateRandomByteSlice(aes.BlockSize)
	if err != nil {
		return nil, err
	}

	res.salt, err = generateRandomByteSlice(saltSize)
	if err != nil {
		return nil, err
	}

	res.cryptCipher, err = aes.NewCipher(res.key)
	if err != nil {
		return nil, err
	}

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

func (c *Crypt) GenerateAuthCookie(userID uint32) string {
	idBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(idBytes, userID)

	bytesToEncrypt := append(c.salt, idBytes...)
	encrypted := make([]byte, aes.BlockSize)
	c.cryptCipher.Encrypt(encrypted, bytesToEncrypt)

	return hex.EncodeToString(encrypted)
}

func (c *Crypt) GetUserIDFromAuthCookie(cookieVal string) (uint32, error) {
	encryptedBytes, err := hex.DecodeString(cookieVal)
	if err != nil {
		return 0, err
	}

	decryptedBytes := make([]byte, aes.BlockSize)
	c.cryptCipher.Decrypt(decryptedBytes, encryptedBytes)

	if !bytes.Equal(decryptedBytes[:saltSize], c.salt) {
		return 0, serviceerrors.ErrInvalidAuthCookie
	}

	if len(decryptedBytes) != saltSize+4 {
		return 0, serviceerrors.ErrInvalidAuthCookieBytesLen
	}

	return binary.BigEndian.Uint32(decryptedBytes[saltSize:]), nil
}
