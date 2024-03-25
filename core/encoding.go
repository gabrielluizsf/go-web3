package core

import (
	"encoding/gob"
	"io"
)

//
// For now we GOB encoding is used for fast bootstrapping of the project
// in a later phase I'm considering using Protobuffers as default encoding / decoding.
//

type Encoder[T any] interface {
	Encode(T) error
}

type Decoder[T any] interface {
	Decode(T) error
}

type GobTransactionEncoder struct {
	w io.Writer
}

func NewGobTransactionEncoder(w io.Writer) *GobTransactionEncoder {
	return &GobTransactionEncoder{
		w: w,
	}
}

func (e *GobTransactionEncoder) Encode(transaction *Transaction) error {
	return gob.NewEncoder(e.w).Encode(transaction)
}

type GobTransactionDecoder struct {
	r io.Reader
}

func NewGobTransactionDecoder(r io.Reader) *GobTransactionDecoder {
	return &GobTransactionDecoder{
		r: r,
	}
}

func (e *GobTransactionDecoder) Decode(transaction *Transaction) error {
	return gob.NewDecoder(e.r).Decode(transaction)
}

type GobBlockEncoder struct {
	w io.Writer
}

func NewGobBlockEncoder(w io.Writer) *GobBlockEncoder {
	return &GobBlockEncoder{
		w: w,
	}
}

func (enc *GobBlockEncoder) Encode(b *Block) error {
	return gob.NewEncoder(enc.w).Encode(b)
}

type GobBlockDecoder struct {
	r io.Reader
}

func NewGobBlockDecoder(r io.Reader) *GobBlockDecoder {
	return &GobBlockDecoder{
		r: r,
	}
}

func (dec *GobBlockDecoder) Decode(b *Block) error {
	return gob.NewDecoder(dec.r).Decode(b)
}
