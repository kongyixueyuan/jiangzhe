package BLC

import (
	"log"
	"github.com/boltdb/bolt"
	"encoding/hex"
	"fmt"
	"os"
	"bytes"
)


//存储未花费交易输出的数据库表
const UTXOTableName  = "UTXOTableName"

type JZ_UTXOSet struct {

	JZ_Blockchain *JZ_Blockchain
}

// 1.重置数据库表
func (utxoSet *JZ_UTXOSet) JZ_ResetUTXOSet()  {

	err := utxoSet.JZ_Blockchain.JZ_DB.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(UTXOTableName))

		// 删除原有UTXO表
		if b != nil {

			err := tx.DeleteBucket([]byte(UTXOTableName))
			if err!= nil {

				log.Panic(err)
			}
		}

		// 新建UTXO表
		b ,_ = tx.CreateBucket([]byte(UTXOTableName))
		if b != nil {

			//找到链上所有UTXO并存入数据库
			txOutputsMap := utxoSet.JZ_Blockchain.JZ_FindUTXOMap()

			for keyHash,outs := range txOutputsMap {

				txHash,_ := hex.DecodeString(keyHash)

				b.Put(txHash,outs.JZ_Serialize())

			}
		}

		return nil

	})
	if err != nil {

		log.Panic(err)
	}
}

// 2.查询某个地址的UTXO
func (utxoSet *JZ_UTXOSet) JZ_FindUTXOsForAddress(address string) []*JZ_UTXO {

	var utxos []*JZ_UTXO

	err := utxoSet.JZ_Blockchain.JZ_DB.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(UTXOTableName))

		// 游标
		c := b.Cursor()
		for k, v := c.First(); k != nil; k,v = c.Next() {

			txOutputs := JZ_DeserializeTXOutputs(v)

			for _, utxo := range txOutputs.JZ_UTXOS {

				if utxo.JZ_Output.JZ_UnLockScriptPubKeyWithAddress(address) {

					utxos = append(utxos,utxo)
				}
			}
		}

		return nil
	})
	if err != nil {

		log.Panic(err)
	}

	return utxos
}

// 3.查询余额
func (utxoSet *JZ_UTXOSet) JZ_GetBalance(address string) int64 {

	UTXOS := utxoSet.JZ_FindUTXOsForAddress(address)

	var amount int64

	for _, utxo := range UTXOS  {

		amount += utxo.JZ_Output.JZ_Value
	}

	return amount
}

// 返回要凑多少钱，对应TXOutput的TX的Hash和index ???Set本身就是UTXO集合，里面的不全是未花费吗？？？？
func (utxoSet *JZ_UTXOSet) JZ_FindUnPackageSpendableUTXOS(address string, txs []*JZ_Transaction) []*JZ_UTXO {

	var unUTXOs []*JZ_UTXO
	spentTXOutputs := make(map[string][]int)

	for _,tx := range txs {

		if tx.JZ_IsCoinbaseTransaction() == false {

			for _, in := range tx.JZ_Vins {

				//是否能够解锁
				if in.JZ_UnlockWithAddress(address) {

					key := hex.EncodeToString(in.JZ_TxHash)
					spentTXOutputs[key] = append(spentTXOutputs[key], in.JZ_Vout)
				}
			}
		}
	}

	for _,tx := range txs {

	Work:
		for index,out := range tx.JZ_Vouts {

			if out.JZ_UnLockScriptPubKeyWithAddress(address) {

				if len(spentTXOutputs) != 0 {

					for hash,indexArray := range spentTXOutputs {

						txHashStr := hex.EncodeToString(tx.JZ_TxHAsh)

						if hash == txHashStr {

							var isUnSpent =true
							for _,outIndex := range indexArray {

								if index == outIndex {

									isUnSpent = false
									continue Work
								}

								if isUnSpent {

									utxo := &JZ_UTXO{tx.JZ_TxHAsh, index, out}
									unUTXOs = append(unUTXOs, utxo)
								}
							}
						} else {

							utxo := &JZ_UTXO{tx.JZ_TxHAsh, index, out}
							unUTXOs = append(unUTXOs, utxo)
						}
					}
				} else {

					utxo := &JZ_UTXO{tx.JZ_TxHAsh, index, out}
					unUTXOs = append(unUTXOs, utxo)
				}
			}
		}
	}

	return unUTXOs
}

