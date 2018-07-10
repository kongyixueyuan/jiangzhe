package BLC

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

const targetBit = 20

type Pow struct {
	JZ_Block  *Block
	JZ_Target *big.Int
}

//给要加入区块链中的区块创建一个pow验证
func JZ_NewPow(block *Block) *Pow {
	//1.创建一个初始值为1的target
	target := big.NewInt(1)

	//2.做移256-targetBit
	target = target.Lsh(target, 256-targetBit)

	return &Pow{block, target}
}

//数据的拼接
func (pow *Pow) jz_prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.JZ_Block.JZ_PrevHash,
			pow.JZ_Block.JZ_Data,
			Int64ToBytes(pow.JZ_Block.JZ_Timestamp),
			Int64ToBytes(int64(targetBit)),
			Int64ToBytes(int64(nonce)),
			Int64ToBytes(int64(pow.JZ_Block.JZ_Height)),
		},
		[]byte{},
	)

	return data
}

//开始验证
func (proofOfWork *Pow) JZ_Run() ([]byte, int64) {
	var nonce = 0
	var hashInt big.Int //存储新生成的Hash
	var hash [32]byte
	for {
		//1.将block的属性拼接成字节数组
		dataBytes := proofOfWork.jz_prepareData(nonce)

		//2.生成hash
		hash = sha256.Sum256(dataBytes)
		fmt.Printf("%x\r", hash)

		//3.将Hash存储到hashInt
		hashInt.SetBytes(hash[:])

		//判断hashInt是否小于Block里面的target
		if hashInt.Cmp(proofOfWork.JZ_Target) == -1 {
			fmt.Printf("%x\n", hash)
			break
		}
		nonce = nonce + 1

	}

	fmt.Println(nonce)

	fmt.Println()

	return hash[:], int64(nonce)

	//4.判断hash有效性，如果满足条件跳出循环
}
