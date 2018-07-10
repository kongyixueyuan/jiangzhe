package BLC

//新建区块链
func (cli *JZ_CLI) JZ_creatBlockchain(address string) {

	blockchain := JZ_CreateBlockchainWithGensisBlock(address)
	defer blockchain.JZ_DB.Close()
}
