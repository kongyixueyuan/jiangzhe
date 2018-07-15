package BLC

import (
	"fmt"
	"os"
)

func (cli *JZ_CLI) JZ_StartNode(nodeID string, minerAdd string)  {

	// 启动服务器
	if minerAdd == "" || JZ_IsValidForAddress([]byte(minerAdd))  {

		//  启动服务器
		fmt.Printf("start Server:localhost:%s\n", nodeID)
		JZ_StartServer(nodeID, minerAdd)

	} else {

		fmt.Println("Server address invalid....")
		os.Exit(0)
	}

}