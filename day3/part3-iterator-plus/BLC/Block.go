package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"strconv"
	"time"
)

//区块的结构体

type Block struct {
	Height    int64
	Data      []byte
	Timestamp int64
	PrevHash  []byte
	Hash      []byte
	Nonce     int64
}

//创建一个新的区块
func NewBlock(height int64, data []byte, prev []byte) *Block {
	block := &Block{
		height,
		data,
		time.Now().Unix(),
		prev,
		nil,
		0,
	}
	/*****生成Hash******/
	block.Hash = block.SetHash()

	/*****POW算法******/
	//创建pow对象
	pow := NewPow(block)

	hash, nonce := pow.Run()
	block.Hash = hash
	block.Nonce = nonce

	return block
}

func CreateGenesisBlock(data string) *Block {
	height := 1
	currentData := data
	prevhash := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	block := NewBlock(
		int64(height),
		[]byte(currentData),
		prevhash,
	)

	return block
}

//生成Hash
func (blc *Block) SetHash() []byte {
	heightBytes := Int64ToBytes(blc.Height)
	timeBytes := []byte(strconv.FormatInt(blc.Timestamp, 2))

	buff := [][]byte{
		heightBytes,
		blc.Data,
		timeBytes,
		blc.PrevHash,
	}

	buffRes := bytes.Join(buff, []byte{})

	hash := sha256.Sum256(buffRes)
	return hash[:]
}

//序列化
func (block *Block) Serialize() []byte {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(block)

	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

//反序列化
func DeserializeBlock(blockBytes []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}
