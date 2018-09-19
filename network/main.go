package main

func main(){
	cli := CLI{}
	//cli.createWallet()
	cli.listAddresses()
	//cli.createBlockchain("1C373J4p4xm2PjFVBBkZhfTHXbVPhD86zB")

	//cli.reindexUTXO()



	cli.send("1C373J4p4xm2PjFVBBkZhfTHXbVPhD86zB", "1DfdGvFgKQphS36JdcsAKj67hKvV4mpTPd", 9)
	//cli.printChan()
	//cli.reindexUTXO()
	//cli.getBalance("1C373J4p4xm2PjFVBBkZhfTHXbVPhD86zB")
	//cli.getBalance("1DfdGvFgKQphS36JdcsAKj67hKvV4mpTPd")

	//cli.printUsage()
	//cli.createBlockchain("1PKqP8Lw7DxPdVkAQ5em9DfZWzhSPyr3iq")
	//cli.Run()
}
