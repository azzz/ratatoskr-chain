package blockchain

import (
	"errors"
	"fmt"
	"github.com/azzz/ratatoskr/pkg/transaction"

	"github.com/azzz/ratatoskr/pkg/block"
	"go.etcd.io/bbolt"
)

type ProofOfWork interface {
	Sign(block block.Block) (block.Block, error)
}

type Blockchain struct {
	tip []byte
	pow ProofOfWork
	db  *bbolt.DB
}

func (bc Blockchain) Iterator() Iterator {
	return Iterator{db: bc.db, head: bc.tip}
}

func (bc Blockchain) Tip() []byte {
	return bc.tip
}

func (bc *Blockchain) AddBlock(transactions []transaction.Transaction) error {
	tx, err := bc.db.Begin(true)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	b := tx.Bucket(blocksBucket)
	if b == nil {
		return errors.New("bucket is missing")
	}

	tip := b.Get(tipKey)
	if tip == nil {
		return errors.New("missing tip")
	}

	block := block.New(transactions, tip)
	signed, err := bc.pow.Sign(block)
	if err != nil {
		return fmt.Errorf("proof-of-work: %w", err)
	}

	encoded, err := signed.Serialize()
	if err != nil {
		return fmt.Errorf("serialize: %w", err)
	}

	if err := b.Put(signed.Hash, encoded); err != nil {
		return fmt.Errorf("save block: %w", err)
	}

	if err := b.Put(tipKey, signed.Hash); err != nil {
		return fmt.Errorf("save tip: %w", err)
	}

	bc.tip = signed.Hash

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func (bc Blockchain) FindUnspendTransactions(address string) []transaction.Transaction {
	return nil
}

var (
	blocksBucket = []byte("blocks")
	tipKey       = []byte("l")
)
