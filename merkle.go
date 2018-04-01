package blockchain

import (
	"github.com/lucasmbaia/blockchain/utils"
)

type MerkleRoot struct {
	*MerkleNode
}

type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Hash  utils.Hash
}

func NewMerkleTree(transactions [][]byte) *MerkleRoot {
	var (
		nodes []MerkleNode
	)

	if len(transactions) == 1 {
		var hash utils.Hash
		copy(hash[:], transactions[0])

		nodes = append(nodes, MerkleNode{Hash: hash})
		return &MerkleRoot{&nodes[0]}
	}

	if len(transactions)%2 != 0 {
		transactions = append(transactions, transactions[len(transactions)-1])
	}

	for _, transaction := range transactions {
		nodes = append(nodes, MerkleNode{Hash: utils.CalcDoubleHash(transaction)})
	}

	for i := 0; i < len(transactions)/2; i++ {
		var newNodes []MerkleNode

		for j := 0; j < len(nodes); j += 2 {
			var hash [utils.HASH_SIZE * 2]byte
			copy(hash[:utils.HASH_SIZE], nodes[j].Hash[:])
			copy(hash[utils.HASH_SIZE:], nodes[j+1].Hash[:])

			doubleHash := utils.CalcDoubleHash(hash[:])

			newNodes = append(newNodes, MerkleNode{
				Left:  &nodes[j],
				Right: &nodes[j+1],
				Hash:  doubleHash,
			})
		}

		nodes = newNodes
	}

	return &MerkleRoot{&nodes[0]}
}
