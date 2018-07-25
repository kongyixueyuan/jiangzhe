package BLC

import "bytes"

type JZ_TXInput struct {
	//交易ID
	JZ_TxHash []byte
	//存储TXOutput在Vouts里的索引
	JZ_Vout int
	//数字签名
	JZ_Signature []byte
	//公钥
	JZ_PublicKey []byte
}

//验证当前输入是否是当前地址的
func (txInput *JZ_TXInput) JZ_UnlockWithAddress(address string) bool  {

	//base58解码
	version_pubKeyHash_checkSumBytes := JZ_Base58Decode([]byte(address))
	//去除版本得到地反编码的公钥两次哈希后的值
	ripemd160Hash := version_pubKeyHash_checkSumBytes[1:len(version_pubKeyHash_checkSumBytes)-4]

	//Ripemd160Hash算法得到公钥两次哈希后的值
	ripemd160HashNew := JZ_Ripemd160Hash(txInput.JZ_PublicKey)

	return bytes.Compare(ripemd160HashNew, ripemd160Hash) == 0
}