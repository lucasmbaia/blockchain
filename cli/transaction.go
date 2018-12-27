package main

import (
	"encoding/hex"
	"fmt"
	"github.com/lucasmbaia/blockchain"
)

const (
	MULTIPLIER = 100000000
)

func newTransaction(private, from, to string, value float64) {
	var (
		err    error
		amount uint64
		w      *blockchain.Wallet
		valid  bool
		bc     *blockchain.Blockchain
		tx     *blockchain.Transaction
		txs    []blockchain.Transaction
	)

	amount = uint64(value * MULTIPLIER)

	if valid, w, err = blockchain.UnlockWallet(private, from); err != nil {
		panic(fmt.Sprintf("Error to unlock waller: %s", err.Error()))
	}

	if !valid {
		panic(fmt.Sprintf("The private key is not allowed to unlock the wallet"))
	}

	bc = blockchain.NewBlockchain([]byte(from))

	if tx, err = bc.NewTransaction([]byte(from), []byte(to), amount); err != nil {
		panic(fmt.Sprintf("Error to generate transaction: %s", err.Error()))
	}

	for _, input := range tx.TXInput {
		_, transaction, _ := bc.FindTransaction(input.TXid)

		txs = append(txs, transaction)
	}

	if err = tx.SignTransaction(w, txs); err != nil {
		panic(fmt.Sprintf("Error to sign transaction: %s", err.Error()))
	}

	fmt.Printf("Transaction ID: %s\n", hex.EncodeToString(tx.ID[:]))

	fmt.Println(blockchain.ValidTransaction(tx, bc))
}

func main() {
	newTransaction("6704c183d5278523ad8a1eb88ba256ad1ea22222bda127fd5972f5acecdee835", "1CyssrDhEvZv2jXci6F5oueZwMXszm6kLs", "1G1yPBfSeLRY2sSzLGNbAUWRFBgSv58F4r", 0.01)
}
