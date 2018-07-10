package BLC

type JZ_UTXO struct {
	//来自交易的哈希
	JZ_TxHash []byte
	//在该交易VOuts里的下标
	JZ_Index int
	//未花费的交易输出
	JZ_Output *JZ_TXOutput
}
