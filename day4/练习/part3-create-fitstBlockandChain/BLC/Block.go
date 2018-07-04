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

	newBlock := &Block{
		1,
		[]byte(data),
		time.Now().Unix(),
		[]byte{0, 0, 0, 0},
		[]byte{},
	}
	newBlock.Hash = newBlock.calcHash()

	return newBlock
}
