package BLC

import (
	"log"
	"encoding/gob"
	"bytes"
	"fmt"
	"encoding/hex"
	"github.com/boltdb/bolt"
)

// Version命令处理器
func JZ_handleVersion(request []byte, blc *JZ_Blockchain)  {

	var buff bytes.Buffer
	var payload JZ_Version

	dataBytes := request[COMMANDLENGTH:]

	// 反序列化
	buff.Write(dataBytes)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {

		log.Panic(err)
	}

	// 提取最大区块高度作比较
	bestHeight := blc.JZ_GetBestHeight()
	foreignerBestHeight := payload.JZ_BestHeight

	if bestHeight > foreignerBestHeight {

		// 向请求节点回复自身Version信息
		JZ_sendVersion(payload.JZ_AddrFrom, blc)
	} else if bestHeight < foreignerBestHeight {

		// 向请求节点要信息
		JZ_sendGetBlocks(payload.JZ_AddrFrom)
	}

	// 添加到已知节点中
	if !JZ_nodeIsKnown(payload.JZ_AddrFrom) {

		knowedNodes = append(knowedNodes, payload.JZ_AddrFrom)
	}
}

func JZ_handleAddr(request []byte, blc *JZ_Blockchain)  {




}

func JZ_handleGetblocks(request []byte, blc *JZ_Blockchain)  {

	var buff bytes.Buffer
	var payload JZ_GetBlocks

	dataBytes := request[COMMANDLENGTH:]

	// 反序列化
	buff.Write(dataBytes)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blocks := blc.JZ_GetBlockHashes()

	JZ_sendInv(payload.JZ_AddrFrom, BLOCK_TYPE, blocks)
}

func JZ_handleGetData(request []byte, blc *JZ_Blockchain)  {

	var buff bytes.Buffer
	var payload JZ_GetData

	dataBytes := request[COMMANDLENGTH:]

	// 反序列化
	buff.Write(dataBytes)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {

		log.Panic(err)
	}

	if payload.JZ_Type == BLOCK_TYPE {

		block, err := blc.JZ_GetBlock([]byte(payload.JZ_Hash))
		if err != nil {

			return
		}

		fmt.Println(block)
		JZ_sendBlock(payload.JZ_AddrFrom, block)
	}

	if payload.JZ_Type == TX_TYPE {

		// 取出交易
		TxHash := hex.EncodeToString(payload.JZ_Hash)
		tx := memTxPool[TxHash]

		JZ_sendTx(payload.JZ_AddrFrom, &tx)
	}
}

func JZ_handleBlock(request []byte, blc *JZ_Blockchain)  {

	//fmt.Println("handleblock:\n")
	//blc.Printchain()

	var buff bytes.Buffer
	var payload JZ_BlockData

	dataBytes := request[COMMANDLENGTH:]

	// 反序列化
	buff.Write(dataBytes)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {

		log.Panic(err)
	}

	block := JZ_DeSerializeBlock(payload.JZ_BlockBytes)
	if block == nil {

		fmt.Printf("Block nil")
	}

	err = blc.JZ_AddBlock(block)
	if err != nil {

		log.Panic(err)
	}
	fmt.Printf("add block %x succ.\n", block.JZ_Hash)
	//blc.Printchain()

	if len(transactionArray) > 0 {

		JZ_sendGetData(payload.JZ_AddrFrom, BLOCK_TYPE, transactionArray[0])
		transactionArray = transactionArray[1:]
	}else {

		//blc.Printchain()

		utxoSet := &JZ_UTXOSet{blc}
		utxoSet.JZ_ResetUTXOSet()
	}
}

