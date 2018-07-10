package BLC

/**
  定义一个区块的结构
*/

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	//1.区块高度
	JZ_Height int64
	//2.上一个区块的HASH
	JZ_PrevBlockHash []byte
	//3.交易数据
	JZ_Txs []*Transaction
	//4.时间戳
	JZ_Timestamp int64
	//5.hash
	JZ_Hash []byte
	//6. nonce
	JZ_Nonce int64
}

//需要将Txs转换成byte
func (block *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte
	//遍历区块中的每一比交易
	for _, tx := range block.JZ_Txs {
		txHashes = append(txHashes, tx.TxHash)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txHash[:]
}

//序列化，把区块对象转成[]byte
func (block *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(block)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

//反序列化，将字节数组转成对象
func DeserializeBlock(blockBytes []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}

func JZ_NewBlock(txs []*Transaction, height int64, prevBlockHash []byte) *Block {
	//创建区块
	block := &Block{
		JZ_Height:        height,
		JZ_PrevBlockHash: prevBlockHash,
		JZ_Txs:           txs,
		JZ_Timestamp:     time.Now().Unix(),
		JZ_Hash:          nil,
		JZ_Nonce:         0,
	}
	//创建pow对象
	pow := NewProofOfWork(block)
	//获取hash和nonce
	hash, nonce := pow.Run()
	block.JZ_Hash = hash[:]
	block.JZ_Nonce = nonce
	return block
}

//生成创世区块
func JZ_CreateGenesisBlock(txs []*Transaction) *Block {
	return JZ_NewBlock(txs, 1, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})

}
