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

	Transactions   []*Transaction
	Nonce		   int
	Height 		   int
}


func NewBlock(transactions []*Transaction, preHash []byte, height int) *Block {
	block := &Block{
		time.Now().Unix(),
		preHash,
		[]byte{},
		transactions,
		0,
		height,
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

	for _, tx := range b.Transactions{
		txHashes = append(txHashes, tx.Serialize())
	}

	mTree := NewMerkleTree(txHashes)
	return mTree.RootNode.Data
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



