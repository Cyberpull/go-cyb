package crypto

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"math"
)

func RandomBytes(max int) (value []byte, err error) {
	value = make([]byte, max)

	if _, err = rand.Read(value); err != nil {
		return
	}

	return
}

func RandomInt(max int) (value int, err error) {
	b, err := RandomBytes(max)

	if err != nil {
		return
	}

	value = int(binary.BigEndian.Uint64(b))

	return
}

func RandomString(max int) (value string, err error) {
	max = int(math.Ceil(float64(max) / 2))
	b, err := RandomBytes(max)

	if err != nil {
		return
	}

	value = hex.EncodeToString(b)

	return
}
