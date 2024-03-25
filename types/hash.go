package types

import (
	"encoding/hex"
	"fmt"
)

const HASH_LENGHT = 32

type Hash [HASH_LENGHT]uint8

func (h Hash) IsZero() bool {
	for i := 0; i < HASH_LENGHT; i++ {
		if h[i] != 0 {
			return false
		}
	}
	return true
}

func (h Hash) ToSlice() []byte {
	b := make([]byte, HASH_LENGHT)
	for i := 0; i < HASH_LENGHT; i++ {
		b[i] = h[i]
	}
	return b
}

func (h Hash) String() string {
	return hex.EncodeToString(h.ToSlice())
}

func HashFromBytes(b []byte) Hash {
	if len(b) != HASH_LENGHT {
		msg := fmt.Sprintf("given bytes with length %d should be 32", len(b))
		panic(msg)
	}

	var value [HASH_LENGHT]uint8
	for i := 0; i < HASH_LENGHT; i++ {
		value[i] = b[i]
	}

	return Hash(value)
}
