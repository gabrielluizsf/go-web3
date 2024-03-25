package network

import (
	"testing"

	"github.com/gabrielluizsf/go-web3/core"
	"github.com/gabrielluizsf/go-web3/util"
	"github.com/stretchr/testify/assert"
)

func TestTransactionMaxLength(t *testing.T) {
	p := NewTransactionPool(1)
	p.Add(util.NewRandomTransaction(10))
	assert.Equal(t, 1, p.all.Count())

	p.Add(util.NewRandomTransaction(10))
	p.Add(util.NewRandomTransaction(10))
	p.Add(util.NewRandomTransaction(10))
	transaction := util.NewRandomTransaction(100)
	p.Add(transaction)
	assert.Equal(t, 1, p.all.Count())
	assert.True(t, p.Contains(transaction.Hash(core.TransactionHasher{})))
}

func TestTransactionPoolAdd(t *testing.T) {
	p := NewTransactionPool(11)
	n := 10

	for i := 1; i <= n; i++ {
		Transaction := util.NewRandomTransaction(100)
		p.Add(Transaction)
		// cannot add twice
		p.Add(Transaction)

		assert.Equal(t, i, p.PendingCount())
		assert.Equal(t, i, p.pending.Count())
		assert.Equal(t, i, p.all.Count())
	}
}

func TestTransactionPoolMaxLength(t *testing.T) {
	maxLen := 10
	p := NewTransactionPool(maxLen)
	n := 100
	transactions := []*core.Transaction{}

	for i := 0; i < n; i++ {
		transaction := util.NewRandomTransaction(100)
		p.Add(transaction)

		if i > n-(maxLen+1) {
			transactions = append(transactions, transaction)
		}
	}

	assert.Equal(t, p.all.Count(), maxLen)
	assert.Equal(t, len(transactions), maxLen)

	for _, Transaction := range transactions {
		assert.True(t, p.Contains(Transaction.Hash(core.TransactionHasher{})))
	}
}

func TestTransactionSortedMapFirst(t *testing.T) {
	m := NewTransactionSortedMap()
	first := util.NewRandomTransaction(100)
	m.Add(first)
	m.Add(util.NewRandomTransaction(10))
	m.Add(util.NewRandomTransaction(10))
	m.Add(util.NewRandomTransaction(10))
	m.Add(util.NewRandomTransaction(10))
	assert.Equal(t, first, m.First())
}

func TestTransactionSortedMapAdd(t *testing.T) {
	m := NewTransactionSortedMap()
	n := 100

	for i := 0; i < n; i++ {
		transaction := util.NewRandomTransaction(100)
		m.Add(transaction)
		// cannot add the same twice
		m.Add(transaction)

		assert.Equal(t, m.Count(), i+1)
		assert.True(t, m.Contains(transaction.Hash(core.TransactionHasher{})))
		assert.Equal(t, len(m.lookup), m.transactions.Len())
		assert.Equal(t, m.Get(transaction.Hash(core.TransactionHasher{})), transaction)
	}

	m.Clear()
	assert.Equal(t, m.Count(), 0)
	assert.Equal(t, len(m.lookup), 0)
	assert.Equal(t, m.transactions.Len(), 0)
}

func TestTransactionSortedMapRemove(t *testing.T) {
	m := NewTransactionSortedMap()

	transaction := util.NewRandomTransaction(100)
	m.Add(transaction)
	assert.Equal(t, m.Count(), 1)

	m.Remove(transaction.Hash(core.TransactionHasher{}))
	assert.Equal(t, m.Count(), 0)
	assert.False(t, m.Contains(transaction.Hash(core.TransactionHasher{})))
}
