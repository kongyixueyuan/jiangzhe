package BLC

import (
	"math/big"
	"bytes"
)
//base58编码集
var b58Alphabet = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

// 字节数组转 Base58,加密
func JZ_Base58Encode(input []byte) []byte {

	var result []byte

	x := big.NewInt(0).SetBytes(input)

	base := big.NewInt(int64(len(b58Alphabet)))
	zero := big.NewInt(0)
	mod := &big.Int{}

	for x.Cmp(zero) != 0 {

		x.DivMod(x, base, mod)
		result = append(result, b58Alphabet[mod.Int64()])
	}

	JZ_ReverseBytes(result)
	for b := range input {

		if b == 0x00 {

			result = append([]byte{b58Alphabet[0]}, result...)
		} else {

			break
		}
	}

	return result
}

// Base58转字节数组，解密
func JZ_Base58Decode(input []byte) []byte {

	result := big.NewInt(0)
	zeroBytes := 0

	for b := range input {

		if b == 0x00 {

			zeroBytes++
		}
	}

	payload := input[zeroBytes:]
	for _, b := range payload {

		charIndex := bytes.IndexByte(b58Alphabet, b)
		result.Mul(result, big.NewInt(58))
		result.Add(result, big.NewInt(int64(charIndex)))
	}

	decoded := result.Bytes()
	//decoded...表示将decoded所有字节追加
	decoded = append(bytes.Repeat([]byte{byte(0x00)}, zeroBytes), decoded...)

	return decoded
}

