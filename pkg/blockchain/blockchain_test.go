package blockchain

import (
	"github.com/azzz/ratatoskr/pkg/transaction"
	"github.com/stretchr/testify/assert"
	"testing"
)

func mustSend(bc Blockchain, from, to string, amount uint) {
	outputs, err := bc.FindSpendableOutputs(from, amount)
	if err != nil {
		panic(err)
	}

	tx, err := transaction.NewUTXOTransaction(from, to, amount, outputs)
	if err != nil {
		panic(err)
	}

	err = bc.AddBlock([]transaction.Transaction{tx})
	if err != nil {
		panic(err)
	}
}

func mustBalance(bc Blockchain, address string) uint {
	acc, err := bc.GetBalance(address)
	if err != nil {
		panic(err)
	}

	return acc
}

func TestBlockchain_FindUnspentTransactions(t *testing.T) {
	var (
		batman   = "batman"
		robin    = "robin"
		superman = "superman"
		flash    = "flash"
	)

	t.Run("Batman sends to Robin", func(t *testing.T) {
		bc := NewInMemoryBlockChain(batman)
		mustSend(bc, batman, robin, 3)
		assert.Equal(t, uint(2), mustBalance(bc, batman))
		assert.Equal(t, uint(3), mustBalance(bc, robin))
	})

	t.Run("Batman and Robin send to Superman", func(t *testing.T) {
		bc := NewInMemoryBlockChain(batman)
		mustSend(bc, batman, robin, 2)
		mustSend(bc, batman, superman, 2)
		mustSend(bc, robin, superman, 1)

		assert.Equal(t, uint(1), mustBalance(bc, batman))
		assert.Equal(t, uint(1), mustBalance(bc, robin))
		assert.Equal(t, uint(3), mustBalance(bc, superman))

		mustSend(bc, superman, flash, 3)
		assert.Equal(t, uint(0), mustBalance(bc, superman))
		assert.Equal(t, uint(3), mustBalance(bc, flash))
	})
}
