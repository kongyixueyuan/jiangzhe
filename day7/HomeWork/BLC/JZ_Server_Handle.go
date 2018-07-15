package BLC

import (
	"log"
	"encoding/gob"
	"bytes"
	"fmt"
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

	if payload.JZ_Type == "tx" {

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

	}
}