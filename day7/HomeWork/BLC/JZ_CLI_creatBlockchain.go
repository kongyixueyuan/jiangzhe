package BLC



//新建区块链
func (cli *JZ_CLI)JZ_creatBlockchain(address string, nodeID string)  {

	blockchain := JZ_CreateBlockchainWithGensisBlock(address, nodeID)
	defer blockchain.JZ_DB.Close()
}
