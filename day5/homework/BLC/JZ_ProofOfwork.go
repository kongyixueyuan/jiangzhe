package BLC

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

//期望计算的Hash值前面至少要有16个零
const targetBits = 18

type JZ_ProofOfWork struct {
	//求工作量的block
	JZ_Block *JZ_Block
	//工作量难度 big.Int大数存储
	JZ_Target *big.Int
}

//创建新的工作量证明对象
func JZ_NewProofOfWork(block *JZ_Block) *JZ_ProofOfWork {
	//1.创建一个初始值为1的target
	target := big.NewInt(1)
	//2.左移bits(Hash) - targetBit 位
	target = target.Lsh(target, 256-targetBits)

	return &JZ_ProofOfWork{block, target}
}

//拼接区块属性，返回字节数组
func (pow *JZ_ProofOfWork) JZ_prepareData(nonce int) []byte {

	data := bytes.Join(
		[][]byte{
			pow.JZ_Block.JZ_PrevBlockHash,
			pow.JZ_Block.JZ_HashTransactions(),
			JZ_IntToHex(pow.JZ_Block.JZ_Timestamp),
			JZ_IntToHex(int64(targetBits)),
			JZ_IntToHex(int64(nonce)),
			JZ_IntToHex(int64(pow.JZ_Block.JZ_Height)),
		},
		[]byte{},
	)

	return data
}

//判断当前区块是否有效
func (proofOfWork *JZ_ProofOfWork) IsValid() bool {

	//比较当前区块哈希值与目标哈希值
	var hashInt big.Int
	hashInt.SetBytes(proofOfWork.JZ_Block.JZ_Hash)

	if proofOfWork.JZ_Target.Cmp(&hashInt) == 1 {

		return true
	}

	return false
}

//运行工作量证明
func (proofOfWork *JZ_ProofOfWork) JZ_Run() ([]byte, int64) {

	//1.将Block属性拼接成字节数组

	//2.生成hash
	//3.判断Hash值有效性，如果满足条件跳出循环

	//用于寻找目标hash值的随机数
	nonce := 0
	//存储新生成的Hash值
	var hashInt big.Int
	var hash [32]byte

	fmt.Println("正在挖矿...")

	for {
		//准备数据
		dataBytes := proofOfWork.JZ_prepareData(nonce)
		//生成Hash
		hash = sha256.Sum256(dataBytes)

		//\r将当前打印行覆盖
		fmt.Printf("\r%x", hash)
		//存储Hash到hashInt
		hashInt.SetBytes(hash[:])
		//验证Hash
		if proofOfWork.JZ_Target.Cmp(&hashInt) == 1 {

			break
		}
		nonce++
	}

	return hash[:], int64(nonce)
}
