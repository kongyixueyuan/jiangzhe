package BLC

// 用于请求区块或交易
type JZ_GetData struct {
	JZ_AddrFrom string
	JZ_Type     string
	JZ_Hash       []byte
}
