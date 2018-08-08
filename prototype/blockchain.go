package main

type BlockChain struct {
	blocks []*Block
}

func NewBlockChain() *BlockChain {
	return &BlockChain{[]*Block{NewGenesisBlock()}}
}

func NewGenesisBlock() *Block {
	return NewBlock("Bonsai", []byte{})
}

func (bc *BlockChain) AddBlock(s string) {
	prevBlock := bc.blocks[len(bc.blocks) - 1]
	newBlock := NewBlock(s, prevBlock.Hash)

	bc.blocks = append(bc.blocks, newBlock)
}