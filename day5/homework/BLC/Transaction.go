package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"log"
	"crypto/ecdsa"
	"crypto/rand"
	"math/big"
	"crypto/elliptic"
)

// UTXO模型
type Transaction struct {
	TxHash []byte      // 交易hash
	Vins   []*TXInput  // 输入
	Vouts  []*TXOutput // 输出
}

// 判断当前的交易是否是Coinbase交易
func (tx *Transaction) IsCoinbaseTransaction() bool {
	return len(tx.TxHash) == 0 && tx.Vins[0].Vout == -1
}

// Coinbase账号 （区块链中的第一笔交易）
// 创世区块创建时的Transaction	address其实就是从cli客户端接收的参数
func NewCoinbaseTransaction(address string, amount int64) *Transaction {
	// 消费记录
	txInput := &TXInput{[]byte{}, -1, nil, []byte{}}
	// 给Coinbase转账1000个token
	txOutput := &TXOutput{amount, []byte(address)}
	txCoinbase := &Transaction{[]byte{}, []*TXInput{txInput}, []*TXOutput{txOutput}}
	// 设置hash值
	txCoinbase.HashTransaction()
	return txCoinbase
}

// 转账时产生的Transaction
func NewSimpleTransaction(from string, to string, amount int, blockchain *Blockchain, txs []*Transaction) *Transaction {
	wallets, _ := NewWallets()
	wallet := wallets.WalletsMap[from];
	// 记录已花费的output
	//{hash1:[0],hash2:[2]}
	money, spendableUTXODic := blockchain.FindSpendableUTXOS(from, amount, txs)

	var txInputs []*TXInput
	var txOutputs []*TXOutput

	for txHash, indexArray := range spendableUTXODic {

		txHashBytes, _ := hex.DecodeString(txHash)
		for _, index := range indexArray {
			txInput := &TXInput{txHashBytes, index, nil, wallet.PublicKey}
			txInputs = append(txInputs, txInput)
		}

	}

	// 转账

	//txOutput := &TXOutput{int64(amount), to}
	txOutput := NewTXOutput(int64(amount), to)
	txOutputs = append(txOutputs, txOutput)

	// 找零
	//txOutput = &TXOutput{int64(money) - int64(amount), from}
	txOutput = NewTXOutput(int64(money)-int64(amount), from)
	txOutputs = append(txOutputs, txOutput)

	tx := &Transaction{[]byte{}, txInputs, txOutputs}

	// 设置hash值
	tx.HashTransaction()
	// 进行签名
	blockchain.SignTransaction(tx, wallet.PrivateKey)

	return tx

}

func (tx *Transaction) HashTransaction() {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash := sha256.Sum256(result.Bytes())
	tx.TxHash = hash[:]
}

//对交易信息进行哈希
func (tx *Transaction) Hash() []byte {

	txCopy := tx

	txCopy.TxHash = []byte{}

	hash := sha256.Sum256(txCopy.Serialize())

	return hash[:]
}

//交易序列化
func (tx *Transaction) Serialize() []byte {

	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {

		log.Panic(err)
	}

	return encoded.Bytes()
}


// 拷贝一份新的Transaction用于签名,包含所有的输入输出，但TXInput.Signature 和 TXIput.PubKey 被设置为 nil                                 T
func (tx *Transaction) TrimmedCopy() Transaction {

	var inputs []*TXInput
	var outputs []*TXOutput

	for _, vin := range tx.Vins {

		inputs = append(inputs, &TXInput{vin.Txhash, vin.Vout, nil, nil})
	}

	for _, vout := range tx.Vouts {

		outputs = append(outputs, &TXOutput{vout.Value, vout.Ripemd160Hash})
	}

	txCopy := Transaction{tx.TxHash, inputs, outputs}

	return txCopy
}

//数字签名
func (tx *Transaction) Sign(privateKey ecdsa.PrivateKey, prevTxs map[string]Transaction) {

	//判断当前交易是否为创币交易，coinbase交易因为没有实际输入，所以没有被签名
	if tx.IsCoinbaseTransaction() {

		return
	}

	for _, vin := range tx.Vins {

		if prevTxs[hex.EncodeToString(vin.Txhash)].TxHash == nil {

			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	//将会被签署的是修剪后的交易副本
	txCopy := tx.TrimmedCopy()

	//遍历交易的每一个输入
	for inID, vin := range txCopy.Vins {

		//交易输入引用的上一笔交易
		prevTx := prevTxs[hex.EncodeToString(vin.Txhash)]
		//Signature 被设置为 nil
		txCopy.Vins[inID].Signature = nil
		//PubKey 被设置为所引用输出的PubKeyHash
		txCopy.Vins[inID].PublicKey = prevTx.Vouts[vin.Vout].Ripemd160Hash
		//设置交易哈希
		txCopy.TxHash = txCopy.Hash()
		//设置完哈希后要重置PublicKey
		txCopy.Vins[inID].PublicKey = nil

		// 签名代码
		r, s, err := ecdsa.Sign(rand.Reader, &privateKey, txCopy.TxHash)
		if err != nil {

			log.Panic(err)
		}
		//一个ECDSA签名就是一对数字，我们对这对数字连接起来就是signature
		signature := append(r.Bytes(), s.Bytes()...)

		tx.Vins[inID].Signature = signature
	}
}

// 验签
func (tx *Transaction) Verify(prevTxs map[string]Transaction) bool {

	if tx.IsCoinbaseTransaction() {

		return true
	}

	for _, vin := range tx.Vins {

		if prevTxs[hex.EncodeToString(vin.Txhash)].TxHash == nil {

			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := tx.TrimmedCopy()

	//用于椭圆曲线算法生成秘钥对
	curve := elliptic.P256()

	// 遍历输入，验证签名
	for inID, vin := range tx.Vins {

		// 这个部分跟Sign方法一样,因为在验证阶段,我们需要的是与签名相同的数据。
		prevTx := prevTxs[hex.EncodeToString(vin.Txhash)]
		txCopy.Vins[inID].Signature = nil
		txCopy.Vins[inID].PublicKey = prevTx.Vouts[vin.Vout].Ripemd160Hash
		txCopy.TxHash = txCopy.Hash()
		txCopy.Vins[inID].PublicKey = nil

		// 私钥
		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])

		// 公钥
		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PublicKey)
		x.SetBytes(vin.PublicKey[:(keyLen / 2)])
		y.SetBytes(vin.PublicKey[(keyLen / 2):])

		// 使用从输入提取的公钥创建ecdsa.PublicKey
		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&rawPubKey, txCopy.TxHash, &r, &s) == false {

			return false
		}
	}

	return true
}