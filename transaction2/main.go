package main

func main(){
	cli := CLI{}
	//cli.getBalance("you")
	//cli.send("bonsai", "you", 2)
	cli.createBlockchain("1PKqP8Lw7DxPdVkAQ5em9DfZWzhSPyr3iq")
	cli.Run()

}
