package BLC

import "fmt"

func (cli *JZ_CLI) JZ_ResetUTXOSet(nodeID string)  {


	fmt.Println("TestMethod")

	blockchain := JZ_GetBlockchain(nodeID)
	defer blockchain.JZ_DB.Close()

	utxoSet := &JZ_UTXOSet{blockchain}
	utxoSet.JZ_ResetUTXOSet()

	fmt.Println(blockchain.JZ_FindUTXOMap())
}
