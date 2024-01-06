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
		head: bc.tip,
		db:   bc.db,
	}
}

func (bc Blockchain) IsEmpty() bool {
	return bc.tip == nil
}

func NewBlockchain(db *bbolt.DB, pow ProofOfWork) (Blockchain, error) {
	var tip []byte

	err := db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(blocksKey)
		if b == nil {
			b, err := tx.CreateBucket(blocksKey)
			if err != nil {
				return fmt.Errorf("create bucket: %w", err)
			}

			b.Put([]byte("l"), nil)
		} else {
			tip = b.Get([]byte("l"))
		}

		return nil
	})

	return Blockchain{
		tip: tip,
		db:  db,
		pow: pow,
	}, err
}

func (bc *Blockchain) AddGenesisBlock(ctx context.Context) error {
	if !bc.IsEmpty() {
		return errors.New("blockchain is not empty")
	}

	block, err := bc.addBlock(ctx, "genesis block")
	bc.tip = block.Hash
	return err
}

func (bc Blockchain) addBlock(ctx context.Context, data string) (Block, error) {
	var lastHash []byte

	err := bc.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(blocksKey)
		lastHash = b.Get([]byte("l"))
		return nil
	})

	if err != nil {
		return Block{}, fmt.Errorf("retrieve last hash: %w", err)
	}

	if lastHash == nil {
		return Block{}, errors.New("missing last block hash, the db is probably broken")
	}

	newBlock := newBlock(data, lastHash)
	if err := bc.pow.Run(ctx, &newBlock); err != nil {
		return Block{}, fmt.Errorf("hash block: %w", err)
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

	return newBlock, nil
}

// AddBlock adds a new block into the chain.
// Pay attention it's a heavy operation as it runs Proof-of-Work for the block.
func (bc *Blockchain) AddBlock(ctx context.Context, data string) error {
	if bc.IsEmpty() {
		return errors.New("blockchain is not initialized yet")
	}

	block, err := bc.addBlock(ctx, data)
	bc.tip = block.Hash
	return err
}

func (bc Blockchain) Tip() []byte {
	return bc.tip
}
