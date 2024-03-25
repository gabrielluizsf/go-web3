package core

import (
	"bytes"
	"testing"

	"github.com/gabrielluizsf/go-web3/crypto"
	"github.com/gabrielluizsf/go-web3/types"
	"github.com/stretchr/testify/assert"
)

func TestVerifyTransactionWithTamper(t *testing.T) {
	transaction := NewTransaction(nil)

	fromPrivKey := crypto.GeneratePrivateKey()
	toPrivKey := crypto.GeneratePrivateKey()
	hackerPrivKey := crypto.GeneratePrivateKey()

	transaction.From = fromPrivKey.PublicKey()
	transaction.To = toPrivKey.PublicKey()
	transaction.Value = 666

	assert.Nil(t, transaction.Sign(fromPrivKey))
	transaction.hash = types.Hash{}

	transaction.To = hackerPrivKey.PublicKey()

	assert.NotNil(t, transaction.Verify())
}

func TestNFTTransaction(t *testing.T) {
	collectionTransaction := CollectionTransaction{
		Fee:      200,
		MetaData: []byte("The beginning of a new collection"),
	}

	privKey := crypto.GeneratePrivateKey()
	transaction := &Transaction{
		TransactionInner: collectionTransaction,
	}
	transaction.Sign(privKey)
	transaction.hash = types.Hash{}

	buf := new(bytes.Buffer)
	assert.Nil(t, transaction.Encode(NewGobTransactionEncoder(buf)))

	transactionDecoded := &Transaction{}
	assert.Nil(t, transactionDecoded.Decode(NewGobTransactionDecoder(buf)))
	assert.Equal(t, transaction, transactionDecoded)
}

func TestNativeTransferTransaction(t *testing.T) {
	fromPrivKey := crypto.GeneratePrivateKey()
	toPrivKey := crypto.GeneratePrivateKey()
	transaction := &Transaction{
		To:    toPrivKey.PublicKey(),
		Value: 666,
	}

	assert.Nil(t, transaction.Sign(fromPrivKey))
}

func TestSignTransaction(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	transaction := &Transaction{
		Data: []byte("foo"),
	}

	assert.Nil(t, transaction.Sign(privKey))
	assert.NotNil(t, transaction.Signature)
}

func TestVerifyTransaction(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	transaction := &Transaction{
		Data: []byte("foo"),
	}

	assert.Nil(t, transaction.Sign(privKey))
	assert.Nil(t, transaction.Verify())

	otherPrivKey := crypto.GeneratePrivateKey()
	transaction.From = otherPrivKey.PublicKey()

	assert.NotNil(t, transaction.Verify())
}

func TestTransactionEncodeDecode(t *testing.T) {
	transaction := randomTransactionWithSignature(t)
	buf := &bytes.Buffer{}
	assert.Nil(t, transaction.Encode(NewGobTransactionEncoder(buf)))
	transaction.hash = types.Hash{}

	transactionDecoded := new(Transaction)
	assert.Nil(t, transactionDecoded.Decode(NewGobTransactionDecoder(buf)))
	assert.Equal(t, transaction, transactionDecoded)
}

func randomTransactionWithSignature(t *testing.T) *Transaction {
	privKey := crypto.GeneratePrivateKey()
	transaction := Transaction{
		Data: []byte("foo"),
	}
	assert.Nil(t, transaction.Sign(privKey))

	return &transaction
}
