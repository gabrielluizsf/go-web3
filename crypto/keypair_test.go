package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeypairSignVerifySuccess(t *testing.T) {
	privateKey := GeneratePrivateKey()
	publicKey := privateKey.PublicKey()

	msg := []byte("Hello")
	sign, err := privateKey.Sign(msg)
	assert.Nil(t, err)
	assert.NotNil(t, sign)
	assert.True(t, sign.Verify(publicKey, msg))
	assert.NotNil(t, privateKey)
}

func TestKeypairSignVerifyFail(t *testing.T) {
	privateKey := GeneratePrivateKey()
	publicKey := privateKey.PublicKey()
	otherPrivateKey := GeneratePrivateKey()
	otherPublicKey := otherPrivateKey.PublicKey()
	msg := []byte("Hello")
	sign, err := privateKey.Sign(msg)
	assert.Nil(t, err)
	assert.NotNil(t, sign)
	assert.False(t, sign.Verify(otherPublicKey, msg))
	assert.NotNil(t, privateKey)
	assert.False(t, sign.Verify(publicKey,[]byte("Hello World")))
}