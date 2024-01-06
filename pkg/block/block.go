package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	tx "github.com/azzz/ratatoskr/pkg/transaction"
	"strconv"
	"strings"
	"time"
)

type Block struct {
	Timestamp     uint64
	Transactions  []tx.Transaction
	PrevBlockHash []byte
	Nonce         uint64
	Hash          []byte
}

func New(transactions []tx.Transaction, prevBlockHash []byte) Block {
	ts := time.Now().UnixMicro()
	return Block{
		Timestamp:     uint64(ts),
		Transactions:  transactions,
		PrevBlockHash: prevBlockHash,
	}
}

func NewGenesis(coinbase tx.Transaction) Block {
	ts := time.Now().UnixMicro()
	return Block{
		Timestamp:     uint64(ts),
		Transactions:  []tx.Transaction{coinbase},
		PrevBlockHash: []byte{},
	}
}

func (b Block) TransactionsHash() []byte {
	var hashes [][]byte

	for _, tx := range b.Transactions {
		hashes = append(hashes, tx.ID)
	}

	hash := sha256.Sum256(bytes.Join(hashes, []byte{}))

	return hash[:]
}

func (b Block) String() string {
	header := []string{
		strconv.Itoa(int(b.Nonce)),
		hex.EncodeToString(b.Hash),
	}

	return strings.Join(header, ";")
}

func (b Block) Serialize() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(b); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func DeserializeBlock(data []byte) (Block, error) {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(&block); err != nil {
		return Block{}, err
	}

	return block, nil
}
