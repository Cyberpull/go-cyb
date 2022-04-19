package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func Hash(data string) (value string, err error) {
	hasher := sha256.New()

	_, err = hasher.Write([]byte(data))
	if err != nil {
		return
	}

	value = hex.EncodeToString(hasher.Sum(nil))

	return
}

func Hmac(data string, cipherKey ...string) (value string, err error) {
	key, err := GetCryptoKey(cipherKey...)
	if err != nil {
		return
	}

	hasher := hmac.New(sha256.New, []byte(key))

	_, err = hasher.Write([]byte(data))
	if err != nil {
		return
	}

	value = hex.EncodeToString(hasher.Sum(nil))

	return
}
