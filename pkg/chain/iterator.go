package chain

import (
	"go.etcd.io/bbolt"
)

type Iterator struct {
	current []byte
	db      *bbolt.DB
}

func (i *Iterator) Next() *Block {
	var block Block

	if i.current == nil {
		return nil
	}

	err := i.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(blocksKey)
		encoded := b.Get(i.current)

		var err error
		block, err = DeserializeBlock(encoded)
		return err
	})

	if err != nil {
		return nil
	}

	i.current = block.Hash

	return &block
}
