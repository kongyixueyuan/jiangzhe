package BLC

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"os"
	"strconv"
	"time"
)

//相关数据库属性
const dbName = "Blockchain.db"
const blockTableName = "Blocks"
const newestBlockKey = "BlockKey"

type JZ_Blockchain struct {
	//最新区块的Hash
	JZ_Tip []byte
	//存储区块的数据库
	JZ_DB *bolt.DB
}

//1.创建创世区块
func JZ_CreateBlockchainWithGensisBlock(address string) *JZ_Blockchain {

	var blc *JZ_Blockchain

	//判断数据库是否存在
	if JZ_IsDBExists(dbName) {

		fmt.Println("创世区块已存在...")
		//os.Exit(1)

		//创建并打开数据库
		db, err := bolt.Open(dbName, 0600, nil)
		if err != nil {
			log.Fatal(err)
		}

		var block *JZ_Block
		err = db.View(func(tx *bolt.Tx) error {

			b := tx.Bucket([]byte(blockTableName))
			if b != nil {

				hash := b.Get([]byte(newestBlockKey))
				blockBytes := b.Get(hash)
				block = JZ_DeSerializeBlock(blockBytes)
				fmt.Printf("\r%x\n", block.JZ_Nonce, hash)

				blc = &JZ_Blockchain{hash, db}
			}

			return nil
		})
		if err != nil {

			log.Panic(err)
		}

		return blc
		//os.Exit(1)
	}

	fmt.Println("正在创建创世区块...")

	//创建并打开数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {

		log.Fatal(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {

		b, err := tx.CreateBucket([]byte(blockTableName))
		if err != nil {

			log.Panic(err)
		}

		if b != nil {

			//创币交易
			txCoinbase := JZ_NewCoinbaseTransaction(address)
			//创世区块
			gensisBlock := JZ_CreateGenesisBlock([]*JZ_Transaction{txCoinbase})
			//存入数据库
			err := b.Put(gensisBlock.JZ_Hash, gensisBlock.JZ_Serialize())
			if err != nil {
				log.Panic(err)
			}

			//存储最新区块hash
			err = b.Put([]byte(newestBlockKey), gensisBlock.JZ_Hash)
			if err != nil {
				log.Panic(err)
			}

			blc = &JZ_Blockchain{gensisBlock.JZ_Hash, db}
		}

		return nil
	})
	//更新数据库失败
	if err != nil {
		log.Fatal(err)
	}

	//创建创世区块时候初始化UTXO表
	utxoSet := &JZ_UTXOSet{blc}
	utxoSet.JZ_ResetUTXOSet()

	return blc
}

//2.新增一个区块到区块链 --> 包含交易的挖矿
func (blc *JZ_Blockchain) JZ_MineNewBlock(from []string, to []string, amount []string) {


	//获取UTXO集
	utxoSet := &JZ_UTXOSet{blc}

	var txs []*JZ_Transaction

	//作为奖励给矿工的奖励  暂时将这笔奖励给from[0]  挖矿成功后再转给挖矿的矿工
	tx := JZ_NewCoinbaseTransaction(from[0])
	txs = append(txs, tx)

	//1.通过相关算法建立Transaction数组
	for index, address := range from {

		value, _ := strconv.Atoi(amount[index])
		tx := JZ_NewTransaction(address, to[index], int64(value), utxoSet, txs)
		txs = append(txs, tx)
	}

	//2.挖矿
	//取上个区块的哈希和高度值
	var block *JZ_Block
	err := blc.JZ_DB.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(blockTableName))
		if b != nil {

			hash := b.Get([]byte(newestBlockKey))
			blockBytes := b.Get(hash)
			block = JZ_DeSerializeBlock(blockBytes)
		}

		return nil
	})
	if err != nil {

		log.Panic(err)
	}

	//建立新区快前需要对交易进行验签
	//已经验证的交易
	verifiedTxs := []*JZ_Transaction{}
	for _, tx := range txs {

		if blc.JZ_VerifyTransaction(tx, txs) == false {

			log.Printf("The Tx:%x verify failed.\n", tx.JZ_TxHAsh)
		}

		verifiedTxs = append(verifiedTxs, tx)
	}

	//3.建立新区块
	block = JZ_NewBlock(txs, block.JZ_Height+1, block.JZ_Hash)

	//4.存储新区块
	err = blc.JZ_DB.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(blockTableName))
		if b != nil {

			//fmt.Printf("444---%x\n\n", block.Txs[0].Vins[0].TxHash)
			//fmt.Println(block)

			err = b.Put(block.JZ_Hash, block.JZ_Serialize())
			if err != nil {

				log.Panic(err)
			}

			err = b.Put([]byte(newestBlockKey), block.JZ_Hash)
			if err != nil {

				log.Panic(err)
			}

			blc.JZ_Tip = block.JZ_Hash
		}

		return nil
	})
	if err != nil {

		log.Panic(err)
		//fmt.Print(err)
	}

}

