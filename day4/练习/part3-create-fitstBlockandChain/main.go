package main

import (
	"./BLC"
	"fmt"
)

func main() {
	//1.创建区块链并自动创建创世区块
	blockChain := BLC.CreateBlockChainWithBlock()

	fmt.Println(blockChain.Block[0])
}
