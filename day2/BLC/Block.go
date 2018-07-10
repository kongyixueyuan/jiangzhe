package BLC

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

//区块的结构体
type Block struct {
	JZ_Height    int64
	JZ_Data      []byte
	JZ_Timestamp int64
	JZ_PrevHash  []byte
	JZ_Hash      []byte
	JZ_Nonce     int64
}

//创建一个新的区块
func JZ_NewBlock(height int64, data []byte, prev []byte) *Block {
	block := &Block{
		height,
		data,
		time.Now().Unix(),
		prev,
		nil,
		0,
	}
	/*****生成Hash******/
	block.JZ_Hash = block.SetHash()

	/*****POW算法******/
	//创建pow对象
	pow := JZ_NewPow(block)
	pow.JZ_Run()

	return block
}

//创建创世区块
func JZ_CreateGenesisBlock(data string) *Block {
	height := 1
	currentData := data
	prevhash := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	block := JZ_NewBlock(
		int64(height),
		[]byte(currentData),
		prevhash,
	)

	return block
}

//生成Hash
func (blc *Block) SetHash() []byte {
	heightBytes := Int64ToBytes(blc.JZ_Height)
	timeBytes := []byte(strconv.FormatInt(blc.JZ_Timestamp, 2))

	buff := [][]byte{
		heightBytes,
		blc.JZ_Data,
		timeBytes,
		blc.JZ_Hash,
	}

	buffRes := bytes.Join(buff, []byte{})

	hash := sha256.Sum256(buffRes)
	return hash[:]
}
