package BLC

import (
	"github.com/boltdb/bolt"
	"log"
	"fmt"
	"time"
	"math/big"
)

//数据库名
const dbName = "blockchain.db"

//表的名字
const blockTablename = "blocks"

type BlockChain struct {
	Tip []byte	//最新的区块
	DB *bolt.DB
}

//迭代器结构体
type BlockchainIterator struct {
	CurrentHash	[]byte
	DB *bolt.DB
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

//遍历输出所有区块的信息
func (blc *BlockChain)PrintChain() {
	//创建迭代器
	blockIterator := blc.Iterator()

	for  {
		//遍历所有区块
		block := blockIterator.Next()
		//输出
		fmt.Printf("Height: %d\n", block.Height)
		fmt.Printf("Data: %s\n", string(block.Data))
		fmt.Printf("Prev: %x\n", block.PrevHash)
		fmt.Printf("Timestamp: %s\n", time.Unix(block.Timestamp, 0).Format("2006-01-02 03:04 PM"))
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("Nonce: %d\n", block.Nonce)
		var hashInt big.Int
		hashInt.SetBytes(blockIterator.CurrentHash)
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}
	}
}

//增加区块到区块链里面
func (blc *BlockChain) AddBlockToBlockChain(data string) *BlockChain {

	err := blc.DB.Update( func(tx *bolt.Tx) error {
		//1.获取表
		b := tx.Bucket([]byte(blockTablename))
		//2. 创建新区块
		if b != nil {
			//2.1获取最新区块
			blockBytes := b.Get(blc.Tip)
			//2.2反序列化
			block := DeserializeBlock(blockBytes)

			//3.将区块序列化并且存储到数据库中
			newBlock := NewBlock(block.Height+1, []byte(data), block.Hash)
			err := b.Put(newBlock.Hash, newBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}
			//4. 更新数据库里面"l"对应的hash
			err = b.Put([]byte("l"), newBlock.Hash)
			if err != nil {
				log.Panic(err)
			}

			//5.更新blockchain的Tip
			blc.Tip = newBlock.Hash
		}
		return nil
	} )


	if err != nil {
		log.Panic(err)
	}

	return blc
}

func CreateGenesisBlockWithChain(data string) *BlockChain {
	//创建或者打开数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	var blockHash []byte	//用于存储最新区块的哈希

	//对数据库进行读写操作
	err = db.Update(func(tx *bolt.Tx) error {
		//查询表
		b := tx.Bucket([]byte(blockTablename))
		if b == nil {
			//创建数据库表
			b, err = tx.CreateBucket([]byte(blockTablename))
			if err != nil {
				log.Panic(err)
			}
		}else{
			//创建创世区块的时候将区块存入数据库
			genesisBlock := CreateGenesisBlock(data)

			//将创世区块存储到表中
			err := b.Put(genesisBlock.Hash, genesisBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}

			//存储最新的区块的哈希
			err = b.Put([]byte("l"), genesisBlock.Hash)
			if err != nil {
				log.Panic(err)
			}

			blockHash = genesisBlock.Hash



		}

		return nil
	})
	//返回区块链对象
	return &BlockChain{blockHash, db}
}
