package blockchain

import (
	"errors"

	"github.com/azzz/ratatoskr/pkg/block"
	"go.etcd.io/bbolt"
)

type Iterator struct {
	head  []byte
	block block.Block
	err   error
	db    *bbolt.DB
}

func (i *Iterator) Err() error {
	return i.err
}

func (i *Iterator) Block() block.Block {
	return i.block
}

func (i *Iterator) Next() bool {
	if i.head == nil {
		return false
	}

	err := i.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(blocksBucket)
		if bucket == nil {
			return errors.New("missing bucket")
		}

		encoded := bucket.Get(i.head)

		var err error
		i.block, err = block.DeserializeBlock(encoded)
		return err
	})

	if err != nil {
		i.err = err
		return false
	}

	i.head = i.block.PrevBlockHash

	return true
}