//3.X 优化区块链遍历方法
func (blc *JZ_Blockchain) JZ_Printchain() {
	//迭代器
	blcIterator := blc.JZ_Iterator()

	//block := blcIterator.Next()
	//fmt.Printf("666---%x\n\n", block.Txs[0].Vins[0].txHash)

	for {

		block := blcIterator.JZ_Next()

		fmt.Println("------------------------------")
		fmt.Printf("Height：%d\n", block.JZ_Height)
		fmt.Printf("PrevBlockHash：%x\n", block.JZ_PrevBlockHash)
		fmt.Printf("Timestamp：%s\n", time.Unix(block.JZ_Timestamp, 0).Format("2006-01-02 03:04:05 PM"))
		fmt.Printf("Hash：%x\n", block.JZ_Hash)
		fmt.Printf("Nonce：%d\n", block.JZ_Nonce)
		fmt.Println("Txs:")
		for _, tx := range block.JZ_Txs {

			fmt.Printf("%x\n", tx.JZ_TxHAsh)
			fmt.Println("Vins:")
			for _, in := range tx.JZ_Vins {
				fmt.Printf("txHash:%x\n", in.JZ_TxHash)
				fmt.Printf("Vout:%d\n", in.JZ_Vout)
				fmt.Printf("Signature:%x\n\n", in.JZ_Signature)
				fmt.Printf("PublicKey:%x\n\n", in.JZ_PublicKey)
			}

			fmt.Println("Vouts:")
			for _, out := range tx.JZ_Vouts {
				fmt.Printf("Value:%d\n", out.JZ_Value)
				fmt.Printf("Ripemd160Hash:%x\n\n", out.JZ_Ripemd160Hash)
			}
		}
		fmt.Println("------------------------------\n\n")

		var hashInt big.Int
		hashInt.SetBytes(block.JZ_PrevBlockHash)

		if big.NewInt(0).Cmp(&hashInt) == 0 {

			break
		}
	}
}

func (blc *JZ_Blockchain) JZ_Iterator() *JZ_BlockchainIterator {

	return &JZ_BlockchainIterator{blc.JZ_Tip, blc.JZ_DB}
}

//4.获取Blockchain对象
func JZ_GetBlockchain() *JZ_Blockchain {

	var blc *JZ_Blockchain
	//判断数据库是否存在
	if JZ_IsDBExists(dbName) {

		//创建并打开数据库
		db, err := bolt.Open(dbName, 0600, nil)
		if err != nil {
			log.Fatal(err)
		}

		err = db.View(func(tx *bolt.Tx) error {

			b := tx.Bucket([]byte(blockTableName))
			if b != nil {

				hash := b.Get([]byte(newestBlockKey))
				blc = &JZ_Blockchain{hash, db}
			}

			return nil
		})
		if err != nil {

			log.Panic(err)
		}
	} else {

		fmt.Println("区块链不存在...")
		os.Exit(1)
	}

	return blc
}

//获取某个交易
func (blc *JZ_Blockchain) JZ_FindTransaction(txHash []byte, txs []*JZ_Transaction) (JZ_Transaction, error) {

	//fmt.Printf("%x----%d\n\n", txHash, len(txs))
	for _, tx := range txs {

		//fmt.Printf("%x\n\n", tx.TxHAsh)
		if bytes.Compare(tx.JZ_TxHAsh, txHash) == 0 {

			return *tx, nil
		}
	}

	blcIterator := blc.JZ_Iterator()

	for {

		block := blcIterator.JZ_Next()

		for _, tx := range block.JZ_Txs {

			//fmt.Printf("%x\n\n", tx.TxHAsh)
			if bytes.Compare(tx.JZ_TxHAsh, txHash) == 0 {

				//fmt.Println("0yes")
				return *tx, nil
			}
		}

		var hashInt big.Int
		hashInt.SetBytes(block.JZ_PrevBlockHash)
		if big.NewInt(0).Cmp(&hashInt) == 0 {

			break
		}
	}

	return JZ_Transaction{}, errors.New("Transaction is not found")
}

