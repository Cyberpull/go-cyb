package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"

	"cyberpull.com/go-cyb/errors"
)

func EncryptAES(plaintext string, cipherKey ...string) (value string, err error) {
	key, err := GetCryptoKey(cipherKey...)

	if err != nil {
		return
	}

	block, err := aes.NewCipher([]byte(key))

	if err != nil {
		return
	}

	cipherText := make([]byte, aes.BlockSize+len(plaintext))
	iv := cipherText[:aes.BlockSize]

	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], []byte(plaintext))

	value = hex.EncodeToString(cipherText)

	return
}

func DecryptAES(encrypted string, cipherKey ...string) (value string, err error) {
	key, err := GetCryptoKey(cipherKey...)
	if err != nil {
		return
	}

	cipherText, err := hex.DecodeString(encrypted)

	if err != nil {
		return
	}

	block, err := aes.NewCipher([]byte(key))

	if err != nil {
		return
	}

	if len(cipherText) < aes.BlockSize {
		err = errors.New("cipherText block size is too short")
		return
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	value = string(cipherText)

	return
}
