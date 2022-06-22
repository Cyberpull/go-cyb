package tests

import (
	"testing"

	_ "cyberpull.com/go-cyb/env"

	"cyberpull.com/go-cyb/crypto"
)

const CryptoKey string = "804w8irdjfry9573l2348jd0"

var (
	cryptoPlainText string = "Testing Data"
	cryptoEncrypted string
)

func TestCrypto_Encrypt(t *testing.T) {
	var err error

	cryptoEncrypted, err = crypto.EncryptAES(cryptoPlainText, CryptoKey)

	if err != nil {
		t.Fatal(err)
	}
}

func TestCrypto_Decrypt(t *testing.T) {
	var err error

	value, err := crypto.DecryptAES(cryptoEncrypted, CryptoKey)

	if err != nil {
		t.Fatal(err)
		return
	}

	if value != cryptoPlainText {
		t.Fatalf(`Expected "%s", got "%s" instead`, cryptoPlainText, value)
	}
}
