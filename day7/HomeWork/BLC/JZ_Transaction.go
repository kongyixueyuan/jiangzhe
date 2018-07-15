package BLC

import (
	"crypto/sha256"
	"bytes"
	"encoding/gob"
	"log"
	"encoding/hex"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/elliptic"
	"math/big"
	"time"
)

type JZ_Transaction struct {
	//1.交易哈希值
	JZ_TxHAsh []byte
	//2.输入
	JZ_Vins []*JZ_TXInput
	//3.输出
	JZ_Vouts []*JZ_TXOutput
}
func JZ_NewCoinbaseTransaction(address string) *JZ_Transaction {

	//输入  由于创世区块其实没有输入，所以交易哈希传空，TXOutput索引传-1，签名随你
	txInput := &JZ_TXInput{[]byte{}, -1, []byte{}, []byte{}}
	//输出  产生一笔奖励给挖矿者
	txOutput := JZ_NewTXOutput(int64(25), address)
	txCoinbase := &JZ_Transaction{
		[]byte{},
		[]*JZ_TXInput{txInput},
		[]*JZ_TXOutput{txOutput},
	}

	txCoinbase.JZ_HashTransactions()

	return txCoinbase
}

//创币交易判断
func (tx *JZ_Transaction) JZ_IsCoinbaseTransaction() bool {

	return len(tx.JZ_Vins[0].JZ_TxHash) == 0 && tx.JZ_Vins[0].JZ_Vout == -1
}

