package BLC

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"time"
)

type JZ_Block struct {
	//1.区块高度
	JZ_Height int64
	//2.上一个区块HAsh
	JZ_PrevBlockHash []byte
	//3.交易数据
	JZ_Txs []*JZ_Transaction
	//4.时间戳
	JZ_Timestamp int64
	//5.Hash
	JZ_Hash []byte
	//6.Nonce  符合工作量证明的随机数
	JZ_Nonce int64
}

//区块序列化
func (block *JZ_Block) JZ_Serialize() []byte {

	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(block)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

//区块反序列化
func JZ_DeSerializeBlock(blockBytes []byte) *JZ_Block {

	var block JZ_Block

	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))

	err := decoder.Decode(&block)

	if err != nil {

		log.Panic(err)
	}

	return &block
}

//1.创建新的区块
func JZ_NewBlock(txs []*JZ_Transaction, height int64, prevBlockHash []byte) *JZ_Block {

	//创建区块
	block := &JZ_Block{
		JZ_Height:        height,
		JZ_PrevBlockHash: prevBlockHash,
		JZ_Txs:           txs,
		JZ_Timestamp:     time.Now().Unix(),
		JZ_Hash:          nil,
		JZ_Nonce:         0}

	//调用工作量证明返回有效的Hash
	pow := JZ_NewProofOfWork(block)
	hash, nonce := pow.JZ_Run()
	block.JZ_Hash = hash[:]
	block.JZ_Nonce = nonce

	fmt.Printf("\r######%d-%x\n", nonce, hash)

	return block
}

//单独方法生成创世区块
func JZ_CreateGenesisBlock(txs []*JZ_Transaction) *JZ_Block {

	return JZ_NewBlock(
		txs,
		1,
		[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	)
}

// 需要将Txs转换成[]byte
func (block *JZ_Block) JZ_HashTransactions() []byte {


	var transactions [][]byte

	for _, tx := range block.JZ_Txs {
		transactions = append(transactions, tx.JZ_Serialize())
	}
	mTree := JZ_NewMerkleTree(transactions)

	return mTree.JZ_RootNode.JZ_Data
}
