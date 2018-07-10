package BLC

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type CLI struct {
	blockchain *Blockchain
}

func jz_printUsage() {
	fmt.Println("Usage:")
	fmt.Println("\tcreateblockchain -address -- 交易数据.")
	//fmt.Println("\taddblock -data DATA -- 交易数据.")
	fmt.Println("\tsend -from FROM -to TO -amount value --发起交易")
	fmt.Println("\tprintchain -- 输出区块信息.")
	fmt.Println("\tgetBalance -address -- 查询余额")
}
func IsvalidArgs() {
	if len(os.Args) < 2 {
		jz_printUsage()
		os.Exit(1)
	}

}
func (cli *CLI) jz_createGenesisBlockChain(address string) {
	blockchain := JZ_CreateBlockchainWithGenenisBlock(address)
	defer blockchain.DB.Close()
}
func (cli *CLI) JZ_Run() {
	//验证命令行参数是否小于2
	IsvalidArgs()
	//addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)

	//创建区块链
	jz_createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)

	//转账
	jz_sendBlockcmd := flag.NewFlagSet("send", flag.ExitOnError)

	//查看余额
	jz_getBalancecmd := flag.NewFlagSet("getBalance", flag.ExitOnError)

	//打印区块信息
	jz_printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	//转账的参数
	jz_flagFrom := jz_sendBlockcmd.String("from", "", "转账源地址")
	jz_flagTo := jz_sendBlockcmd.String("to", "", "转账目的地址")
	jz_flagAmount := jz_sendBlockcmd.String("amount", "", "转账金额")

	//创建初始区块并生成初始地址

	//生成新的区块的参数（接收一个address）
	flagCreateBlockchainWithAddress := jz_createBlockchainCmd.String("address", "", "genesis data")

	//查询余额
	getBalanceWithAddress := jz_getBalancecmd.String("address", "", "查询某个地址对应的余额")

	switch os.Args[1] {
	case "send":
		err := jz_sendBlockcmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

	case "printchain":
		err := jz_printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := jz_createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getBalance":
		err := jz_getBalancecmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		jz_printUsage()
		os.Exit(1)
	}
	if jz_sendBlockcmd.Parsed() {
		if *jz_flagFrom == "" || *jz_flagTo == "" || *jz_flagAmount == "" {
			jz_printUsage()
			os.Exit(1)
		}
		from := JSONToArray(*jz_flagFrom)
		to := JSONToArray(*jz_flagTo)
		amount := JSONToArray(*jz_flagAmount)
		cli.JZ_Send(from, to, amount)
	}
	if jz_printChainCmd.Parsed() {
		cli.JZ_Printchain()
	}
	if jz_createBlockchainCmd.Parsed() {
		if *flagCreateBlockchainWithAddress == "" {
			fmt.Println("交易数据不能为空---地址不能为空。。。。。。")
			jz_printUsage()
			os.Exit(1)
		}
		//生成创世区块的方法
		/**
		*flagCreateBlockchainWithAddress 为cli客户端发送的数据
		 */
		cli.jz_createGenesisBlockChain(*flagCreateBlockchainWithAddress)
	}
	if jz_getBalancecmd.Parsed() {
		if *getBalanceWithAddress == "" {
			fmt.Println("地址不能为空")
			jz_printUsage()
			os.Exit(1)
		}
		cli.JZ_GetBalance(*getBalanceWithAddress)
	}
}

func (cli *CLI) JZ_Printchain() {

	if JZ_DBExists() == false {
		fmt.Println("数据库不存在")
		os.Exit(1)
	}

	blockchain := JZ_BlockchainObject()

	defer blockchain.DB.Close()

	blockchain.Printchain()
}

// 转账
func (cli *CLI) JZ_Send(from []string, to []string, amount []string) {

	if JZ_DBExists() == false {
		fmt.Println("数据库不存在")
		os.Exit(1)
	}

	blockchain := JZ_BlockchainObject()
	defer blockchain.DB.Close()
	blockchain.MineNewBlock(from, to, amount)

}

// 查询余额
func (cli *CLI) JZ_GetBalance(address string) {

	blockchain := JZ_BlockchainObject()
	defer blockchain.DB.Close()

	amount := blockchain.GetBalance(address)

	fmt.Printf("%s一共有%d个Token\n", address, amount)

}
