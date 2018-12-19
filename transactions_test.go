package blockchain

import (
    "testing"
    "fmt"
    "github.com/lucasmbaia/blockchain/utils"
    "encoding/hex"
    "errors"
)

const (
    TXID_TRANSACTION_UTXO = "01ecaecc96b148589be10ef3f8fffc70dcc14970ed77144a787c087dfcd0b5e2"
    PRIVATE_KEY		  = "6704c183d5278523ad8a1eb88ba256ad1ea22222bda127fd5972f5acecdee835"
    ADDRESS		  = "1CyssrDhEvZv2jXci6F5oueZwMXszm6kLs"
)

func Test_RawTransaction(t *testing.T) {
    var (
	txID	  utils.Hash
	prevTxID  []byte
	err	  error
    )

    if prevTxID, err = hex.DecodeString(TXID_TRANSACTION_UTXO); err != nil {
	t.Fatal(err)
    }
    copy(txID[:], prevTxID)

    var transaction = &Transaction{
	Version:      1,
	ID:	      txID,
	TXInputCount: 1,
	TXOutputCount:	2,
	TXOutput:	[]TXOutput{
	    {Address: []byte("1N8A7WLw4TA8Vbu6GQ9JSotbvtTPY5uKhf"), Value: 600000},
	    {Address: []byte("1EoPUX89wRFGHRNGJFTzqfNcRgdzgs5JhK"), Value: 390000},
	},
    }

    fmt.Println(transaction.RawTransaction())
}

func Test_SignTransaction(t *testing.T) {
    var (
	w     *Wallet
	txID	  utils.Hash
	prevTxID  []byte
	err	  error
	valid	  bool
    )

    if valid , w, err = UnlockWallet(PRIVATE_KEY, ADDRESS); err != nil {
	t.Fatal(err)
    }

    if !valid {
	t.Fatal(errors.New("Invalid"))
    }

    if prevTxID, err = hex.DecodeString(TXID_TRANSACTION_UTXO); err != nil {
	t.Fatal(err)
    }
    copy(txID[:], prevTxID)

    var transaction = &Transaction{
	Version:      1,
	ID:	      txID,
	TXInputCount: 1,
	TXOutputCount:	2,
	TXOutput:	[]TXOutput{
	    {Address: []byte("1N8A7WLw4TA8Vbu6GQ9JSotbvtTPY5uKhf"), Value: 600000},
	    {Address: []byte("1EoPUX89wRFGHRNGJFTzqfNcRgdzgs5JhK"), Value: 390000},
	},
    }

   transaction.SignTransaction(w) 
}
