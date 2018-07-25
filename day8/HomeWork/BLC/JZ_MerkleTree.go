package BLC

import "crypto/sha256"

//默克尔树
type JZ_MerkleTree struct {
	//根节点
	JZ_RootNode *JZ_MerkleNode
}

//默克尔树节点
type JZ_MerkleNode struct {
	//做节点
	JZ_Left *JZ_MerkleNode
	//右节点
	JZ_Right *JZ_MerkleNode
	//节点数据
	JZ_Data []byte
}

//新建节点
func JZ_NewMerkleNode(left, right *JZ_MerkleNode, data []byte) *JZ_MerkleNode {

	mNode := JZ_MerkleNode{}

	if left == nil && right == nil {

		hash := sha256.Sum256(data)
		mNode.JZ_Data = hash[:]
	} else {

		prevHashes := append(left.JZ_Data, right.JZ_Data...)
		hash := sha256.Sum256(prevHashes)
		mNode.JZ_Data = hash[:]
	}

	mNode.JZ_Left = left
	mNode.JZ_Right = right

	return &mNode
}

// 1 2 3 --> 1 2 3 3
//新建默克尔树
func JZ_NewMerkleTree(datas [][]byte) *JZ_MerkleTree {

	var nodes []*JZ_MerkleNode

	//如果是奇数，添加最后一个交易哈希拼凑为偶数个交易
	if len(datas) % 2 != 0 {

		datas = append(datas, datas[len(datas)-1])
	}

	//将每一个交易哈希构造为默克尔树节点
	for _, data := range datas {

		node := JZ_NewMerkleNode(nil, nil, data)
		nodes = append(nodes, node)
	}

	//将所有节点两两组合生成新节点，直到最后只有一个更节点
	for i := 0; i < len(datas)/2; i++ {

		var newLevel []*JZ_MerkleNode

		for j := 0; j < len(nodes); j += 2 {

			node := JZ_NewMerkleNode(nodes[j], nodes[j+1], nil)
			newLevel = append(newLevel, node)
		}

		nodes = newLevel
	}

	//取根节点返回
	mTree := JZ_MerkleTree{nodes[0]}

	return &mTree
}


