package BLC

type TXInput struct {
	Txhash    []byte //交易的hash
	Vout      int    //存储TXOutput在Vout里面的索引
	ScriptSig string
}
