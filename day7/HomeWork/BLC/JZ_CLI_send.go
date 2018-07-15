package BLC

import "fmt"

//转账
func (cli *JZ_CLI) JZ_send(from []string, to []string, amount []string, nodeID string, mineNow bool)  {

	blockchain := JZ_GetBlockchain(nodeID)
	defer blockchain.JZ_DB.Close()

	//打包交易并挖矿
	if mineNow {

		blockchain.JZ_MineNewBlock(from, to, amount, nodeID)

		//转账成功以后，需要更新UTXOSet
		utxoSet := &JZ_UTXOSet{blockchain}
		utxoSet.JZ_Update()
	}else {

		// 把交易发送到矿工节点去进行验证
		fmt.Println("由矿工节点处理......")
	}
}