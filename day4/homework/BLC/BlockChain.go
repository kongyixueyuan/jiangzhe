package BLC

import (
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"os"
	"strconv"
	"time"
)

// 数据库名
const dbName = "blockchain.db"

// 表名
const blockTableName = "blocks"

// 用于构造区块链的结构体
type Blockchain struct {
	Tip []byte //最新的区块hash
	DB  *bolt.DB
}

// 通过迭代器遍历输出所有的区块信息
func (blc *Blockchain) Printchain() {
	blockchainIterator := blc.Iterator()
	for {
		block := blockchainIterator.Next()
		fmt.Printf("Height: %d\n", block.JZ_Height)
		fmt.Printf("PreBlockHash: %x\n", block.JZ_PrevBlockHash)
		fmt.Printf("Timestamp: %s\n", time.Unix(block.JZ_Timestamp, 0).Format("2006-01-02 03:04:05 PM"))
		fmt.Printf("Hash: %x\n", block.JZ_Hash)
		fmt.Printf("Nonce: %d\n", block.JZ_Nonce)
		fmt.Println("Transaction:")
		for _, tx := range block.JZ_Txs {
			fmt.Printf("TxHash:%x\n", tx.TxHash)
			fmt.Println("Vins:")
			for _, in := range tx.Vins {
				fmt.Printf("%x\n", in.Txhash)
				fmt.Printf("%d\n", in.Vout)
				fmt.Printf("%s\n", in.ScriptSig)
			}
			fmt.Println("Vouts:")
			for _, out := range tx.Vouts {
				fmt.Println(out.Value)
				fmt.Println(out.ScriptPubKey)
			}
		}

		var hashInt big.Int
		hashInt.SetBytes(block.JZ_PrevBlockHash)

		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}
	}

}

// 判断数据库是否存在
func JZ_DBExists() bool {
	/*
		os.Stat 查看文件信息
	*/
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		return false
	}
	return true
}

