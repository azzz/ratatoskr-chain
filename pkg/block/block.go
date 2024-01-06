package block

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"strconv"
	"strings"
	"time"
)

type Block struct {
	Timestamp     uint64
	Data          []byte
	PrevBlockHash []byte
	Nonce         uint64
	Hash          []byte
}

func New(data string, prevBlockHash []byte) Block {
	ts := time.Now().UnixMicro()
	return Block{
		Timestamp:     uint64(ts),
		Data:          []byte(data),
		PrevBlockHash: prevBlockHash,
	}
}

func NewGenesis() Block {
	ts := time.Now().UnixMicro()
	return Block{
		Timestamp:     uint64(ts),
		Data:          []byte("genesis block"),
		PrevBlockHash: []byte{},
	}
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
