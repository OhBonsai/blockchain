package main

import "github.com/boltdb/bolt"

type BlockChain struct {
	tip []byte
	db *bolt.DB
}

type BlockChainItr struct {
	curHash []byte
	db *bolt.DB
}

const dbFile = "./persistent/blockchain.db"
const blocksBucket = "blocks"

func NewBlockChain() *BlockChain {

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		if b == nil {
			genesis := NewGenesisBlock()
			b, _ := tx.CreateBucket([]byte(blocksBucket))
			err = b.Put(genesis.Hash, genesis.Serializer())
			err = b.Put([]byte("l"), genesis.Hash)
			tip = genesis.Hash
		}else{
			tip = b.Get([]byte("l"))
		}
		return nil
	})

	bc := BlockChain{tip , db}
	return &bc
}

func NewGenesisBlock() *Block {
	return NewBlock("Bonsai", []byte{})
}

func (bc *BlockChain) AddBlock(s string) {

	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})

	newBlock := NewBlock(s, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err = b.Put(newBlock.Hash, newBlock.Serializer())
		err = b.Put([]byte("l"), newBlock.Hash)
		bc.tip = newBlock.Hash

		return nil
	})
}


func (bc *BlockChain) Iterator() *BlockChainItr {
	bci := &BlockChainItr{bc.tip, bc.db}
	return bci
}


func (i *BlockChainItr) Next() *Block{
	var block *Block

	_ = i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.curHash)
		block = DeserializeBlock(encodedBlock)

		return nil
	})

	i.curHash = block.PrevBlockHash
	return block
}
