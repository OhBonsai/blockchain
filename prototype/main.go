package main

import "fmt"

func main(){
	bc := NewBlockChain()

	bc.AddBlock("Send 1 BTC to Bonsai")
	bc.AddBlock("Send 1 BTC to Sarah")

	for _, b := range bc.blocks {
		fmt.Printf("Prev.hash is %x \n", b.PrevBlockHash)
		fmt.Printf("Data is %s \n", b.Data)
		fmt.Printf("Hash is %x \n", b.PrevBlockHash)

	}
}
