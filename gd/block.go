package gd

import (
	"encoding/json"
	"github.com/goldoge/goldoge/common/hash"
)

type Payload struct {
	Owner string
	Star  string
}

type Block struct {
	Hash              string
	Height            int64
	Payload           Payload
	Time              int64
	PreviousBlockHash string
}

func (block *Block) Validate() bool {
	return block.Hash == RecalculateHash(*block)
}

func RecalculateHash(block Block) string {
	block.Hash = ""
	out, err := json.Marshal(block)
	if err != nil {
		panic(err)
	}
	return hash.SHA256(out)
}
