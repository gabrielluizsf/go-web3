package core

import (
	"fmt"
	"sync"

	"github.com/gabrielluizsf/go-web3/crypto"
	"github.com/gabrielluizsf/go-web3/types"
	"github.com/go-kit/log"
)

type Blockchain struct {
	logger log.Logger
	store  Storage
	// TODO: double check this!
	lock             sync.RWMutex
	headers          []*Header
	blocks           []*Block
	TransactionStore map[types.Hash]*Transaction
	blockStore       map[types.Hash]*Block

	accountState *AccountState

	stateLock       sync.RWMutex
	collectionState map[types.Hash]*CollectionTransaction
	mintState       map[types.Hash]*MintTransaction
	validator       Validator
	// TODO: make this an interface.
	contractState *State
}

func NewBlockchain(l log.Logger, genesis *Block) (*Blockchain, error) {
	// We should create all states inside the scope of the newblockchain.

	// TODO: read this from disk later on
	accountState := NewAccountState()

	coinbase := crypto.PublicKey{}
	accountState.CreateAccount(coinbase.Address())

	bc := &Blockchain{
		contractState:    NewState(),
		headers:          []*Header{},
		store:            NewMemorystore(),
		logger:           l,
		accountState:     accountState,
		collectionState:  make(map[types.Hash]*CollectionTransaction),
		mintState:        make(map[types.Hash]*MintTransaction),
		blockStore:       make(map[types.Hash]*Block),
		TransactionStore: make(map[types.Hash]*Transaction),
	}
	bc.validator = NewBlockValidator(bc)
	err := bc.addBlockWithoutValidation(genesis)

	return bc, err
}

func (bc *Blockchain) SetValidator(v Validator) {
	bc.validator = v
}

func (bc *Blockchain) AddBlock(b *Block) error {
	if err := bc.validator.ValidateBlock(b); err != nil {
		return err
	}

	return bc.addBlockWithoutValidation(b)
}

func (bc *Blockchain) handleNativeTransfer(transaction *Transaction) error {
	bc.logger.Log(
		"msg", "handle native token transfer",
		"from", transaction.From,
		"to", transaction.To,
		"value", transaction.Value)

	return bc.accountState.Transfer(transaction.From.Address(), transaction.To.Address(), transaction.Value)
}

func (bc *Blockchain) handleNativeNFT(transaction *Transaction) error {
	hash := transaction.Hash(TransactionHasher{})

	switch t := transaction.TransactionInner.(type) {
	case CollectionTransaction:
		bc.collectionState[hash] = &t
		bc.logger.Log("msg", "created new NFT collection", "hash", hash)
	case MintTransaction:
		_, ok := bc.collectionState[t.Collection]
		if !ok {
			return fmt.Errorf("collection (%s) does not exist on the blockchain", t.Collection)
		}
		bc.mintState[hash] = &t

		bc.logger.Log("msg", "created new NFT mint", "NFT", t.NFT, "collection", t.Collection)
	default:
		return fmt.Errorf("unsupported transaction type %v", t)
	}

	return nil
}

func (bc *Blockchain) GetBlockByHash(hash types.Hash) (*Block, error) {
	bc.lock.Lock()
	defer bc.lock.Unlock()

	block, ok := bc.blockStore[hash]
	if !ok {
		return nil, fmt.Errorf("block with hash (%s) not found", hash)
	}

	return block, nil
}

func (bc *Blockchain) GetBlock(height uint32) (*Block, error) {
	if height > bc.Height() {
		return nil, fmt.Errorf("given height (%d) too high", height)
	}

	bc.lock.Lock()
	defer bc.lock.Unlock()

	return bc.blocks[height], nil
}

func (bc *Blockchain) GetHeader(height uint32) (*Header, error) {
	if height > bc.Height() {
		return nil, fmt.Errorf("given height (%d) too high", height)
	}

	bc.lock.Lock()
	defer bc.lock.Unlock()

	return bc.headers[height], nil
}

func (bc *Blockchain) GetTransactionByHash(hash types.Hash) (*Transaction, error) {
	bc.lock.Lock()
	defer bc.lock.Unlock()

	Transaction, ok := bc.TransactionStore[hash]
	if !ok {
		return nil, fmt.Errorf("could not find transaction with hash (%s)", hash)
	}

	return Transaction, nil
}

func (bc *Blockchain) HasBlock(height uint32) bool {
	return height <= bc.Height()
}

// [0, 1, 2 ,3] => 4 len
// [0, 1, 2 ,3] => 3 height
func (bc *Blockchain) Height() uint32 {
	bc.lock.RLock()
	defer bc.lock.RUnlock()

	return uint32(len(bc.headers) - 1)
}

func (bc *Blockchain) handleTransaction(transaction *Transaction) error {
	// If we have data inside execute that data on the VM.
	if len(transaction.Data) > 0 {
		bc.logger.Log("msg", "executing code", "len", len(transaction.Data), "hash", transaction.Hash(&TransactionHasher{}))

		vm := NewVM(transaction.Data, bc.contractState)
		if err := vm.Run(); err != nil {
			return err
		}
	}

	// If the TransactionInner of the transaction is not nil we need to handle
	// the native NFT implemtation.
	if transaction.TransactionInner != nil {
		if err := bc.handleNativeNFT(transaction); err != nil {
			return err
		}
	}

	// Handle the native transaction here
	if transaction.Value > 0 {
		if err := bc.handleNativeTransfer(transaction); err != nil {
			return err
		}
	}

	return nil
}

func (bc *Blockchain) addBlockWithoutValidation(b *Block) error {
	bc.stateLock.Lock()
	for i := 0; i < len(b.Transactions); i++ {
		if err := bc.handleTransaction(b.Transactions[i]); err != nil {
			bc.logger.Log("error", err.Error())

			b.Transactions[i] = b.Transactions[len(b.Transactions)-1]
			b.Transactions = b.Transactions[:len(b.Transactions)-1]

			continue
		}
	}
	bc.stateLock.Unlock()

	// fmt.Println("========ACCOUNT STATE==============")
	// fmt.Printf("%+v\n", bc.accountState.accounts)
	// fmt.Println("========ACCOUNT STATE==============")

	bc.lock.Lock()
	bc.headers = append(bc.headers, b.Header)
	bc.blocks = append(bc.blocks, b)
	bc.blockStore[b.Hash(BlockHasher{})] = b

	for _, transaction := range b.Transactions {
		bc.TransactionStore[transaction.Hash(TransactionHasher{})] = transaction
	}
	bc.lock.Unlock()

	bc.logger.Log(
		"msg", "new block",
		"hash", b.Hash(BlockHasher{}),
		"height", b.Height,
		"transactions", len(b.Transactions),
	)

	return bc.store.Put(b)
}
