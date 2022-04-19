package crypto

import (
	"os"

	_ "cyberpull.com/go-cyb/env"

	"cyberpull.com/go-cyb/errors"
)

func GetCryptoKey(key ...string) (value string, err error) {
	if len(key) > 0 {
		value = key[0]
		return
	}

	value = os.Getenv("CRYPTO_KEY")

	if value == "" {
		err = errors.New(`"CRYPTO_KEY" is required`)
	}

	return
}
