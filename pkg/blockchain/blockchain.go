package blockchain

import (
	"errors"
	"github.com/azzz/ratatoskr/pkg/transaction"

	"github.com/azzz/ratatoskr/pkg/block"
)

type Blockchain struct {
	pow   ProofOfWork
	store Store
}

func (bc Blockchain) Iterator() Iterator {
	return Iterator{store: bc.store, head: bc.store.Tip()}
}

func (bc Blockchain) Tip() []byte {
	return bc.store.Tip()
}

func (bc *Blockchain) AddBlock(transactions []transaction.Transaction) error {
	return bc.store.AddBlock(func(tip []byte) (block.Block, error) {
		newBlock := block.New(transactions, tip)
		return bc.pow.Sign(newBlock)
	})
}

func (bc *Blockchain) addGenesis(address string) error {
	return bc.store.AddBlock(func(tip []byte) (block.Block, error) {
		if tip != nil {
			return block.Block{}, errors.New("genesis can be added only to an empty blockchain")
		}

		coinbase := transaction.NewCoinbaseTx(address, "")
		genesis := block.NewGenesis(coinbase)

		return bc.pow.Sign(genesis)
	})
}

func (bc Blockchain) FindUnspendTransactions(address string) []transaction.Transaction {
	return nil
}
