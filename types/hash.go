package types

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

const HASH_LENGHT = 32

type Hash [HASH_LENGHT]uint8

func (h Hash) Slice() []byte {
	slice := make([]byte, HASH_LENGHT)
	for i := 0; i < HASH_LENGHT; i++ {
		slice[i] = h[i]
	}
	return slice
}

func (h Hash) String() string {
	return hex.EncodeToString(h.Slice())
}

func (h Hash) IsZero() bool {
	for i := 0; i < HASH_LENGHT; i++ {
		if h[i] != 0 {
			return false
		}
	}
	return true
}

func HashFromBytes(b []byte) Hash {
	bytesLenght := len(b)
	if bytesLenght != HASH_LENGHT {
		msg := fmt.Sprintf("Given bytes with lenght %d should be %d", bytesLenght, HASH_LENGHT)
		panic(msg)
	}

	var value Hash
	for i := 0; i < HASH_LENGHT; i++ {
		value[i] = b[i]
	}
	return value
}

func RandomBytes(size int) []byte {
	token := make([]byte, size)
	rand.Read(token)
	return token
}

func RandomHash() Hash {
	return HashFromBytes(RandomBytes(HASH_LENGHT))
}
