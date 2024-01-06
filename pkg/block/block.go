package block

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"strconv"
	"strings"
)

type Block struct {
	Candidate
	Nonce uint64
	Hash  []byte
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
