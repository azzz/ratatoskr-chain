package blockchain

import (
	"github.com/azzz/ratatoskr/pkg/block"
)

type Iterator struct {
	head  []byte
	block block.Block
	err   error
	store Store
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

	var err error
	i.block, err = i.store.FindBlock(i.head)
	if err != nil {
		i.err = err
		return false
	}

	i.head = i.block.PrevBlockHash
	return true
}
