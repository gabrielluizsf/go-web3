package network

import (
	"sync"

	"github.com/gabrielluizsf/go-web3/core"
	"github.com/gabrielluizsf/go-web3/types"
)

type TransactionPool struct {
	all     *TransactionSortedMap
	pending *TransactionSortedMap
	// The maxLength of the total pool of transactions.
	// When the pool is full we will prune the oldest transaction.
	maxLength int
}

func NewTransactionPool(maxLength int) *TransactionPool {
	return &TransactionPool{
		all:       NewTransactionSortedMap(),
		pending:   NewTransactionSortedMap(),
		maxLength: maxLength,
	}
}

func (p *TransactionPool) Add(transaction *core.Transaction) {
	// prune the oldest transaction that is sitting in the all pool
	if p.all.Count() == p.maxLength {
		oldest := p.all.First()
		p.all.Remove(oldest.Hash(core.TransactionHasher{}))
	}

	if !p.all.Contains(transaction.Hash(core.TransactionHasher{})) {
		p.all.Add(transaction)
		p.pending.Add(transaction)
	}
}

func (p *TransactionPool) Contains(hash types.Hash) bool {
	return p.all.Contains(hash)
}

// Pending returns a slice of transactions that are in the pending pool
func (p *TransactionPool) Pending() []*core.Transaction {
	return p.pending.transactions.Data
}

func (p *TransactionPool) ClearPending() {
	p.pending.Clear()
}

func (p *TransactionPool) PendingCount() int {
	return p.pending.Count()
}

type TransactionSortedMap struct {
	lock         sync.RWMutex
	lookup       map[types.Hash]*core.Transaction
	transactions *types.List[*core.Transaction]
}

func NewTransactionSortedMap() *TransactionSortedMap {
	return &TransactionSortedMap{
		lookup:       make(map[types.Hash]*core.Transaction),
		transactions: types.NewList[*core.Transaction](),
	}
}

func (t *TransactionSortedMap) First() *core.Transaction {
	t.lock.RLock()
	defer t.lock.RUnlock()

	first := t.transactions.Get(0)
	return t.lookup[first.Hash(core.TransactionHasher{})]
}

func (t *TransactionSortedMap) Get(h types.Hash) *core.Transaction {
	t.lock.RLock()
	defer t.lock.RUnlock()

	return t.lookup[h]
}

func (t *TransactionSortedMap) Add(transaction *core.Transaction) {
	hash := transaction.Hash(core.TransactionHasher{})

	t.lock.Lock()
	defer t.lock.Unlock()

	if _, ok := t.lookup[hash]; !ok {
		t.lookup[hash] = transaction
		t.transactions.Insert(transaction)
	}
}

func (t *TransactionSortedMap) Remove(h types.Hash) {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.transactions.Remove(t.lookup[h])
	delete(t.lookup, h)
}

func (t *TransactionSortedMap) Count() int {
	t.lock.RLock()
	defer t.lock.RUnlock()

	return len(t.lookup)
}

func (t *TransactionSortedMap) Contains(h types.Hash) bool {
	t.lock.RLock()
	defer t.lock.RUnlock()

	_, ok := t.lookup[h]
	return ok
}

func (t *TransactionSortedMap) Clear() {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.lookup = make(map[types.Hash]*core.Transaction)
	t.transactions.Clear()
}
