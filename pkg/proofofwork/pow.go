package proofofwork

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"math"
	"math/big"

	"github.com/azzz/ratatoskr/pkg/block"
)

const (
	maxNonce   = math.MaxUint64
	targetBits = 16
)

// ProofOfWork algorithm than ensures if the first targetBits are zeros
type ProofOfWork struct {
	target *big.Int
}

func New() ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	return ProofOfWork{
		target: target,
	}
}

type ProofOfWorkResult struct {
	Nonce uint64
	Hash  []byte
}

func (p ProofOfWork) Validate(b block.Block) bool {
	var hashInt big.Int

	headers := p.prepareHeaders(b.Candidate, b.Nonce)
	data := bytes.Join(headers, []byte{})
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(p.target) == -1 && hash == [32]byte(b.Hash)
}

func (p ProofOfWork) Sign(b block.Candidate) (block.Block, error) {
	var (
		nonce   uint64 = 0
		hashInt big.Int
		headers = p.prepareHeaders(b, 0)
	)

	for nonce < maxNonce {
		headers[4] = uintToBytes(nonce)
		data := bytes.Join(headers, []byte{})
		hash := sha256.Sum256(data)

		hashInt.SetBytes(hash[:])

		// hash is lower than target
		if hashInt.Cmp(p.target) == -1 {
			return block.Block{
				Candidate: b,
				Nonce:     nonce,
				Hash:      hash[:],
			}, nil
		} else {
			nonce++
		}
	}

	return block.Block{}, errors.New("reached the maximal nonce without a result")
}

func (p ProofOfWork) prepareHeaders(b block.Candidate, nonce uint64) [][]byte {
	return [][]byte{
		b.PrevBlockHash,
		b.Data,
		uintToBytes(b.Timestamp),
		p.target.Bytes(),
		uintToBytes(nonce),
	}
}

func uintToBytes(u uint64) []byte {
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, u)
	return data
}
