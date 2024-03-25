package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"

	"github.com/gabrielluizsf/go-web3/types"
)

type Hasher[T any] interface {
	Hash(T) types.Hash
}

type BlockHasher struct{}

func (BlockHasher) Hash(b *Header) types.Hash {
	h := sha256.Sum256(b.Bytes())
	return types.Hash(h)
}

type TransactionHasher struct{}

// Hash will hash the whole bytes of the transaction no exception.
func (TransactionHasher) Hash(transaction *Transaction) types.Hash {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, transaction.Data)
	binary.Write(buf, binary.LittleEndian, transaction.To)
	binary.Write(buf, binary.LittleEndian, transaction.Value)
	binary.Write(buf, binary.LittleEndian, transaction.From)
	binary.Write(buf, binary.LittleEndian, transaction.Nonce)

	return types.Hash(sha256.Sum256(buf.Bytes()))
}