//2.普通交易
func JZ_NewTransaction(from string, to string, amount int64, utxoSet *JZ_UTXOSet, txs []*JZ_Transaction, nodeID string) *JZ_Transaction {

	//获取钱包集合
	wallets, _ := JZ_NewWallets(nodeID)
	wallet := wallets.JZ_Wallets[from]

	money, spendableUTXODic := utxoSet.JZ_FindSpendableUTXOs(from, amount, txs)

	//输入输出
	var txInputs []*JZ_TXInput
	var txOutputs []*JZ_TXOutput

	for txHash, indexArr := range spendableUTXODic {

		//字符串转换为[]byte
		txHashBytes, _ := hex.DecodeString(txHash)
		for _, index := range indexArr {

			//交易输入
			txInput := &JZ_TXInput{
				txHashBytes,
				index,
				nil,
				wallet.JZ_PublicKey,
			}

			txInputs = append(txInputs, txInput)
		}
	}

	//转账
	txOutput := JZ_NewTXOutput(int64(amount), to)
	txOutputs = append(txOutputs, txOutput)

	//找零
	txOutput = JZ_NewTXOutput(int64(money)-int64(amount), from)
	txOutputs = append(txOutputs, txOutput)

	//交易构造
	tx := &JZ_Transaction{
		[]byte{},
		txInputs,
		txOutputs,
	}

	tx.JZ_HashTransactions()

	//进行签名
	utxoSet.JZ_Blockchain.JZ_SignTransaction(tx, wallet.JZ_PrivateKey, txs)

	return tx

func (tx *JZ_Transaction) JZ_Sign(privateKey ecdsa.PrivateKey, prevTxs map[string]JZ_Transaction) {

	//判断当前交易是否为创币交易，coinbase交易因为没有实际输入，所以没有被签名
	if tx.JZ_IsCoinbaseTransaction() {

		return
	}

	for _, vin := range tx.JZ_Vins {

		if prevTxs[hex.EncodeToString(vin.JZ_TxHash)].JZ_TxHAsh == nil {

			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	//将会被签署的是修剪后的交易副本
	txCopy := tx.JZ_TrimmedCopy()

	//遍历交易的每一个输入
	for inID, vin := range txCopy.JZ_Vins {

		//交易输入引用的上一笔交易
		prevTx := prevTxs[hex.EncodeToString(vin.JZ_TxHash)]
		//Signature 被设置为 nil
		txCopy.JZ_Vins[inID].JZ_Signature = nil
		//PubKey 被设置为所引用输出的PubKeyHash
		txCopy.JZ_Vins[inID].JZ_PublicKey = prevTx.JZ_Vouts[vin.JZ_Vout].JZ_Ripemd160Hash
		//设置交易哈希
		txCopy.JZ_TxHAsh = txCopy.JZ_Hash()
		//设置完哈希后要重置PublicKey
		txCopy.JZ_Vins[inID].JZ_PublicKey = nil

		// 签名代码
		r, s, err := ecdsa.Sign(rand.Reader, &privateKey, txCopy.JZ_TxHAsh)
		if err != nil {

			log.Panic(err)
		}
		//一个ECDSA签名就是一对数字，我们对这对数字连接起来就是signature
		signature := append(r.Bytes(), s.Bytes()...)

		tx.JZ_Vins[inID].JZ_Signature = signature
	}
}

// 验签
func (tx *JZ_Transaction) JZ_Verify(prevTxs map[string]JZ_Transaction) bool {

	if tx.JZ_IsCoinbaseTransaction() {

		return true
	}

	for _, vin := range tx.JZ_Vins {

		if prevTxs[hex.EncodeToString(vin.JZ_TxHash)].JZ_TxHAsh == nil {

			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := tx.JZ_TrimmedCopy()

	//用于椭圆曲线算法生成秘钥对
	curve := elliptic.P256()

	// 遍历输入，验证签名
	for inID, vin := range tx.JZ_Vins {

		// 这个部分跟Sign方法一样,因为在验证阶段,我们需要的是与签名相同的数据。
		prevTx := prevTxs[hex.EncodeToString(vin.JZ_TxHash)]
		txCopy.JZ_Vins[inID].JZ_Signature = nil
		txCopy.JZ_Vins[inID].JZ_PublicKey = prevTx.JZ_Vouts[vin.JZ_Vout].JZ_Ripemd160Hash
		txCopy.JZ_TxHAsh = txCopy.JZ_Hash()
		txCopy.JZ_Vins[inID].JZ_PublicKey = nil

		// 私钥
		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.JZ_Signature)
		r.SetBytes(vin.JZ_Signature[:(sigLen / 2)])
		s.SetBytes(vin.JZ_Signature[(sigLen / 2):])

		// 公钥
		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.JZ_PublicKey)
		x.SetBytes(vin.JZ_PublicKey[:(keyLen / 2)])
		y.SetBytes(vin.JZ_PublicKey[(keyLen / 2):])

		// 使用从输入提取的公钥创建ecdsa.PublicKey
		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&rawPubKey, txCopy.JZ_TxHAsh, &r, &s) == false {

			return false
		}
	}

	return true
}

// 拷贝一份新的Transaction用于签名,包含所有的输入输出，但TXInput.Signature 和 TXIput.PubKey 被设置为 nil                                 T
func (tx *JZ_Transaction) JZ_TrimmedCopy() JZ_Transaction {

	var inputs []*JZ_TXInput
	var outputs []*JZ_TXOutput

	for _, vin := range tx.JZ_Vins {

		inputs = append(inputs, &JZ_TXInput{vin.JZ_TxHash, vin.JZ_Vout, nil, nil})
	}

	for _, vout := range tx.JZ_Vouts {

		outputs = append(outputs, &JZ_TXOutput{vout.JZ_Value, vout.JZ_Ripemd160Hash})
	}

	txCopy := JZ_Transaction{tx.JZ_TxHAsh, inputs, outputs}

	return txCopy
}

//对交易信息进行哈希
func (tx *JZ_Transaction) JZ_Hash() []byte {

	txCopy := tx

	txCopy.JZ_TxHAsh = []byte{}

	hash := sha256.Sum256(txCopy.JZ_Serialize())

	return hash[:]
}

//交易序列化
func (tx *JZ_Transaction) JZ_Serialize() []byte {

	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {

		log.Panic(err)
	}

	return encoded.Bytes()
}

//将交易信息转换为字节数组
func (tx *JZ_Transaction) JZ_HashTransactions() {

	//交易信息序列化
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(tx)
	if err != nil {

		log.Panic(err)
	}

	//是创币交易的哈希不同
	timeSpBytes := JZ_IntToHex(time.Now().Unix())
	//设置hash
	txHash := sha256.Sum256(bytes.Join([][]byte{timeSpBytes, result.Bytes()}, []byte{}))
	tx.JZ_TxHAsh = txHash[:]
}
