package core

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type Blockchain struct {
	store     Storage
	headers   []*Header
	validator Validator
}

func (bc *Blockchain) SetValidator(v Validator) {
	bc.validator = v
}

func NewBlockchain(genesis *Block, store Storage) (*Blockchain, error) {
	bc := &Blockchain{
		headers: []*Header{},
		store:   store,
	}
	bc.SetValidator(NewBlockValidator(bc))
	err := bc.addBlockWithoutValidation(genesis)
	return bc, err
}

func (bc *Blockchain) AddBlock(b *Block) error {
	if err := bc.validator.ValidateBlock(b); err != nil {
		return err
	}

	return bc.addBlockWithoutValidation(b)
}

func (bc *Blockchain) addBlockWithoutValidation(b *Block) error {
	logrus.WithFields(logrus.Fields{
		"height": b.Height,
		"hash": b.Hash(BlockHasher{}),
	}).Info("adding new block")
	bc.headers = append(bc.headers, b.Header)
	return bc.store.Put(b)
}

func (bc *Blockchain) HasBlock(height uint32) bool {
	return height <= bc.Height()
}

func (bc *Blockchain) GetHeader(height uint32) (*Header, error) {
	if height > bc.Height() {
		return nil, fmt.Errorf("given height (%d) too high", height)
	}
	return bc.headers[height], nil
}

func (bc *Blockchain) Height() uint32 {
	return uint32(len(bc.headers) - 1)
}