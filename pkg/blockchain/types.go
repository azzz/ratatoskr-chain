package blockchain

import "github.com/azzz/ratatoskr/pkg/block"

type Store interface {
	Tip() []byte
	FindBlock(hash []byte) (block.Block, error)
	AddBlock(pow powHandler) error
}

type ProofOfWork interface {
	Sign(block block.Block) (block.Block, error)
}
