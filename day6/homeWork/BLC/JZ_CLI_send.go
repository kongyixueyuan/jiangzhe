package BLC

//转账
func (cli *JZ_CLI) JZ_send(from []string, to []string, amount []string) {

	blockchain := JZ_GetBlockchain()
	defer blockchain.JZ_DB.Close()

	//打包交易并挖矿
	blockchain.JZ_MineNewBlock(from, to, amount)

	//转账成功以后，需要更新UTXOSet
	utxoSet := &JZ_UTXOSet{blockchain}
	utxoSet.JZ_Update()
}