// 增加区块到区块链里面
func (blc *Blockchain) AddBlockToBlockchain(txs []*Transaction) {

	err := blc.DB.Update(func(tx *bolt.Tx) error {
		// 1. 获取表
		b := tx.Bucket([]byte(blockTableName))
		// 2. 创建新区块
		if b != nil {
			// 从数据库中获取最新区块
			blockBytes := b.Get(blc.Tip)
			// 反序列化
			block := DeserializeBlock(blockBytes)

			// 3. 将区块序列化，存储到数据库中（设置新增区块的高度和上一个区块的哈希）
			newBlock := JZ_NewBlock(txs, block.JZ_Height+1, block.JZ_Hash)

			// 保存生成新的区块到数据库中
			err := b.Put(newBlock.JZ_Hash, newBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}
			// 4. 更新数据库里面 "H" 对应的hash
			err = b.Put([]byte("H"), newBlock.JZ_Hash)
			if err != nil {
				log.Panic(err)
			}
			// 5. 更新blockchain的Tip
			blc.Tip = newBlock.JZ_Hash
		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

}

// 创建带有创世区块的区块链
func JZ_CreateBlockchainWithGenenisBlock(address string) *Blockchain {

	//判断数据库是否存在
	if JZ_DBExists() {
		fmt.Println("创世区块已经存在")
		os.Exit(1)
	}

	fmt.Println("正在创建创世区块")

	// 创建或者打开数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	var jz_genesisHash []byte

	// 更新数据库
	err = db.Update(func(tx *bolt.Tx) error {
		//创建表
		b, err := tx.CreateBucket([]byte(blockTableName))
		if err != nil {
			log.Panic(err)
		}

		if b != nil {
			// 创建创世区块
			// 创建了一个coinbase Transaction
			txCoinbase := JZ_NewCoinbaseTransaction(address)
			genesisBlock := JZ_CreateGenesisBlock([]*Transaction{txCoinbase})

			// 将创世区块存储到表中
			err := b.Put(genesisBlock.JZ_Hash, genesisBlock.Serialize())

			if err != nil {
				log.Panic(err)
			}

			// 存储最新的区块的Hash
			err = b.Put([]byte("H"), genesisBlock.JZ_Hash)

			if err != nil {
				log.Panic(err)
			}

			jz_genesisHash = genesisBlock.JZ_Hash

		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return &Blockchain{jz_genesisHash, db}

}

// 返回Blockchain对象
func JZ_BlockchainObject() *Blockchain {

	// 创建或者打开数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	var tip []byte

	err = db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(blockTableName))

		if b != nil {
			// 读取最新区块的hash
			tip = b.Get([]byte("H"))
		}

		return nil
	})
	return &Blockchain{tip, db}

}

// 如果一个地址对应的TXOutput未花费，那么这个Transaction就应该添加到数组中返回
func (blockchain *Blockchain) UnUTXOs(address string, txs []*Transaction) []*UTXO {
	// 用于存储未花费的Transaction
	var unUTXOs []*UTXO

	spentTXOutputs := make(map[string][]int)

	for _, tx := range txs {

		// Vins
		if tx.IsCoinbaseTransaction() == false {
			for _, in := range tx.Vins {
				// 是否能够解锁
				if in.UnLockWithAddress(address) {

					key := hex.EncodeToString(in.Txhash)
					spentTXOutputs[key] = append(spentTXOutputs[key], in.Vout)

				}

			}
		}

	}

	for _, tx := range txs {
	label:
		for index, out := range tx.Vouts {

			if out.UnLockScriptPubKeyWithAddress(address) {
				fmt.Println(address)
				fmt.Println(spentTXOutputs)
				if len(spentTXOutputs) == 0 {
					utxo := &UTXO{tx.TxHash, index, out}
					unUTXOs = append(unUTXOs, utxo)
				} else {
					for hash, indexArray := range spentTXOutputs {
						txHashStr := hex.EncodeToString(tx.TxHash)
						if hash == txHashStr {

							var isSpentUTXO bool

							for _, outIndex := range indexArray {
								if index == outIndex {
									isSpentUTXO = true
									continue label
								}

								if isSpentUTXO == false {
									utxo := &UTXO{tx.TxHash, index, out}
									unUTXOs = append(unUTXOs, utxo)
								}
							}
						} else {
							utxo := &UTXO{tx.TxHash, index, out}
							unUTXOs = append(unUTXOs, utxo)
						}

					}

				}

			}

		}

	}

	blockIterator := blockchain.Iterator()

	for {

		block := blockIterator.Next()
		fmt.Println(block)
		fmt.Println()
		for i := len(block.JZ_Txs) - 1; i >= 0; i-- {

			tx := block.JZ_Txs[i]
			// txHash

			// Vins
			if tx.IsCoinbaseTransaction() == false {
				for _, in := range tx.Vins {
					// 是否能够解锁
					if in.UnLockWithAddress(address) {

						key := hex.EncodeToString(in.Txhash)
						spentTXOutputs[key] = append(spentTXOutputs[key], in.Vout)

					}

				}
			}

			// Vouts
		label1:
			for index, out := range tx.Vouts {

				if out.UnLockScriptPubKeyWithAddress(address) {

					if spentTXOutputs != nil {

						if len(spentTXOutputs) != 0 {

							var isSpentUTXO bool

							for txHash, indexArray := range spentTXOutputs {

								for _, i := range indexArray {
									if index == i && txHash == hex.EncodeToString(tx.TxHash) {
										isSpentUTXO = true
										continue label1
									}
								}

							}

							if isSpentUTXO == false {
								utxo := &UTXO{tx.TxHash, index, out}
								unUTXOs = append(unUTXOs, utxo)
							}

						} else {
							utxo := &UTXO{tx.TxHash, index, out}
							unUTXOs = append(unUTXOs, utxo)
						}

					}

				}

			}

		}
		fmt.Println(spentTXOutputs)
		var hashInt big.Int
		hashInt.SetBytes(block.JZ_PrevBlockHash)

		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}

	}

	return unUTXOs

}

// 转账时查找可用的UTXO
func (blockchain *Blockchain) FindSpendableUTXOS(from string, amount int, txs []*Transaction) (int64, map[string][]int) {

	// 获取所有的UTXO
	utxos := blockchain.UnUTXOs(from, txs)

	spendableUTXO := make(map[string][]int)

	var value int64

	// 遍历utxos
	for _, utxo := range utxos {

		value = value + utxo.Output.Value

		hash := hex.EncodeToString(utxo.TxHash)
		spendableUTXO[hash] = append(spendableUTXO[hash], utxo.Index)

		if value >= int64(amount) {
			break
		}

	}

	if value < int64(amount) {
		fmt.Printf("%s's fund is not enough\n", from)
		os.Exit(1)
	}

	return value, spendableUTXO
}

// 查询余额
func (blockchain *Blockchain) GetBalance(address string) int64 {
	utxos := blockchain.UnUTXOs(address, []*Transaction{})
	var amount int64
	for _, out := range utxos {
		amount = amount + out.Output.Value
	}
	return amount
}

// 挖掘新的区块
func (blockchain *Blockchain) MineNewBlock(from []string, to []string, amount []string) {
	// 通过相关算法建立Transantion数组
	var txs []*Transaction
	for index, address := range from {
		value, _ := strconv.Atoi(amount[index])
		tx := JZ_NewSimpleTransaction(address, to[index], value, blockchain, txs)
		txs = append(txs, tx)
	}

	var block *Block
	blockchain.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			hash := b.Get([]byte("H"))
			blockBytes := b.Get(hash)
			block = DeserializeBlock(blockBytes)
		}
		return nil

	})

	// 建立新的区块
	block = JZ_NewBlock(txs, block.JZ_Height+1, block.JZ_Hash)

	// 将新区块存储到数据库
	blockchain.DB.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(blockTableName))

		if b != nil {

			b.Put(block.JZ_Hash, block.Serialize())
			b.Put([]byte("H"), block.JZ_Hash)
			blockchain.Tip = block.JZ_Hash
		}
		return nil

	})

}
