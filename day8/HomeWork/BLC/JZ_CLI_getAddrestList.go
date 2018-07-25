package BLC

import "fmt"

func (cli *JZ_CLI) JZ_getAddressList(nodeID string)  {

	fmt.Println("All addresses:")

	wallets, _ := JZ_NewWallets(nodeID)
	for address, _ := range wallets.JZ_Wallets {

		fmt.Println(address)
	}
}
