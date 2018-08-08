package main

import (
	"strconv"
	"bytes"
	"crypto/sha256"
	"time"
)

type Block struct {
	Timestamp      int64
	PrevBlockHash  []byte
	Hash           []byte

	Data           []byte
}


func NewBlock(data string, preHash []byte) *Block {
	block := &Block{
		time.Now().Unix(),
		preHash,
		[]byte{},
		[]byte(data),
	}
	block.SetHash()
	return block
}


func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})

	hash := sha256.Sum256(headers)

	b.Hash = hash[:]
}



