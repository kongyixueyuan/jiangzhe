package BLC

import (
	"os"
	"io/ioutil"
	"crypto/elliptic"
	"encoding/gob"
	"log"
	"bytes"
	"fmt"
)

//存储钱包集的文件名
const WalletFile = "Wallets.dat"

type Wallets struct {
	WalletsMap map[string]*Wallet
}

//创建钱包集合
func NewWallets() (*Wallets, error) {

	//判断文件是否存在
	if _, err := os.Stat(WalletFile); os.IsNotExist(err) {

		wallets := &Wallets{}
		wallets.WalletsMap = make(map[string] *Wallet)

		return wallets, err
	}


	var wallets Wallets
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

//创建新钱包
func (wallets *Wallets) CreateWallet()  {

	wallet := NewWallet()
	fmt.Printf("Your new addres：%s\n",wallet.GetAddress())
	wallets.WalletsMap[string(wallet.GetAddress())] = wallet

	//保存到本地
	wallets.SaveWallets()
}

func (wallets *Wallets) SaveWallets()  {

	var context bytes.Buffer

	//注册是为了可以序列化任何类型
	gob.Register(elliptic.P256())
	encoder :=gob.NewEncoder(&context)
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