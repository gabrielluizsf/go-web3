package core

import (
	"testing"

	"github.com/gabrielluizsf/go-web3/crypto"
	"github.com/stretchr/testify/assert"
)

var (
	privateKey  = crypto.GeneratePrivateKey()
	data        = []byte("Hello World")
	transaction = NewTransaction(data)
)

func TestSignTransaction(t *testing.T) {
	assert.Nil(t, transaction.Sign(privateKey))
	assert.NotNil(t, transaction.Signature)
	assert.Equal(t, transaction.From, privateKey.PublicKey())
}

func TestVerifyTransaction(t *testing.T) {
	assert.Nil(t, transaction.Sign(privateKey))
	assert.Nil(t, transaction.Verify())
	otherPrivKey := crypto.GeneratePrivateKey()
	transaction.From = otherPrivKey.PublicKey()

	assert.NotNil(t, transaction.Verify())
}
