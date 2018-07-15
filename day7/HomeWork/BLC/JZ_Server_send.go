package BLC

import (
	"fmt"
	"io"
	"bytes"
	"log"
	"net"
)

//COMMAND_VERSION
func JZ_sendVersion(toAddress string, blc *JZ_Blockchain)  {


	bestHeight := blc.JZ_GetBestHeight()
	payload := JZ_gobEncode(JZ_Version{NODE_VERSION, bestHeight, nodeAddress})

	request := append(JZ_commandToBytes(COMMAND_VERSION), payload...)

	JZ_sendData(toAddress, request)
}



//COMMAND_GETBLOCKS
func JZ_sendGetBlocks(toAddress string)  {

	payload := JZ_gobEncode(JZ_GetBlocks{nodeAddress})

	request := append(JZ_commandToBytes(COMMAND_GETBLOCKS), payload...)

	JZ_sendData(toAddress, request)

}

// 主节点将自己的所有的区块hash发送给钱包节点
//COMMAND_BLOCK
//
func JZ_sendInv(toAddress string, kind string, hashes [][]byte) {

	payload := JZ_gobEncode(JZ_Inv{nodeAddress,kind,hashes})

	request := append(JZ_commandToBytes(COMMAND_INV), payload...)

	JZ_sendData(toAddress, request)

}

func JZ_sendGetData(toAddress string, kind string ,blockHash []byte) {

	payload := JZ_gobEncode(JZ_GetData{nodeAddress,kind,blockHash})

	request := append(JZ_commandToBytes(COMMAND_GETDATA), payload...)

	JZ_sendData(toAddress, request)
}


func JZ_sendBlock(toAddress string, blockBytes []byte)  {


	payload := JZ_gobEncode(JZ_BlockData{nodeAddress,blockBytes})

	request := append(JZ_commandToBytes(COMMAND_BLOCK), payload...)

	JZ_sendData(toAddress, request)
}

// 客户端向服务器发送数据
func JZ_sendData(to string, data []byte) {

	fmt.Println("Client send message to server...")

	conn, err := net.Dial("tcp", to)
	if err != nil {

		panic("error")
	}
	defer conn.Close()

	// 要发送的数据
	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {

		log.Panic(err)
	}
}
