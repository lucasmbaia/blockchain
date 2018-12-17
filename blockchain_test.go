package blockchain

import (
	"fmt"
	"math/big"
	"testing"
)

func Test_NewBlockchain_Add_Block(t *testing.T) {
	bc := NewBlockchain([]byte("1CyssrDhEvZv2jXci6F5oueZwMXszm6kLs"))

	bc.AddBlock([]byte("Send 1 BTC to Lucas"))
	bc.AddBlock([]byte("Send 2 more BTC to Lucas"))
}

func Test_Iterator_PrintBlochckain(t *testing.T) {
	var bc = NewBlockchain([]byte("18xB6w3WrNDriHruodzxsRFsiqYx6VufYY"))
	var bci = bc.Iterator()
	var stop = big.NewInt(0)

	for {
		block, _ := bci.Next()

		fmt.Printf("Block Index: %d, Block Hash: %x, Block Data: %s, Block Prev. Hash: %x, Block Bits: %d, Block Create: %s, Block Nonce: %d\n", block.Index, block.Hash, block.Data, block.Header.PrevBlock[:], block.Header.Bits, block.Header.Timestamp, block.Header.Nonce)

		if HashToBig(&block.Header.PrevBlock).Cmp(stop) == 0 {
			break
		}
	}
}

func Test_UnspentTransaction(t *testing.T) {
  bc := NewBlockchain([]byte("18xB6w3WrNDriHruodzxsRFsiqYx6VufYY"))

  if tx, err := bc.UnspentTransaction([]byte("18xB6w3WrNDriHruodzxsRFsiqYx6VufYY")); err != nil {
    t.Fatal(err)
  } else {
    fmt.Println(tx)
  }
}
