package main

func main(){
	bc := NewBlockChain("bonsai")

	defer bc.db.Close()

	cli := CLI{bc}
	cli.Run()
}
