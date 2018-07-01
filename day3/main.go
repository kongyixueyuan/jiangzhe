package main

import (
	"./BLC"
)

func main() {
	//创世区块链
	blockChain := BLC.CreateGenesisBlockWithChain("jiangzhe")
	defer blockChain.DB.Close()


	//添加区块
	blockChain.AddBlockToBlockChain("zhangmengbiaa", )
	blockChain.AddBlockToBlockChain("zhangmengbiaa", )

	blockChain.PrintChain()

}
