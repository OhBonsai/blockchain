package main

import "github.com/boltdb/bolt"

type BlockChainItr struct {
	curHash []byte
	db *bolt.DB
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
