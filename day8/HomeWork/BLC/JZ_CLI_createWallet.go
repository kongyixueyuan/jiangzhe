package BLC

import "fmt"

func (cli *JZ_CLI)JZ_createWallet(nodeID string)  {

	wallets, _ := JZ_NewWallets(nodeID)
	wallets.JZ_CreateWallet(nodeID)

	fmt.Println(len(wallets.JZ_Wallets))
}
