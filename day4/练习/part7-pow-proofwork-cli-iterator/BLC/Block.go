package BLC

import (
	"bytes"
	"crypto/sha256"
	"time"
)

type Block struct {
	Height    int64  //区块的高度
	Data      []byte //交易数据
	Timestamp int64  //时间戳
	PrevHash  []byte //上一个区块的哈希
	Hash      []byte //Hash
	Nonce     int64
}

//计算区块生成哈希的方法
func (blc *Block) calcHash() []byte {
	byteJoin := [][]byte{
		Int64ToBytes(blc.Height),
		blc.Data,
		Int64ToBytes(blc.Timestamp),
		blc.PrevHash,
	}
	newBytes := bytes.Join(byteJoin, []byte{})
	HashBytes := sha256.Sum256(newBytes)
	return HashBytes[:]
}

//创建区块
func FirstBlock(data string) *Block {
	firstBlock := NewBlock(1, data, []byte{0, 0, 0, 0})
	firstBlock.Hash = firstBlock.calcHash()
	return firstBlock
}

//创建一个函数用于创建新的区块
func NewBlock(height int64, data string, prevHash []byte) *Block {
	newBlock := &Block{
		int64(height),
		[]byte(data),
		time.Now().Unix(),
		prevHash,
		[]byte{},
		0,
	}

	//创建一个工作量证明的对象(创建对象的时候就要设置好难度值)
	powObj := NewProofWork(newBlock)
	//进行挖矿挖到后返回一个nonce和挖到的hash
	nonce, hash := powObj.Run()

	newBlock.Nonce = nonce
	newBlock.Hash = hash
	return newBlock
}
