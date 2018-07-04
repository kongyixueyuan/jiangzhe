package BLC

import (
	"fmt"
	"flag"
	"os"
	"log"
)

type CLI struct {

}

func isValid() {
	if len(os.Args) < 4 {
		printUsage()
		os.Exit(1)
	}
}

func (cli CLI) Run() {
	isValid()
	/**
		1.先写usage
	 */
	createBlockchain := flag.NewFlagSet("createBlockchain", flag.ExitOnError)
	flagCreateBlockchainWithAddress := createBlockchain.String("data", "", "创建创世区块链")

	switch os.Args[1] {
	case "createBlockchain":
		err := createBlockchain.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	}

	if createBlockchain.Parsed() {
		if *flagCreateBlockchainWithAddress == "" {
			fmt.Println("交易数据不能为空---地址不能为空。。。。。。")
			printUsage()
			os.Exit(1)
		}
		fmt.Println(cli.createGenesisBlockChain(*flagCreateBlockchainWithAddress))
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("\tcreateBlockchain -data \"string\" ")
}

func flagSetting() {
	//flag
}

func (cli CLI) createGenesisBlockChain(data string) *BlockChain {
	return CreateBlockChainWithBlock(data)
}