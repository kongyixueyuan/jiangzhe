package main

import (
	"./BLC"
	"fmt"
)

func main() {
	//1.创建区块链并自动创建创世区块
	blockChain := BLC.CreateBlockChainWithBlock()
	blockChain.AddBlockToBlockChain(int64(len(blockChain.Block)+1), "second block", blockChain.Block[len(blockChain.Block)-1].PrevHash)

	fmt.Println(blockChain)
}
