package core

import (
	"crypto/sha256"

	"github.com/gabrielluizsf/go-web3/types"
)

type Hasher[T any] interface {
	Hash(T) types.Hash
}

type BlockHasher struct{}

func (BlockHasher) Hash(b *Block) types.Hash{
	hash := sha256.Sum256(b.HeaderData())
	return types.Hash(hash)
}
