package core

import (
	"encoding/gob"
	"fmt"
	"math/rand"

	"github.com/gabrielluizsf/go-web3/crypto"
	"github.com/gabrielluizsf/go-web3/types"
)

type TransactionType byte

const (
	TransactionTypeCollection TransactionType = iota // 0x0
	TransactionTypeMint                              // 0x01
)

type CollectionTransaction struct {
	Fee      int64
	MetaData []byte
}

type MintTransaction struct {
	Fee             int64
	NFT             types.Hash
	Collection      types.Hash
	MetaData        []byte
	CollectionOwner crypto.PublicKey
	Signature       crypto.Signature
}

type Transaction struct {
	// Only used for native NFT logic
	TransactionInner any
	// Any arbitrary data for the VM
	Data      []byte
	To        crypto.PublicKey
	Value     uint64
	From      crypto.PublicKey
	Signature *crypto.Signature
	Nonce     int64

	// cached version of the Transaction data hash
	hash types.Hash
}

func NewTransaction(data []byte) *Transaction {
	return &Transaction{
		Data:  data,
		Nonce: rand.Int63n(1000000000000000),
	}
}

func (transaction *Transaction) Hash(hasher Hasher[*Transaction]) types.Hash {
	if transaction.hash.IsZero() {
		transaction.hash = hasher.Hash(transaction)
	}
	return transaction.hash
}

func (transaction *Transaction) Sign(privKey crypto.PrivateKey) error {
	hash := transaction.Hash(TransactionHasher{})
	sig, err := privKey.Sign(hash.ToSlice())
	if err != nil {
		return err
	}

	transaction.From = privKey.PublicKey()
	transaction.Signature = sig

	return nil
}

func (transaction *Transaction) Verify() error {
	if transaction.Signature == nil {
		return fmt.Errorf("transaction has no signature")
	}

	hash := transaction.Hash(TransactionHasher{})
	if !transaction.Signature.Verify(transaction.From, hash.ToSlice()) {
		return fmt.Errorf("invalid transaction signature")
	}

	return nil
}

func (transaction *Transaction) Decode(dec Decoder[*Transaction]) error {
	return dec.Decode(transaction)
}

func (transaction *Transaction) Encode(enc Encoder[*Transaction]) error {
	return enc.Encode(transaction)
}

func init() {
	gob.Register(CollectionTransaction{})
	gob.Register(MintTransaction{})
}
