package BLC

// 向其他节点展示自己拥有的区块和交易
type JZ_Inv struct {
	// 自己的地址
	JZ_AddrFrom string
	// 类型 block tx
	JZ_Type     string
	// hash二维数组
	JZ_Items    [][]byte
}
