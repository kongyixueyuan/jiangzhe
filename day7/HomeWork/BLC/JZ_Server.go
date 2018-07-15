package BLC

import (
	"net"
	"fmt"
	"log"
	"io/ioutil"
)


func JZ_StartServer(nodeID string, minerAdd string) {

	// 当前节点IP地址
	nodeAddress = fmt.Sprintf("localhost:%s", nodeID)

	ln, err := net.Listen(PROTOCOL, nodeAddress)
	if err != nil {

		log.Panic(err)
	}
	defer ln.Close()

	blc := JZ_GetBlockchain(nodeID)
	//fmt.Println("startserver\n")
	//blc.Printchain()

	// 第一个终端：端口为3000,启动的就是主节点
	// 第二个终端：端口为3001，钱包节点
	// 第三个终端：端口号为3002，矿工节点
	if nodeAddress != knowedNodes[0] {

		// 该节点不是主节点，钱包节点向主节点请求数据
		JZ_sendVersion(knowedNodes[0], blc)
	}

	for {

		// 接收客户端发来的数据
		connc, err := ln.Accept()
		if err != nil {

			log.Panic(err)
		}

		go JZ_handleConnection(connc, blc)
	}
}

// 客户端命令处理器
func JZ_handleConnection(conn net.Conn, blc *JZ_Blockchain) {

	//fmt.Println("handleConnection:\n")
	//blc.Printchain()

	// 读取客户端发送过来的所有的数据
	request, err := ioutil.ReadAll(conn)
	if err != nil {

		log.Panic(err)
	}

	fmt.Printf("Receive a Message:%s\n", request[:COMMANDLENGTH])

	command := JZ_bytesToCommand(request[:COMMANDLENGTH])

	switch command {

	case COMMAND_VERSION:
		JZ_handleVersion(request, blc)

	case COMMAND_ADDR:
		JZ_handleAddr(request, blc)

	case COMMAND_BLOCK:
		JZ_handleBlock(request, blc)

	case COMMAND_GETBLOCKS:
		JZ_handleGetblocks(request, blc)

	case COMMAND_GETDATA:
		JZ_handleGetData(request, blc)

	case COMMAND_INV:
		JZ_handleInv(request, blc)

	case COMMAND_TX:
		JZ_handleTx(request, blc)

	default:
		fmt.Println("Unknown command!")
	}

	defer conn.Close()
}


