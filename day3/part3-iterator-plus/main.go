package main

import "blockChain/HomeWork/day3/part3-iterator-plus/BLC"

func main() {
	//创世区块链
	//blockChain := packages.CreateGenesisBlockWithChain("jiangzhe")
	//defer blockChain.DB.Close()
	//
	////添加区块
	//blockChain.AddBlockToBlockChain("zhangmengbiaa", )
	//blockChain.AddBlockToBlockChain("zhangmengbiaa", )
	//
	//blockChain.PrintChain()

	BLC.CLI{}.Run()
}
