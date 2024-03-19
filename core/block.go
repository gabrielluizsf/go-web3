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

func (b *Block) Hash(hasher Hasher[*Block]) types.Hash {
	if b.hash.IsZero() {
		b.hash = hasher.Hash(b)
	}
	return b.hash
}

func (b *Block) Sign(privKey crypto.PrivateKey) error {
	sign, err := privKey.Sign(b.HeaderData())
	if err != nil{
		return err
	}
	b.Validator = privKey.PublicKey()
	b.Signature = sign
	return nil
}

func (b *Block) Verify() error{
	if b.Signature == nil{
		return errors.New("block has no signature")
	}
	if !b.Signature.Verify(b.Validator, b.HeaderData()){
		return errors.New("block has invalid signature")
	}
	return nil
}

func (b *Block) Decode(r io.Reader, dec Decoder[*Block]) error {
	return dec.Decode(r, b)
}

func (b *Block) Encode(w io.Writer, enc Encoder[*Block]) error {
	return enc.Encode(w, b)
}

func (b *Block) HeaderData() []byte{
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	enc.Encode(b.Header)
	return buf.Bytes()
}