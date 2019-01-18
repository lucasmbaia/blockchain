package cli

import (
	//"encoding/hex"
	"fmt"
	"github.com/lucasmbaia/blockchain"
)

const (
	MULTIPLIER = 100000000
)

func (c *CLI) newTransaction(private, from, to string, value float64) {
	var infos = blockchain.Infos{
		Private:  private,
		From:	  from,
		To:	  to,
		Value:	  value,
	}

	var (
		body []byte
		err   error
	)

	if body, err = infos.Serialize(); err != nil {
		panic(err)
	}

	if err = transmit(gossip{
		Option:	"local_transaction",
		Body:	body,
	}); err != nil {
		panic(fmt.Sprintf("Deu ruim: %s\n", err.Error()))
	}
}
/*func (c *CLI) newTransaction(private, from, to string, value float64) {
	var (
		err	error
		amount	uint64
		w	*blockchain.Wallet
		valid	bool
		bc	*blockchain.Blockchain
		tx	*blockchain.Transaction
		txs	[]blockchain.Transaction
		g	gossip
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

	g = gossip {
		Option:	"transaction",
		Body:	tx.Serialize(),
	}

	if err = transmit(g); err != nil {
		panic(fmt.Sprintf("Error to transmit transaction: %s", err))
	}

	fmt.Printf("Transaction ID: %s\n", hex.EncodeToString(tx.ID[:]))
}*/
