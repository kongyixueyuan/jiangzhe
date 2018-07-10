package BLC

import (
	"github.com/boltdb/bolt"
	"log"
)

//区块链迭代器
type JZ_BlockchainIterator struct {
	//当前遍历hash
	JZ_CurrHash []byte
	//区块链数据库
	JZ_DB *bolt.DB
}

func (blcIterator *JZ_BlockchainIterator) JZ_Next() *JZ_Block {

	var block *JZ_Block

	err := blcIterator.JZ_DB.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(blockTableName))

		if b != nil {

			currentBlockBytes := b.Get(blcIterator.JZ_CurrHash)

			// 获取到当前迭代器里面的currentHash所对应的区块
			block = JZ_DeSerializeBlock(currentBlockBytes)

			// 更新迭代器里面CurrentHash
			blcIterator.JZ_CurrHash = block.JZ_PrevBlockHash
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return block
}
