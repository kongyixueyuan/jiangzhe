package BLC

//打印区块链
func (cli *JZ_CLI) JZ_printchain() {

	blockchain := JZ_GetBlockchain()
	defer blockchain.JZ_DB.Close()

	blockchain.JZ_Printchain()
}
