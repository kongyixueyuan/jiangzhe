package main

import (
	"./BLC"
	"fmt"
)

func main() {
	//1.创建一个创世区块
	block := BLC.FirstBlock("Genesis Block")

	fmt.Printf("%x\n", block.Hash)
}
