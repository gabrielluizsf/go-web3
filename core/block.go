package core

import (
	"bytes"
	"encoding/gob"
	"errors"
	"io"

	"github.com/gabrielluizsf/go-web3/crypto"
	"github.com/gabrielluizsf/go-web3/types"
)

type Header struct {
	Version       uint32
	DataHash      types.Hash
	PrevBlockHash types.Hash
	Timestamp     int64
	Height        uint32
}

func (h *Header) Bytes() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	enc.Encode(h)
	return buf.Bytes()
}

type Block struct {
	*Header
	Transactions []Transaction
	Validator    crypto.PublicKey
	Signature    *crypto.Signature
	hash         types.Hash
}

func NewBlock(h *Header, transactions []Transaction) *Block {
	return &Block{
		Header:       h,
		Transactions: transactions,
	}
}

func (b *Block) AddTransaction(transaction *Transaction){
	b.Transactions = append(b.Transactions, *transaction)
}

func (b *Block) Hash(hasher Hasher[*Header]) types.Hash {
	if b.hash.IsZero() {
		b.hash = hasher.Hash(b.Header)
	}
	return b.hash
}

func (b *Block) Sign(privKey crypto.PrivateKey) error {
	sign, err := privKey.Sign(b.Header.Bytes())
	if err != nil {
		return err
	}
	b.Validator = privKey.PublicKey()
	b.Signature = sign
	return nil
}

func (b *Block) Verify() error {
	if b.Signature == nil {
		return errors.New("block has no signature")
	}
	if !b.Signature.Verify(b.Validator, b.Header.Bytes()) {
		return errors.New("block has invalid signature")
	}
	for _, transaction := range b.Transactions{
		if err := transaction.Verify(); err != nil{
			return err
		}
	}
	return nil
}

func (b *Block) Decode(r io.Reader, dec Decoder[*Block]) error {
	return dec.Decode(r, b)
}

func (b *Block) Encode(w io.Writer, enc Encoder[*Block]) error {
	return enc.Encode(w, b)
}
