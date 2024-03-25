package util

import (
	"crypto/rand"
	"testing"
	"time"

	"github.com/gabrielluizsf/go-web3/core"
	"github.com/gabrielluizsf/go-web3/crypto"
	"github.com/gabrielluizsf/go-web3/types"
	"github.com/stretchr/testify/assert"
)

func RandomBytes(size int) []byte {
	token := make([]byte, size)
	rand.Read(token)
	return token
}

func RandomHash() types.Hash {
	return types.HashFromBytes(RandomBytes(32))
}

// NewRandomTransaction return a new random transaction whithout signature.
func NewRandomTransaction(size int) *core.Transaction {
	return core.NewTransaction(RandomBytes(size))
}

func NewRandomTransactionWithSignature(t *testing.T, privKey crypto.PrivateKey, size int) *core.Transaction {
	transaction := NewRandomTransaction(size)
	assert.Nil(t, transaction.Sign(privKey))
	return transaction
}

func NewRandomBlock(t *testing.T, height uint32, prevBlockHash types.Hash) *core.Block {
	transactionSigner := crypto.GeneratePrivateKey()
	transaction := NewRandomTransactionWithSignature(t, transactionSigner, 100)
	header := &core.Header{
		Version:       1,
		PrevBlockHash: prevBlockHash,
		Height:        height,
		Timestamp:     time.Now().UnixNano(),
	}
	b, err := core.NewBlock(header, []*core.Transaction{transaction})
	assert.Nil(t, err)
	dataHash, err := core.CalculateDataHash(b.Transactions)
	assert.Nil(t, err)
	b.Header.DataHash = dataHash

	return b
}

func NewRandomBlockWithSignature(t *testing.T, pk crypto.PrivateKey, height uint32, prevHash types.Hash) *core.Block {
	b := NewRandomBlock(t, height, prevHash)
	assert.Nil(t, b.Sign(pk))

	return b
}
