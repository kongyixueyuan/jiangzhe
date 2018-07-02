package BLC

import (
	"math/big"
	"bytes"
	"crypto/sha256"
	"fmt"
)

// 256位Hash里面前面至少要有16个零
const targetBit = 16

type ProofOfWork struct {
	//当前要验证的区块
	Block *Block
	//大数存储
	target *big.Int
}

// 数据拼接，返回字节数组
func (pow *ProofOfWork) prepareData(nonce int) []byte {

	data := bytes.Join(
		[][]byte{
			pow.Block.PreBlockHash,
			pow.Block.HashTransactions(),
			IntToHex(pow.Block.Timestamp),
			IntToHex(int64(targetBit)),
			IntToHex(int64(nonce)),
			IntToHex(int64(pow.Block.Height)),
		},
		[]byte{},
	)

	return data

}

// 创建新的工作量证明对象并且给定难度值
func NewProofOfWork(block *Block) *ProofOfWork {

	// 创建一个初始值为1的target
	target := big.NewInt(1)

	// 左移256 - targetBit
	target = target.Lsh(target, 256 - targetBit)

	return &ProofOfWork{block, target}
}

func (proofOfWork *ProofOfWork) IsValid() bool {

	var hashInt big.Int
	hashInt.SetBytes(proofOfWork.Block.Hash)

	return proofOfWork.target.Cmp(&hashInt) == 1

}


func (proofOfWork *ProofOfWork) Run() ([]byte, int64) {

	nonce := 0

	var hashInt big.Int // 用来存储新生成的hash
	var hash [32]byte

	for {
		// 准备数据
		dataBytes := proofOfWork.prepareData(nonce)

		hash = sha256.Sum256(dataBytes)
		fmt.Printf("\r%x", hash)

		// 将hash存储到hashInt
		hashInt.SetBytes(hash[:])

		// 判断hashInt是否小于Block里面的target
		if proofOfWork.target.Cmp(&hashInt) == 1 {
			break
		}

		nonce = nonce + 1

	}

	return hash[:], int64(nonce)

}