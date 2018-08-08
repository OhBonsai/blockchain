package main

import (
	"time"
)

type Block struct {
	Timestamp      int64
	PrevBlockHash  []byte
	Hash           []byte

	Data           []byte
	Nonce		   int
}


func NewBlock(data string, preHash []byte) *Block {
	block := &Block{
		time.Now().Unix(),
		preHash,
		[]byte{},
		[]byte(data),
		0,
	}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}



