package core

import (
	"testing"

	"github.com/gabrielluizsf/go-web3/crypto"
	"github.com/gabrielluizsf/go-web3/types"
	"github.com/stretchr/testify/assert"
)

var (
	genesisBlockHeight = uint32(0)
)

func TestNewBlockchain(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	assert.Equal(t, bc.Height(), genesisBlockHeight)
	assert.NotNil(t, bc.validator)
}

func TestHasBlock(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	assert.True(t, bc.HasBlock(genesisBlockHeight))
	assert.False(t, bc.HasBlock(genesisBlockHeight+100))
	assert.False(t, bc.HasBlock(genesisBlockHeight+50*2))
}

func TestAddBlock(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	lenBlocks := 1000
	for i := 0; i < lenBlocks; i++ {
		h := uint32(i + 1)
		block := randomBlockWithSign(t, h, getPrevBlockHash(t, bc, h))
		assert.Nil(t, bc.AddBlock(block))
	}
	assert.Equal(t, bc.Height(), uint32(lenBlocks))
	assert.Equal(t, len(bc.headers), lenBlocks+1)

	assert.NotNil(t, bc.AddBlock(randomBlock(98, types.Hash{})))
}

func TestGetHeader(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	lenBlocks := 1000
	for i := 0; i < lenBlocks; i++ {
		h := uint32(i + 1)
		block := randomBlockWithSign(t, h, getPrevBlockHash(t, bc, h))
		assert.Nil(t, bc.AddBlock(block))
		header, err := bc.GetHeader(block.Height)
		assert.Nil(t, err)
		assert.Equal(t, header, block.Header)
	}
}

func TestAddBlockToHigh(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	assert.Nil(t, bc.AddBlock(randomBlockWithSign(t, 1, getPrevBlockHash(t, bc, genesisBlockHeight+uint32(1)))))
	assert.NotNil(t, bc.AddBlock(randomBlockWithSign(t, 3, types.Hash{})))
}

func getPrevBlockHash(t *testing.T, bc *Blockchain, height uint32) types.Hash {
	prevHeader, err := bc.GetHeader(height - 1)
	assert.Nil(t, err)

	return BlockHasher{}.Hash(prevHeader)
}

func randomBlockWithSign(t *testing.T, h uint32, prevBlockHash types.Hash) *Block {
	block := randomBlock(h, prevBlockHash)
	privateKey := crypto.GeneratePrivateKey()
	err := block.Sign(privateKey)
	assert.Nil(t, err)
	return block
}

func newBlockchainWithGenesis(t *testing.T) *Blockchain {
	genesisBlock := randomBlockWithSign(t, genesisBlockHeight, types.Hash{})
	bc, err := NewBlockchain(genesisBlock, NewMemoryStore())
	assert.Nil(t, err)
	return bc
}
