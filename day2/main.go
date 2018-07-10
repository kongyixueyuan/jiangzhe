package main

import (
	"./BLC"
)

func main() {
	//创世区块链
	blockChain := BLC.JZ_CreateGenesisBlockWithChain("jiangzhe")

	//添加区块
	blockChain.JZ_AddBlockToBlockChain(
		int64(len(blockChain.Block)+1),
		"zhangmengbiaa",
		blockChain.Block[len(blockChain.Block)-1].JZ_Hash,
	)

	blockChain.JZ_AddBlockToBlockChain(
		int64(len(blockChain.Block)+1),
		"zhangmengbi",
		blockChain.Block[len(blockChain.Block)-1].JZ_Hash,
	)

	blockChain.JZ_AddBlockToBlockChain(
		int64(len(blockChain.Block)+1),
		"zhangmengbi",
		blockChain.Block[len(blockChain.Block)-1].JZ_Hash,
	)

	blockChain.JZ_AddBlockToBlockChain(
		int64(len(blockChain.Block)+1),
		"zhangmengbi",
		blockChain.Block[len(blockChain.Block)-1].JZ_Hash,
	)
}
