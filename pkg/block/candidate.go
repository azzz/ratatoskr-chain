package block

import (
	"time"
)

type Candidate struct {
	Timestamp     uint64
	Data          []byte
	PrevBlockHash []byte
}

func NewCandidate(data string, prevBlockHash []byte) Candidate {
	ts := time.Now().UnixMicro()
	return Candidate{
		Timestamp:     uint64(ts),
		Data:          []byte(data),
		PrevBlockHash: prevBlockHash,
	}
}

func NewGenesis() Candidate {
	ts := time.Now().UnixMicro()
	return Candidate{
		Timestamp:     uint64(ts),
		Data:          []byte("genesis block"),
		PrevBlockHash: []byte{},
	}
}
