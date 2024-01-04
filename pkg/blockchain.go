package pkg

import (
	"context"
	"errors"
	"fmt"
)

type ProofOfWork interface {
	Run(ctx context.Context, block *Block) (ProofOfWorkResult, error)
}

type Blockchain struct {
	blocks []Block
	pow    ProofOfWork
}

func NewBlockchain(pow ProofOfWork) Blockchain {
	return Blockchain{
		pow: pow,
	}
}

// Init the new blockchain with creating a genesis block
func (bc *Blockchain) Init(ctx context.Context) error {
	genesis := NewGenesisBlock()

	pow, err := bc.pow.Run(ctx, &genesis)
	if err != nil {
		return fmt.Errorf("generate hash: %w", err)
	}

	genesis.Proove(pow)
	bc.blocks = append(bc.blocks, genesis)

	return nil
}

func (bc *Blockchain) Block(i int) (Block, bool) {
	if i >= len(bc.blocks) {
		return Block{}, false
	}

	return bc.blocks[i], true
}

func (bc *Blockchain) Add(ctx context.Context, data string) error {
	if len(bc.blocks) == 0 {
		return errors.New("missing genesis block")
	}

	prev := bc.blocks[len(bc.blocks)-1]
	block := newBlock(data, prev.Hash)

	pow, err := bc.pow.Run(ctx, &block)
	if err != nil {
		return fmt.Errorf("generate hash: %w", err)
	}

	block.Proove(pow)

	bc.blocks = append(bc.blocks, block)

	return nil
}
