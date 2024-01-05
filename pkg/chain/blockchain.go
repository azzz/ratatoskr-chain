package chain

import (
	"context"
	"errors"
	"fmt"

	"go.etcd.io/bbolt"
)

type ProofOfWork interface {
	Run(ctx context.Context, block *Block) error
}

type Blockchain struct {
	tip []byte
	pow ProofOfWork
	db  *bbolt.DB
}

var (
	blocksKey = []byte("blocks")
)

func (bc Blockchain) Iterator() Iterator {
	return Iterator{
		current: bc.tip,
		db:      bc.db,
	}
}

func NewBlockChainFromState(ctx context.Context, db *bbolt.DB) (Blockchain, error) {
	pow := NewSimpleHashCash(24)
	bc := Blockchain{
		db: db, pow: pow,
	}

	err := db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(blocksKey)

		if b == nil {
			genesis := NewGenesisBlock()
			if err := pow.Run(ctx, &genesis); err != nil {
				return fmt.Errorf("hash genesis block: %w", err)
			}

			b, err := tx.CreateBucket(blocksKey)
			if err != nil {
				return fmt.Errorf("create bucket: %w", err)
			}

			data, err := genesis.Serialize()
			if err != nil {
				return fmt.Errorf("serialize genesis: %w", err)
			}

			if err := b.Put(genesis.Hash, data); err != nil {
				return fmt.Errorf("save block: %w", err)
			}

			if err := b.Put([]byte("l"), genesis.Hash); err != nil {
				return fmt.Errorf("save last hash: %w", err)
			}

			bc.tip = genesis.Hash
		} else {
			l := b.Get([]byte("l"))
			if l == nil {
				return errors.New("missing last hash key, the DB is probably broken")
			}

			bc.tip = l
		}

		return nil
	})

	return bc, err
}

// AddBlock adds a new block into the chain.
// Pay attention it's a heavy operation as it runs Proof-of-Work for the block.
func (bc Blockchain) AddBlock(ctx context.Context, data string) error {
	var lastHash []byte

	err := bc.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(blocksKey)
		lastHash = b.Get([]byte("l"))
		return nil
	})

	if err != nil {
		return fmt.Errorf("retrieve last hash: %w", err)
	}

	if lastHash == nil {
		return errors.New("missing last block hash, the db is probably broken")
	}

	newBlock := newBlock(data, lastHash)
	if err := bc.pow.Run(ctx, &newBlock); err != nil {
		return fmt.Errorf("hash block: %w", err)
	}

	err = bc.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(blocksKey)

		data, err := newBlock.Serialize()
		if err != nil {
			return fmt.Errorf("serialize genesis: %w", err)
		}

		if err := b.Put(newBlock.Hash, data); err != nil {
			return fmt.Errorf("save block: %w", err)
		}

		if err := b.Put([]byte("l"), newBlock.Hash); err != nil {
			return fmt.Errorf("save last hash: %w", err)
		}

		return nil
	})

	return nil
}
