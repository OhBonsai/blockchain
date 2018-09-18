package main

func main(){
	cli := CLI{}
	//cli.createWallet()
	cli.listAddresses()
	//cli.createBlockchain("1FTZSTbL1wD5anM8HEuyYKHqxuSboURrDC")

	//cli.reindexUTXO()



	cli.send("1FTZSTbL1wD5anM8HEuyYKHqxuSboURrDC", "1DeTab2SBjxJ3QD2mqigY7acZsMvpvuv11", 9)
	//cli.printChan()
	//cli.reindexUTXO()
	//cli.getBalance("1FTZSTbL1wD5anM8HEuyYKHqxuSboURrDC")
	//cli.getBalance("1DeTab2SBjxJ3QD2mqigY7acZsMvpvuv11")

	//cli.printUsage()
	//cli.createBlockchain("1PKqP8Lw7DxPdVkAQ5em9DfZWzhSPyr3iq")
	//cli.Run()
}
