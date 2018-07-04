package main

import "blockChain/HomeWork/day4/part6-pow-proofwork-cli/BLC"

func main() {
	//1.通过cli命令去创建创世区块链
	// ./bc createBlockchain -data "jiangzhe"

	//2.通过cli命令去新加一个区块
	// ./bc addBlock -data "zhangmengbi"

	//3.通过cli命令去挖矿
	// ./bc mine

	BLC.CLI{}.Run()
}
