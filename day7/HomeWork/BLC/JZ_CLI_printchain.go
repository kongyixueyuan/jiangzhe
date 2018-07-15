package BLC


//打印区块链
func (cli *JZ_CLI) JZ_printchain(nodeID string) {

	blockchain := JZ_GetBlockchain(nodeID)
	defer blockchain.JZ_DB.Close()

	blockchain.JZ_Printchain()
}
