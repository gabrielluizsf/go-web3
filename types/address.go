package types

import (
	"encoding/hex"
	"fmt"
)

const ADDRESS_MAX_LENGHT = 20

type Address [ADDRESS_MAX_LENGHT]uint8

func (a Address) Slice() []byte {
	slice := make([]byte, ADDRESS_MAX_LENGHT)
	for i := 0; i < ADDRESS_MAX_LENGHT; i++ {
		slice[i] = a[i]
	}
	return slice
}

func (a Address) String() string {
	return hex.EncodeToString(a.Slice())
}

func AddressFromBytes(b []byte) Address {
	bytesLenght := len(b)
	if bytesLenght != ADDRESS_MAX_LENGHT {
		msg := fmt.Sprintf("Given bytes with lenght %d should be %d", bytesLenght, ADDRESS_MAX_LENGHT)
		panic(msg)
	}

	var addr Address
	for i := 0; i < ADDRESS_MAX_LENGHT; i++ {
		addr[i] = b[i]
	}
	return addr
}
