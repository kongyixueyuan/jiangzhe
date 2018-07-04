package BLC

type BlockChain struct {
	Block []*Block
}

//创建区块链并自动创建创世区块的方法
func CreateBlockChainWithBlock() *BlockChain {
	return &BlockChain{[]*Block{FirstBlock("GenesisBlock")}}
}

//将区块添加进区块链中
func (blc *BlockChain) AddBlockToBlockChain(height int64, data string, prevhash []byte) *BlockChain {
	//新创建一个区块
	newBlock := NewBlock(height, data, prevhash)
	blc.Block = append(blc.Block, newBlock)
	return blc
}
