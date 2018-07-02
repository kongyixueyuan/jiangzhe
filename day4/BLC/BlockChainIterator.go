package BLC

import "github.com/boltdb/bolt"
import "log"

//创建迭代器结构器
type BlockChainIterator struct {
	CurrentHash []byte	//迭代器对象当前的Hash
	DB *bolt.DB		//返回数据库的连接池
}


func (blockchainIterator *BlockChainIterator) Next() *Block {
	var block *Block
	err := blockchainIterator.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			currentBlockBytes := b.Get(blockchainIterator.CurrentHash)
			// 获取到当前迭代器里面的currentHash所对于的区块
			block = DeserializeBlock(currentBlockBytes)

			// 更新迭代里面的currentHash
			blockchainIterator.CurrentHash = block.PreBlockHash
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return block
}