package main

import (
	"time"
	"bytes"
	"encoding/gob"
	"log"
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

func (b *Block) Serializer() []byte {
	var result bytes.Buffer
	encoders := gob.NewEncoder(&result)
	err := encoders.Encode(b)

	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

func DeserializeBlock(d []byte)  *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewBuffer(d))
	err := decoder.Decode(&block)

	if err != nil {
		log.Panic(err)
	}

	return &block

}


