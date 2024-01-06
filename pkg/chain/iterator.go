package chain

import (
	"go.etcd.io/bbolt"
)

type Iterator struct {
	head  []byte
	block Block
	err   error
	db    *bbolt.DB
}

func (i *Iterator) Err() error {
	return i.err
}

func (i *Iterator) Block() Block {
	return i.block
}

func (i *Iterator) Next() bool {
	var block Block

	if i.head == nil {
		return false
	}

	err := i.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(blocksKey)
		encoded := b.Get(i.head)

		var err error
		block, err = DeserializeBlock(encoded)
		return err
	})

	if err != nil {
		i.err = err
		return false
	}

	i.block = block
	i.head = block.PrevBlockHash

	return true
}
