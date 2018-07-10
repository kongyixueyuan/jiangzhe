package BLC

import "fmt"

func (cli *JZ_CLI) JZ_createWallet() {

	wallets, _ := JZ_NewWallets()
	wallets.JZ_CreateWallet()

	fmt.Println(len(wallets.JZ_Wallets))
}
