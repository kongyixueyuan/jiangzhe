package BLC

import "fmt"

//查询余额
func (cli *JZ_CLI) JZ_getBlance(address string) {

	fmt.Println("地址：" + address)

	blockchain := JZ_GetBlockchain()
	defer blockchain.JZ_DB.Close()

	//amount := blockchain.GetBalance(address)

	utxoSet := &JZ_UTXOSet{blockchain}
	amount := utxoSet.JZ_GetBalance(address)

	fmt.Printf("%s一共有%d个Token\n", address, amount)

}
