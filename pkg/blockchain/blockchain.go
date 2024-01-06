package blockchain

import (
	"errors"
	"fmt"

	"github.com/azzz/ratatoskr/pkg/block"
	"go.etcd.io/bbolt"
)

type ProofOfWork interface {
	Sign(block block.Candidate) (block.Block, error)
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

func (bc *Blockchain) AddBlock(value string) error {
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

	candidate := block.NewCandidate(value, tip)
	block, err := bc.pow.Sign(candidate)
	if err != nil {
		return fmt.Errorf("proof-of-work: %w", err)
	}

	encoded, err := block.Serialize()
	if err != nil {
		return fmt.Errorf("serialize: %w", err)
	}

	if err := b.Put(block.Hash, encoded); err != nil {
		return fmt.Errorf("save block: %w", err)
	}

	if err := b.Put(tipKey, block.Hash); err != nil {
		return fmt.Errorf("save tip: %w", err)
	}

	bc.tip = block.Hash

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

var (
	blocksBucket = []byte("blocks")
	tipKey       = []byte("l")
)
