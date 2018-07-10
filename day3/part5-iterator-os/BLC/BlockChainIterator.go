package BLC

import "github.com/boltdb/bolt"
import "log"

//迭代器结构体
type BlockchainIterator struct {
	CurrentHash []byte
	DB          *bolt.DB
}

//迭代器
func (blockchain *BlockChain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{blockchain.Tip, blockchain.DB}
}

func (blockchainIterator *BlockchainIterator) Next() *Block {
	var block *Block

	err := blockchainIterator.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTablename))
		if b != nil {
			//查询迭代器中的当前的hash
			blockByte := b.Get(blockchainIterator.CurrentHash)

			//序列化
			block = DeserializeBlock(blockByte)

			//更新迭代器中当前的hash
			blockchainIterator.CurrentHash = block.PrevHash
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return block
}
