package blockchain

import (
	"encoding/hex"
	"github.com/azzz/ratatoskr/pkg/block"
	"github.com/azzz/ratatoskr/pkg/proofofwork"
)

type InMemoryStore struct {
	blocks map[string]block.Block
	tip    []byte
}

func NewInMemoryBlockChain() Blockchain {
	store := &InMemoryStore{
		blocks: make(map[string]block.Block),
		tip:    nil,
	}

	bc := Blockchain{
		pow:   proofofwork.New(),
		store: store,
	}

	return bc
}

func (s *InMemoryStore) Tip() []byte {
	return s.tip
}

func (s *InMemoryStore) FindBlock(hash []byte) (block.Block, error) {
	bc, ok := s.blocks[hex.EncodeToString(hash)]
	if !ok {
		return bc, BlockNotFoundErr
	}

	return bc, nil
}

func (s *InMemoryStore) AddBlock(pow powHandler) error {
	newBlock, err := pow(s.tip)
	if err != nil {
		return err
	}

	s.blocks[hex.EncodeToString(newBlock.Hash)] = newBlock
	s.tip = newBlock.Hash
	return nil
}
