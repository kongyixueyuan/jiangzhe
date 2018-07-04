package BLC

type BlockChain struct {
	Block []*Block
}

//创建区块链并自动创建创世区块的方法
func CreateBlockChainWithBlock() *BlockChain {
	return &BlockChain{[]*Block{ FirstBlock("GenesisBlock") }}
}