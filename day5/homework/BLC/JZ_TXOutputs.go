package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
)

type JZ_TXOutputs struct {
	JZ_UTXOS []*JZ_UTXO
}

// 序列化成字节数组
func (txOutputs *JZ_TXOutputs) JZ_Serialize() []byte {

	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(txOutputs)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

// 反序列化
func JZ_DeserializeTXOutputs(txOutputsBytes []byte) *JZ_TXOutputs {

	var txOutputs JZ_TXOutputs

	decoder := gob.NewDecoder(bytes.NewReader(txOutputsBytes))
	err := decoder.Decode(&txOutputs)
	if err != nil {

		log.Panic(err)
	}

	return &txOutputs
}
