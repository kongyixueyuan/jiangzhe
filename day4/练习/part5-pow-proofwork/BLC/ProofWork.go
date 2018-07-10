package BLC

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

const diff = 32

type ProofWork struct {
	Block  *Block   //当前要验证的区块
	target *big.Int //大数据存储
}

//数据拼接
func (pow *ProofWork) prepareData(nonce int64) []byte {
	data := [][]byte{
		Int64ToBytes(pow.Block.Height),
		pow.Block.PrevHash,
		pow.Block.Data,
		Int64ToBytes(pow.Block.Timestamp),
		pow.Block.PrevHash,
		Int64ToBytes(nonce),
	}

	bytesData := bytes.Join(data, []byte{})

	return bytesData
}

//工作量证明的对象
func NewProofWork(block *Block) *ProofWork {
	//设置难度
	//1.创建一个big.int类型的1
	diffVal := big.NewInt(1)

	//2.将1左移(256-diff)位
	diff := diffVal.Lsh(diffVal, 256-diff)
	return &ProofWork{block, diff}
}

//开始挖矿
func (powObj *ProofWork) Run() (n int64, h []byte) {
	var hashInt big.Int
	var nonce int64
	nonce = 0
	for {
		//计算哈希
		hashByte := sha256.Sum256(powObj.prepareData(nonce))
		//将hashbyte转换为大数
		hashInt.SetBytes(hashByte[:])
		//比较大小
		//0001000			000099999
		if powObj.target.Cmp(&hashInt) == 1 {
			fmt.Printf("hash: %x\r\n", hashByte)
			fmt.Printf("nonce: %d\r\n", nonce)
			break
		}

		fmt.Printf("hash: %x\r", hashByte)
		nonce++
	}

	return nonce, hashInt.Bytes()
}
