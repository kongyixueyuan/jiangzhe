package BLC

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

//存储钱包集的文件名
const WalletFile = "Wallets.dat"

type JZ_Wallets struct {
	JZ_Wallets map[string]*JZ_Wallet
}

//1.创建钱包集合
func JZ_NewWallets() (*JZ_Wallets, error) {

	//判断文件是否存在
	if _, err := os.Stat(WalletFile); os.IsNotExist(err) {

		wallets := &JZ_Wallets{}
		wallets.JZ_Wallets = make(map[string]*JZ_Wallet)

		return wallets, err
	}

	var wallets JZ_Wallets
	//读取文件
	fileContent, err := ioutil.ReadFile(WalletFile)
	if err != nil {

		log.Panic(err)
	}

	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {

		log.Panic(err)
	}

	return &wallets, err
}

//2.创建新钱包
func (wallets *JZ_Wallets) JZ_CreateWallet() {

	wallet := JZ_NewWallet()
	fmt.Printf("Your new addres：%s\n", wallet.JZ_GetAddress())
	wallets.JZ_Wallets[string(wallet.JZ_GetAddress())] = wallet

	//保存到本地
	wallets.JZ_SaveWallets()
}

//3.保存钱包集信息到文件
func (wallets *JZ_Wallets) JZ_SaveWallets() {

	var context bytes.Buffer

	//注册是为了可以序列化任何类型
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&context)
	err := encoder.Encode(&wallets)
	if err != nil {

		log.Panic(err)
	}

	// 将序列化以后的数覆盖写入到文件
	err = ioutil.WriteFile(WalletFile, context.Bytes(), 0664)
	if err != nil {

		log.Panic(err)
	}
}
