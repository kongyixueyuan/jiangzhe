package BLC

import (
	"crypto/sha256"
	"bytes"
	"encoding/gob"
	"log"
	"time"
	"fmt"
)

//区块的结构体
type Block struct {
	Height       int64          //区块的高度
	Txs          []*Transaction //区块的交易记录
	PreBlockHash []byte         //上一个区块的哈希
	Timestamp    int64          //时间戳
	Hash         []byte
	Nonce        int64
}

//Transaction需要转换成[]byte
func (block *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	//遍历出一个区块中的多笔交易
	for _, tx := range block.Txs  {
		txHashes = append(txHashes, tx.TxHash)
	}


	//拼接交易哈希 字节数组
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txHash[:]
}

//将区块序列化成字节数组
func (block *Block) Serialize() []byte {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(block)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

//将区块反序列化的函数
func DeserializeBlock(blockBytes []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))

	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}

//创建新的区块
func NewBlock(txs []*Transaction, height int64, preHash []byte) *Block {

	// 创建区块
	block := &Block{height, txs, preHash, time.Now().Unix(), nil, 0}

	// 调用工作量证明的方法并且返回有效的Hash和Nonce
	pow := NewProofOfWork(block)

	// 挖矿验证
	hash, nonce := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	fmt.Println()

	return block

}



// 单独写一个方法，生成创世区块
func CreateGenenisBlock(txs []*Transaction) *Block {

	return NewBlock(txs, 1, []byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})

}

