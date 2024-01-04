package pkg

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"math"
	"math/big"
)

const (
	maxNonce = math.MaxUint64
)

// ProofOfWork algorithm than ensures if the first targetBits are zeros
type SimpleHashCash struct {
	target     *big.Int
	targetBits uint
}

func NewSimpleHashCash(targetBits uint) SimpleHashCash {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	return SimpleHashCash{
		target: target,
	}
}

func uintToBytes(u uint64) []byte {
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, u)
	return data
}

func (p SimpleHashCash) dataToHash(block *Block, nonce uint64) []byte {
	data := bytes.Join(
		[][]byte{
			block.PrevBlockHash,
			block.Data,
			uintToBytes(block.Timestamp),
			p.target.Bytes(),
			uintToBytes(nonce),
		},
		[]byte{},
	)

	return data
}

type ProofOfWorkResult struct {
	Nonce uint64
	Hash  []byte
}

func (p SimpleHashCash) Run(ctx context.Context, block *Block) (ProofOfWorkResult, error) {
	var (
		done    = make(chan ProofOfWorkResult) // used to return a result from the goroutine
		errDone = make(chan error)             // used to return an error from the goroutine
	)

	go func() {
		var (
			nonce   uint64 = 0
			hashInt big.Int
		)

		for nonce < maxNonce {
			data := p.dataToHash(block, nonce)
			hash := sha256.Sum256(data)

			hashInt.SetBytes(hash[:])

			// hash is lower than target
			if hashInt.Cmp(p.target) == -1 {
				done <- ProofOfWorkResult{
					Nonce: nonce,
					Hash:  hash[:],
				}

				break
			} else {
				nonce++
			}
		}

		errDone <- errors.New("reached maximum nonce without result")
	}()

	select {
	case result := <-done:
		return result, nil
	case err := <-errDone:
		return ProofOfWorkResult{}, err
	case <-ctx.Done():
		return ProofOfWorkResult{}, ctx.Err()
	}
}
