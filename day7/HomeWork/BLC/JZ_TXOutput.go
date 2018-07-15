package BLC

import (
	"bytes"
)

type JZ_TXOutput struct {
	//面值
	JZ_Value int64
	//用户名
	JZ_Ripemd160Hash []byte  //用户名  公钥两次哈希后的值
}

func JZ_NewTXOutput(value int64,address string) *JZ_TXOutput {

	txOutput := &JZ_TXOutput{value,nil}

	// 设置Ripemd160Hash
	txOutput.JZ_Lock(address)

	return txOutput
}

//锁定
func (txOutput *JZ_TXOutput) JZ_Lock(address string) {

	version_pubKeyHash_checkSumBytes := JZ_Base58Decode([]byte(address))
	txOutput.JZ_Ripemd160Hash = version_pubKeyHash_checkSumBytes[1:len(version_pubKeyHash_checkSumBytes)-4]
}

//解锁
func (txOutput *JZ_TXOutput) JZ_UnLockScriptPubKeyWithAddress(address string) bool {

	version_pubKeyHash_checkSumBytes := JZ_Base58Decode([]byte(address))
	ripemd160Hash := version_pubKeyHash_checkSumBytes[1:len(version_pubKeyHash_checkSumBytes) - 4]

	//fmt.Println(txOutput.Ripemd160Hash, ripemd160Hash)
	return bytes.Compare(txOutput.JZ_Ripemd160Hash, ripemd160Hash) == 0
}


