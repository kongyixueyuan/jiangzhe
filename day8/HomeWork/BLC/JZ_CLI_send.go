package BLC

import (
	"fmt"
	"strconv"
)

//转账
func (cli *JZ_CLI) JZ_send(from []string, to []string, amount []string, nodeID string, mineNow bool)  {

	blockchain := JZ_GetBlockchain(nodeID)
	utxoSet := &JZ_UTXOSet{blockchain}
	defer blockchain.JZ_DB.Close()

	//打包交易并挖矿
	if mineNow {

		blockchain.JZ_MineNewBlock(from, to, amount, nodeID)

		//转账成功以后，需要更新UTXOSet
		utxoSet.JZ_Update()
	}else {

		// 把交易发送到矿工节点去进行验证
		fmt.Println("miner deal with the Tx...")

		// 遍历每一笔转账构造交易
		var txs []*JZ_Transaction
		nodeAddress = fmt.Sprintf("localhost:%s", nodeID)
		for index, address := range from {

			value, _ := strconv.Atoi(amount[index])
			tx := JZ_NewTransaction(address, to[index], int64(value), utxoSet, txs, nodeID)
			txs = append(txs, tx)

			// 将交易发送给主节点
			JZ_sendTx(knowedNodes[0], tx)
		}
	}
}