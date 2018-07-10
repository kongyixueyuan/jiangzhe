package BLC

import "fmt"

func (cli *JZ_CLI) JZ_getAddressList() {

	fmt.Println("All addresses:")

	wallets, _ := JZ_NewWallets()
	for address, _ := range wallets.JZ_Wallets {

		fmt.Println(address)
	}
}
