package pkg

import (
	"encoding/hex"
	"strconv"
	"strings"
	"time"
)

type Block struct {
	Timestamp     uint64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce         uint64
}

func NewGenesisBlock() Block {
	return newBlock("genesis block", []byte{})
}

// newBlock creates a new block without hash and nonce
func newBlock(data string, prevBlockHash []byte) Block {
	ts := time.Now().UnixNano()
	return Block{
		Timestamp:     uint64(ts),
		Data:          []byte(data),
		PrevBlockHash: prevBlockHash,
	}
}

func (bc Block) Prooved() bool {
	return bc.Hash != nil
}

func (b *Block) Proove(pow ProofOfWorkResult) {
	b.Hash = pow.Hash
	b.Nonce = pow.Nonce
}

func (b Block) String() string {
	header := []string{
		strconv.Itoa(int(b.Nonce)),
		hex.EncodeToString(b.Hash),
	}

	str := []string{
		strings.Join(header, ":"),
		string(b.Data),
	}

	return strings.Join(str, ";")
}
