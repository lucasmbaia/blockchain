package main

import (
    "github.com/lucasmbaia/blockchain"
    "fmt"
)

func generateWallet() {
    var (
	w   *blockchain.Wallet
	err error
    )

    if w, err = blockchain.NewWallet(); err != nil {
	panic(fmt.Sprintf("Error to generate new wallet: %s", err.Error()))
    }

    fmt.Printf("Private Key: %s\nAddress: %s\n", w.PrivateToHex(), string(w.Address))
}

func getBalance(private, address string) {
    var (
	err	error
	valid	bool
	bc	*blockchain.Blockchain
	txOut	[]blockchain.TXOutput
	amount	float64
    )

    if valid, _, err = blockchain.UnlockWallet(private, address); err != nil {
	panic(fmt.Sprintf("Error to unlock wallet: %s", err.Error()))
    }

    if !valid {
	panic(fmt.Sprintf("The private key is not allowed to unlock the wallet"))
    }

    bc = blockchain.NewBlockchain([]byte(address))

    if txOut, err = bc.FindUTXO([]byte(address)); err != nil {
	panic(fmt.Sprintf("Error to get balance: %s", err.Error()))
    }

    amount = 0.0
    for _, out := range txOut {
	amount += out.Value
    }

    fmt.Printf("Total of balance: %g\n", amount)
}

func main() {
    //generateWallet()
    getBalance("6704c183d5278523ad8a1eb88ba256ad1ea22222bda127fd5972f5acecdee835", "1CyssrDhEvZv2jXci6F5oueZwMXszm6kLs")
}
