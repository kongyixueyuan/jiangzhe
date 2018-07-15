package BLC

// 表示向节点请求一个块哈希的表，该请求会返回所有块的哈希
type JZ_GetBlocks struct {
	//请求节点地址
	JZ_AddrFrom string
}