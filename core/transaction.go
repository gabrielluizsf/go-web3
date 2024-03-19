package core

import (
	"errors"

	"github.com/gabrielluizsf/go-web3/crypto"
)

type Transaction struct {
	Data []byte

	PublicKey crypto.PublicKey
	Signature *crypto.Signature
}

func NewTransaction(data []byte) *Transaction {
	return &Transaction{
		Data: data,
	}
}

func (t *Transaction) Sign(privKey crypto.PrivateKey) error {
	sign, err := privKey.Sign(t.Data)
	if err != nil {
		return err
	}
	t.PublicKey = privKey.PublicKey()
	t.Signature = sign
	return nil
}

func (t *Transaction) Verify() error {
	if t.Signature == nil {
		return errors.New("transaction has no signature")
	}
	if !t.Signature.Verify(t.PublicKey, t.Data) {
		return errors.New("invalid transaction signature")
	}
	return nil
}
