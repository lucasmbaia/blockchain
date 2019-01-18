package blockchain

import (
	"fmt"
	"github.com/lucasmbaia/blockchain/utils"
	"math/big"
	"testing"
)

func Test_NewBlockchain_Add_Block(t *testing.T) {
	bc := NewBlockchain([]byte("1CyssrDhEvZv2jXci6F5oueZwMXszm6kLs"))

	for {
		bc.NewBlock([]byte("Send 1 BTC to Lucas"))
	}
	//bc.AddBlock([]byte("Send 2 more BTC to Lucas"))
}

func Test_Iterator_PrintBlochckain(t *testing.T) {
	var bc = NewBlockchain([]byte("1CyssrDhEvZv2jXci6F5oueZwMXszm6kLs"))
	var bci = bc.Iterator()
	var stop = big.NewInt(0)
	var count = 0

	for {
		block, _ := bci.Next()

		fmt.Printf("Block Index: %d, Block Hash: %x, Block Data: %s, Block Prev. Hash: %x, Block Bits: %d, Block Create: %s, Block Nonce: %d\n", block.Index, block.Hash, block.Data, block.Header.PrevBlock[:], block.Header.Bits, block.Header.Timestamp, block.Header.Nonce)

		fmt.Println(bc.ValidBlock(block))

		for _, transactions := range block.Transactions {
			fmt.Printf("Transaction ID: %x\n", transactions.ID)
		}

		if HashToBig(&block.Header.PrevBlock).Cmp(stop) == 0 {
			break
		}
		count++
	}

	fmt.Println(count)
}

func Test_UnspentTransaction(t *testing.T) {
	bc := NewBlockchain([]byte("1CyssrDhEvZv2jXci6F5oueZwMXszm6kLs"))
	pubHash := utils.AddressHashSPK([]byte("1CyssrDhEvZv2jXci6F5oueZwMXszm6kLs"))

	if tx, err := bc.UnspentTransaction(pubHash); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(tx)
	}
}

func Test_NewTransaction(t *testing.T) {
	bc := NewBlockchain([]byte("1CyssrDhEvZv2jXci6F5oueZwMXszm6kLs"))

	if transaction, err := bc.NewTransaction([]byte("1CyssrDhEvZv2jXci6F5oueZwMXszm6kLs"), []byte("18xB6w3WrNDriHruodzxsRFsiqYx6VufYY"), 100000); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(transaction)
	}
}

func Test_FindSpendable(t *testing.T) {
	bc := NewBlockchain([]byte("1CyssrDhEvZv2jXci6F5oueZwMXszm6kLs"))
	pubHash := utils.AddressHashSPK([]byte("1CyssrDhEvZv2jXci6F5oueZwMXszm6kLs"))

	if unspentOut, _, total, err := bc.FindSpendable(pubHash, 100000); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(unspentOut, total)
	}
}
