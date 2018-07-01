package BLC

type BlockChain struct {
	Block []*Block
}

func (blc *BlockChain) AddBlockToBlockChain(height int64, data string, prev []byte) *BlockChain {
	newBlock := NewBlock(height, []byte(data), prev)
	blc.Block = append(blc.Block, newBlock)
	return blc
}

func CreateGenesisBlockWithChain(data string) *BlockChain {
	blockChain := &BlockChain{[]*Block{ CreateGenesisBlock(data) }}
	return blockChain
}
