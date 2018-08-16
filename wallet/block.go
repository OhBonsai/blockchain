package main

import (
	"time"
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
)

type Block struct {
	Timestamp      int64
	PrevBlockHash  []byte
	Hash           []byte

	Transactions   []*Transaction
	Nonce		   int
}


func NewBlock(transactions []*Transaction, preHash []byte) *Block {
	block := &Block{
		time.Now().Unix(),
		preHash,
		[]byte{},
		transactions,
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

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions{
		txHashes = append(txHashes, tx.ID)
	}

	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txHash[:]
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



