package blockchain

import (
	"testing"
)

func Test_NewMerkleTree(t *testing.T) {
	var ctbx = NewCoinbase([]byte(""), "coinbase transaction", 0)
	var transactions []*Transaction
	transactions = append(transactions, ctbx)

	NewMerkleTree(transactions)
}
