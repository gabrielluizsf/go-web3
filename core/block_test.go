package core

import (
	"testing"
	"time"

	"github.com/gabrielluizsf/go-web3/crypto"
	"github.com/gabrielluizsf/go-web3/types"
	"github.com/stretchr/testify/assert"
)

func randomBlock(height uint32) *Block {
	header := &Header{
		Version:       1,
		PrevBlockHash: types.RandomHash(),
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
	transactions := []Transaction{{Data: msg, PublicKey: publicKey, Signature: sign}}

	return NewBlock(header, transactions)
}

func TestHashBlock(t *testing.T) {
	block := randomBlock(0)
	assert.NotNil(t, block)
	assert.Nil(t, block.Sign(privateKey))
	assert.NotNil(t, block.Signature)
}

func TestVerifyBlock(t *testing.T){
	block := randomBlock(0)
	assert.NotNil(t, block)

	assert.Nil(t, block.Sign(privateKey))
	assert.Nil(t, block.Verify())

	otherBlock := randomBlock(0)
	otherBlock.Height = 100
	assert.NotNil(t, otherBlock.Verify())
}