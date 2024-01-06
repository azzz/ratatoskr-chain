package blockchain

import (
	"errors"
	"fmt"
	"github.com/azzz/ratatoskr/pkg/block"
	"github.com/azzz/ratatoskr/pkg/proofofwork"
	"go.etcd.io/bbolt"
)

var (
	blocksBucket = []byte("blocks")
	tipKey       = []byte("l")
)

type PersistentStore struct {
	db *bbolt.DB
}

type powHandler func(tip []byte) (block.Block, error)

func (s PersistentStore) Tip() []byte {
	var tip []byte
	_ = s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(blocksBucket)
		if b == nil {
			return nil
		}

		tip = b.Get(tipKey)
		return nil
	})

	return tip
}

var (
	BlockNotFoundErr = errors.New("block not found")
)

func (s PersistentStore) FindBlock(hash []byte) (block.Block, error) {
	var encoded []byte
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(blocksBucket)
		if b == nil {
			return nil
		}

		encoded = b.Get(hash)
		return nil
	})

	if err != nil {
		return block.Block{}, err
	}

	if encoded == nil {
		return block.Block{}, BlockNotFoundErr
	}

	return block.DeserializeBlock(encoded)
}

func (s PersistentStore) AddBlock(pow powHandler) error {
	tx, err := s.db.Begin(true)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	b := tx.Bucket(blocksBucket)
	if b == nil {
		b, err = tx.CreateBucket(blocksBucket)
		if err != nil {
			return fmt.Errorf("create bucket: %w", err)
		}
	}

	tip := b.Get(tipKey)

	bc, err := pow(tip)
	if err != nil {
		return fmt.Errorf("proof-of-work: %w", err)
	}

	encoded, err := bc.Serialize()
	if err != nil {
		return fmt.Errorf("serialize: %w", err)
	}

	if err := b.Put(bc.Hash, encoded); err != nil {
		return fmt.Errorf("save block: %w", err)
	}

	if err := b.Put(tipKey, bc.Hash); err != nil {
		return fmt.Errorf("save tip: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func CreatePersistent(db *bbolt.DB, address string) (Blockchain, error) {
	var (
		pow   = proofofwork.New()
		store = PersistentStore{db: db}
	)

	if store.Tip() != nil {
		return Blockchain{}, errors.New("blockchain already created")
	}

	bc := Blockchain{pow, store}

	return bc, bc.addGenesis(address)
}

func LoadPersistent(db *bbolt.DB) (Blockchain, error) {
	var (
		pow   = proofofwork.New()
		store = PersistentStore{db: db}
	)

	return Blockchain{pow, store}, nil
}
