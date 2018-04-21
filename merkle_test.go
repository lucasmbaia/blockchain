package blockchain

import (
	"testing"
)

func Test_NewMerkleTree(t *testing.T) {
	var ctbx = NewCoinbase("", "coinbase transaction")
	var transactions []*Transaction
	transactions = append(transactions, ctbx)

	NewMerkleTree(transactions)
}