//转账时查找可用的用于消费的UTXO
func (utxoSet *JZ_UTXOSet) JZ_FindSpendableUTXOs(address string,amount int64,txs []*JZ_Transaction) (int64,map[string][]int)  {

	unPackageUTXOS := utxoSet.JZ_FindUnPackageSpendableUTXOS(address, txs)

	spentableUTXO := make(map[string][]int)

	var value int64 = 0

	for _, UTXO := range unPackageUTXOS {

		value += UTXO.JZ_Output.JZ_Value
		txHash := hex.EncodeToString(UTXO.JZ_TxHash)
		spentableUTXO[txHash] = append(spentableUTXO[txHash], UTXO.JZ_Index)

		if value >= amount{

			return  value, spentableUTXO
		}
	}

	// 钱还不够
	err := utxoSet.JZ_Blockchain.JZ_DB.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(UTXOTableName))

		if b != nil {

			c := b.Cursor()
		UTXOBREAK:
			for k, v := c.First(); k != nil; k, v = c.Next() {

				txOutputs := JZ_DeserializeTXOutputs(v)

				for _, utxo := range txOutputs.JZ_UTXOS {

					value += utxo.JZ_Output.JZ_Value
					txHash := hex.EncodeToString(utxo.JZ_TxHash)
					spentableUTXO[txHash] = append(spentableUTXO[txHash], utxo.JZ_Index)

					if value >= amount {

						break UTXOBREAK
					}
				}
			}
		}

		return nil
	})
	if err != nil {

		log.Panic(err)
	}

	if value < amount{

		fmt.Printf("%s found.余额不足...", value)
		os.Exit(1)
	}

	return  value, spentableUTXO
}

//更新UTXO 
func (utxoSet *JZ_UTXOSet) JZ_Update()  {

	// 1.找出最新区块
	block := utxoSet.JZ_Blockchain.JZ_Iterator().JZ_Next()

	// 未花费的UTXO  键为对应交易哈希，值为TXOutput数组
	outsMap := make(map[string] *JZ_TXOutputs)
	// 新区快的交易输入,这些交易输入引用的TXOutput被消耗，应该从UTXOSet删除
	ins := []*JZ_TXInput{}

	// 2.遍历区块交易找出交易输入
	for _, tx := range block.JZ_Txs {

		//遍历交易输入，
		for _, in := range tx.JZ_Vins {

			ins = append(ins, in)
		}
	}

	// 2.遍历交易输出
	for _, tx := range block.JZ_Txs {

		utxos := []*JZ_UTXO{}

		for index, out := range tx.JZ_Vouts {

			//未花费标志
			isUnSpent := true
			for _, in := range ins {

				if in.JZ_Vout == index && bytes.Compare(tx.JZ_TxHAsh, in.JZ_TxHash) == 0 &&
					bytes.Compare(out.JZ_Ripemd160Hash, JZ_Ripemd160Hash(in.JZ_PublicKey)) == 0 {

						isUnSpent = false
						continue
				}
			}

			if isUnSpent {

				utxo := &JZ_UTXO{tx.JZ_TxHAsh,index,out}
				utxos = append(utxos,utxo)
			}
		}

		if len(utxos) > 0 {

			txHash := hex.EncodeToString(tx.JZ_TxHAsh)
			outsMap[txHash] = &JZ_TXOutputs{utxos}
		}
	}

	//3. 删除已消耗的TXOutput
	err := utxoSet.JZ_Blockchain.JZ_DB.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(UTXOTableName))
		if b != nil {

			for _, in := range ins {

				txOutputsBytes := b.Get(in.JZ_TxHash)

				//如果该交易输入无引用的交易哈希
				if len(txOutputsBytes) == 0 {

					continue
				}
				txOutputs := JZ_DeserializeTXOutputs(txOutputsBytes)

				// 判断是否需要
				isNeedDelete := false

				//缓存来自该交易还未花费的UTXO
				utxos := []*JZ_UTXO{}

				for _, utxo := range txOutputs.JZ_UTXOS {

					if in.JZ_Vout == utxo.JZ_Index && bytes.Compare(utxo.JZ_Output.JZ_Ripemd160Hash, JZ_Ripemd160Hash(in.JZ_PublicKey)) == 0 {

						isNeedDelete = true
					}else {

						//txOutputs中剩余未花费的txOutput
						utxos = append(utxos,utxo)
					}
				}

				if isNeedDelete {

					b.Delete(in.JZ_TxHash)

					if len(utxos) > 0 {

						preTXOutputs := outsMap[hex.EncodeToString(in.JZ_TxHash)]
						preTXOutputs.JZ_UTXOS = append(preTXOutputs.JZ_UTXOS, utxos...)
						outsMap[hex.EncodeToString(in.JZ_TxHash)] = preTXOutputs
					}
				}
			}

			// 4.新增交易输出到UTXOSet
			for keyHash, outPuts := range outsMap {

				keyHashBytes, _ := hex.DecodeString(keyHash)
				b.Put(keyHashBytes, outPuts.JZ_Serialize())
			}
		}

		return nil
	})
	if err != nil{

		log.Panic(err)
	}
}