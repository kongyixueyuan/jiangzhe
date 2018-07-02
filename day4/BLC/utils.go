package BLC

import (
	"encoding/binary"
	"encoding/json"
	"log"
	"bytes"
	"os"
)

//用于判断数据库是否存在的方法
func DBExists() bool {
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		return false
	}

	return true
}


//此方法用于返回BlockChain的对象
func BlockChainObj() *BlockChain {

}


// 将int64转换为字节数组
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Fatal(err)
	}

	return buff.Bytes()
}

// 标准的JSON字符串转数组
func JSONToArray(jsonString string) []string {

	//json 到 []string
	var sArr []string
	if err := json.Unmarshal([]byte(jsonString), &sArr); err != nil {
		log.Panic(err)
	}

	return sArr
}