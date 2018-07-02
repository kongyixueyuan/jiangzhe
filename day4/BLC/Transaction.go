package BLC

//区块中交易的字段
type Transaction struct {
	TxHash []byte
	Vins   []*TXInput  //未花费交易输入
	Vouts  []*TXOutput //未花费交易输出
}

//判断当前的交易是否是第一笔交易(coinbase)
func (tx *Transaction) IsCoinBaseTransaction() {

}

//创世区块创建时的Transaction(交易)
func NewCoinbaseTransaction(addr string) *Transaction {

}

func (tx *Transaction) HashTransaction() {

}
