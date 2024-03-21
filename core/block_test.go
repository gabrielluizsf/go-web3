package core

import (
	"testing"
	"time"

	"github.com/gabrielluizsf/go-web3/crypto"
	"github.com/gabrielluizsf/go-web3/types"
	"github.com/stretchr/testify/assert"
)

func TestHashBlock(t *testing.T) {
	block := randomBlock(0, types.Hash{})
	assert.NotNil(t, block)
	assert.Nil(t, block.Sign(privateKey))
	assert.NotNil(t, block.Signature)
}

func TestVerifyBlock(t *testing.T) {
	block := randomBlock(0, types.Hash{})
	assert.NotNil(t, block)

	assert.Nil(t, block.Sign(privateKey))
	assert.Nil(t, block.Verify())

	otherBlock := randomBlock(0, types.Hash{})
	otherBlock.Height = 100
	assert.NotNil(t, otherBlock.Verify())
}

func TestAddTransaction(t *testing.T){
	b := randomBlockWithSign(t, 0, types.Hash{})
	b.AddTransaction(transaction)
	assert.Equal(t, b.Transactions[1],  *transaction)
}

func randomBlock(height uint32, prevBlockHash types.Hash) *Block {
	header := &Header{
		Version:       1,
		PrevBlockHash: prevBlockHash,
		Height:        height,
		Timestamp:     time.Now().UnixNano(),
	}
	privateKey := crypto.GeneratePrivateKey()
	publicKey := privateKey.PublicKey()
	msg := []byte("Hello World")
	sign, err := privateKey.Sign(msg)
	if err != nil {
		return nil
	}
	transactions := []Transaction{{Data: msg, From: publicKey, Signature: sign}}

	return NewBlock(header, transactions)
}
