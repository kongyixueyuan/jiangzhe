package BLC


//localhost:3000 主节点的地址
var knowedNodes = []string{"localhost:8700"}
var nodeAddress string //全局变量，节点地址
// 存储hash值
var transactionArray [][]byte
// 交易内存池
var memTxPool = make(map[string]JZ_Transaction)
// 矿工地址
var miningAddress string
// 挖矿需要满足的最小交易数
const minMinerTxCount = 1