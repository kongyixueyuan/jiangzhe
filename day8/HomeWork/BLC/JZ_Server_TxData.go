package BLC

// 同步中传递的交易类型
type JZ_TxData struct {
	// 节点地址
	JZ_AddFrom string
	// 交易
	JZ_TransactionBytes []byte
}
