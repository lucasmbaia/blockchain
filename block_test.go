package blockchain

import (
	"fmt"
	"testing"
)

func Test_NewBlockchain_Add_Block(t *testing.T) {
	bc := NewBlockchain()

	bc.AddBlock([]byte("Send 1 BTC to Lucas"))
	bc.AddBlock([]byte("Send 2 more BTC to Lucas"))

	for _, block := range bc.blocks {
		fmt.Printf("Index: %d\n", block.Index)
		fmt.Printf("Version Header: %d, Prev. Hash: %x, Bits: %d, Timestamp: %s, Nonce: %d\n", block.Header.Version, block.Header.PrevBlock, block.Header.Bits, block.Header.Timestamp, block.Header.Nonce)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Println()
	}
}