func JZ_handleTx(request []byte, blc *JZ_Blockchain)  {

	var buff bytes.Buffer
	var payload JZ_TxData

	dataBytes := request[COMMANDLENGTH:]
	buff.Write(dataBytes)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {

		log.Panic(err)
	}

	tx := JZ_DeserializeTransaction(payload.JZ_TransactionBytes)
	memTxPool[hex.EncodeToString(tx.JZ_TxHAsh)] = tx

	// 自身为主节点，需要将交易转发给矿工节点
	if nodeAddress == knowedNodes[0] {

		for _, node := range knowedNodes {

			if node != nodeAddress && node != payload.JZ_AddFrom {

				JZ_sendInv(node, TX_TYPE, [][]byte{tx.JZ_TxHAsh})
			}
		}
	} else {

		//fmt.Println(len(memTxPool), len(miningAddress))
		if len(memTxPool) >= minMinerTxCount && len(miningAddress) > 0 {

		MineTransactions:

			var txs []*JZ_Transaction
			// 创币交易，作为挖矿奖励
			coinbaseTx := JZ_NewCoinbaseTransaction(miningAddress)
			txs = append(txs, coinbaseTx)

			var verifyTxs []*JZ_Transaction

			for id := range memTxPool {

				tx := memTxPool[id]
				if blc.JZ_VerifyTransaction(&tx, verifyTxs) {

					txs = append(txs, &tx)
					verifyTxs = append(verifyTxs, &tx)
				}else {

					log.Panic("the transaction  invalid...\n")
				}
			}

			fmt.Println("All transactions verified succ!\n")

			// 建立新区块
			var block *JZ_Block
			// 取出上一个区块
			err = blc.JZ_DB.View(func(tx *bolt.Tx) error {

				b := tx.Bucket([]byte(blockTableName))
				if b != nil {

					hash := b.Get([]byte(newestBlockKey))
					block = JZ_DeSerializeBlock(b.Get(hash))
				}

				return nil
			})
			if err != nil {

				log.Panic(err)
			}

			//构造新区块
			block = JZ_NewBlock(txs, block.JZ_Height+1, block.JZ_Hash)

			fmt.Println("New block is mined!")

			// 添加到数据库
			err = blc.JZ_DB.Update(func(tx *bolt.Tx) error {

				b := tx.Bucket([]byte(blockTableName))
				if b != nil {

					b.Put(block.JZ_Hash, block.JZ_Serialize())
					b.Put([]byte(newestBlockKey), block.JZ_Hash)
					blc.JZ_Tip = block.JZ_Hash

				}
				return nil
			})
			if err != nil {

				log.Panic(err)
			}

			utxoSet := JZ_UTXOSet{blc}
			//utxoSet.Update()
			utxoSet.JZ_ResetUTXOSet()

			// 去除内存池中打包到区块的交易
			for _, tx := range txs {

				fmt.Println("delete...")
				TxHash := hex.EncodeToString(tx.JZ_TxHAsh)
				delete(memTxPool, TxHash)
			}

			// 发送区块给其他节点
			//sendBlock(knowedNodes[0], block.Serialize())
			for _, node := range knowedNodes {

				if node != nodeAddress {

					JZ_sendBlock(node, block.JZ_Serialize())
				}
			}

			if len(memTxPool) > 0 {

				goto MineTransactions
			}
		}
	}
}


func JZ_handleInv(request []byte, blc *JZ_Blockchain)  {

	var buff bytes.Buffer
	var payload JZ_Inv

	dataBytes := request[COMMANDLENGTH:]

	// 反序列化
	buff.Write(dataBytes)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	// Ivn 3000 block hashes [][]
	if payload.JZ_Type == BLOCK_TYPE {

		fmt.Println(payload.JZ_Items)

		blockHash := payload.JZ_Items[0]
		JZ_sendGetData(payload.JZ_AddrFrom, BLOCK_TYPE , blockHash)

		if len(payload.JZ_Items) >= 1 {

			transactionArray = payload.JZ_Items[1:]
		}
	}

	if payload.JZ_Type == TX_TYPE {

		TxHash := payload.JZ_Items[0]

		// 添加到交易池
		if memTxPool[hex.EncodeToString(TxHash)].JZ_TxHAsh == nil {

			JZ_sendGetData(payload.JZ_AddrFrom, TX_TYPE, TxHash)
		}
	}
}