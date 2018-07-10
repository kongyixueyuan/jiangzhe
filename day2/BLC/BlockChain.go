package BLC

type BlockChain struct {
	Block []*Block
}

//添加新的区块
func (blc *BlockChain) JZ_AddBlockToBlockChain(height int64, data string, prev []byte) *BlockChain {
	//生成新的区块
	newBlock := JZ_NewBlock(height, []byte(data), prev)
	//将新区块加入到区块链中
	blc.Block = append(blc.Block, newBlock)
	return blc
}

//创建创世区块和区块链
func JZ_CreateGenesisBlockWithChain(data string) *BlockChain {
	blockChain := &BlockChain{[]*Block{JZ_CreateGenesisBlock(data)}}
	return blockChain
}