//交易签名
func (blc *JZ_Blockchain) JZ_SignTransaction(tx *JZ_Transaction, privKey ecdsa.PrivateKey, txs []*JZ_Transaction) {

	var prevTX JZ_Transaction
	var err error
	prevTXs := make(map[string]JZ_Transaction)

	if tx.JZ_IsCoinbaseTransaction() {

		prevTX, err = blc.JZ_FindTransaction(tx.JZ_TxHAsh, txs)
		prevTXs[hex.EncodeToString(prevTX.JZ_TxHAsh)] = prevTX

	} else {

		for _, vin := range tx.JZ_Vins {

			//找到当前交易输入引用的所有交易
			fmt.Printf("txHas0:%x\n", vin.JZ_TxHash)
			prevTX, err = blc.JZ_FindTransaction(vin.JZ_TxHash, txs)
			if err != nil {

				log.Panic(err)
			}

			prevTXs[hex.EncodeToString(prevTX.JZ_TxHAsh)] = prevTX
		}
	}

	tx.JZ_Sign(privKey, prevTXs)
}

// 交易验签
func (blc *JZ_Blockchain) JZ_VerifyTransaction(tx *JZ_Transaction, txs []*JZ_Transaction) bool {

	var prevTX JZ_Transaction
	var err error
	prevTXs := make(map[string]JZ_Transaction)

	if tx.JZ_IsCoinbaseTransaction() {

		prevTX, err = blc.JZ_FindTransaction(tx.JZ_TxHAsh, txs)
		if err != nil {

			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.JZ_TxHAsh)] = prevTX

	} else {

		for _, vin := range tx.JZ_Vins {

			prevTX, err = blc.JZ_FindTransaction(vin.JZ_TxHash, txs)
			if err != nil {

				log.Panic(err)
			}
			prevTXs[hex.EncodeToString(prevTX.JZ_TxHAsh)] = prevTX
		}

	}

	return tx.JZ_Verify(prevTXs)

	//return true
}

// 查找未花费的UTXO[string]*TXOutputs 返回字典  键为所属交易的哈希，值为TXOutput数组
func (blc *JZ_Blockchain) JZ_FindUTXOMap() map[string]*JZ_TXOutputs {

	blcIterator := blc.JZ_Iterator()

	// 存储已花费的UTXO的信息
	spentableUTXOsMap := make(map[string][]*JZ_TXInput)

	utxoMaps := make(map[string]*JZ_TXOutputs)

	for {

		block := blcIterator.JZ_Next()

		for i := len(block.JZ_Txs) - 1; i >= 0; i-- {

			txOutputs := &JZ_TXOutputs{[]*JZ_UTXO{}}
			tx := block.JZ_Txs[i]

			// coinbase
			if tx.JZ_IsCoinbaseTransaction() == false {

				for _, txInput := range tx.JZ_Vins {

					txHash := hex.EncodeToString(txInput.JZ_TxHash)
					spentableUTXOsMap[txHash] = append(spentableUTXOsMap[txHash], txInput)
				}
			}

			txHash := hex.EncodeToString(tx.JZ_TxHAsh)

		WorkOutLoop:
			for index, out := range tx.JZ_Vouts {

				txInputs := spentableUTXOsMap[txHash]

				if len(txInputs) > 0 {

					isUnSpent := true

					for _, in := range txInputs {

						outPublicKey := out.JZ_Ripemd160Hash
						inPublicKey := in.JZ_PublicKey

						if bytes.Compare(outPublicKey, JZ_Ripemd160Hash(inPublicKey)) == 0 {

							if index == in.JZ_Vout {

								isUnSpent = false
								continue WorkOutLoop
							}
						}

					}

					if isUnSpent {

						utxo := &JZ_UTXO{tx.JZ_TxHAsh, index, out}
						txOutputs.JZ_UTXOS = append(txOutputs.JZ_UTXOS, utxo)
					}

				} else {

					utxo := &JZ_UTXO{tx.JZ_TxHAsh, index, out}
					txOutputs.JZ_UTXOS = append(txOutputs.JZ_UTXOS, utxo)
				}

			}

			// 设置键值对
			utxoMaps[txHash] = txOutputs
		}

		// 找到创世区块时退出
		var hashInt big.Int
		hashInt.SetBytes(block.JZ_PrevBlockHash)
		if hashInt.Cmp(big.NewInt(0)) == 0 {

			break
		}
	}

	return utxoMaps
}

//判断数据库是否存在
func JZ_IsDBExists(dbName string) bool {

	//if _, err := os.Stat(dbName); os.IsNotExist(err) {
	//
	//	return false
	//}

	_, err := os.Stat(dbName)
	if err == nil {

		return true
	}
	if os.IsNotExist(err) {

		return false
	}

	return true
}
