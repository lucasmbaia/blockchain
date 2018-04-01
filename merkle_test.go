package blockchain

import (
	"testing"
)

func Test_NewMerkleTree(t *testing.T) {
	transactions := [][]byte{{'A'}, {'B'}, {'C'}}

	NewMerkleTree(transactions)
}
