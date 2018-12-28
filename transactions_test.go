package blockchain

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/lucasmbaia/blockchain/utils"
	"testing"
)

const (
	TXID_TRANSACTION_UTXO = "16ee4aa4ad7787fa37936ad442e4f3e3d1f1a28d9f68461f20f09baa2d8d4f08"
	PRIVATE_KEY           = "6704c183d5278523ad8a1eb88ba256ad1ea22222bda127fd5972f5acecdee835"
	ADDRESS               = "1CyssrDhEvZv2jXci6F5oueZwMXszm6kLs"
)

func Test_RawTransaction(t *testing.T) {
	var (
		txID     utils.Hash
		prevTxID []byte
		err      error
	)

	if prevTxID, err = hex.DecodeString(TXID_TRANSACTION_UTXO); err != nil {
		t.Fatal(err)
	}
	copy(txID[:], prevTxID)

	var transaction = &Transaction{
		Version:      1,
		TXInputCount: 1,
		TXInput: []TXInput{
			{Coinbase: "03250507", Vout: -1, Sequence: "00000000"},
		},
		TXOutputCount: 1,
		TXOutput: []TXOutput{
			{Address: []byte("1G1yPBfSeLRY2sSzLGNbAUWRFBgSv58F4r"), Value: 600000},
			//{Address: []byte("1CyssrDhEvZv2jXci6F5oueZwMXszm6kLs"), Value: 390000},
		},
	}

	fmt.Println(transaction.RawTransaction())
}

func Test_SignTransaction(t *testing.T) {
	var (
		w        *Wallet
		txID     utils.Hash
		prevTxID []byte
		err      error
		valid    bool
	)

	if valid, w, err = UnlockWallet(PRIVATE_KEY, ADDRESS); err != nil {
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
		TXInputCount: 1,
		TXInput: []TXInput{
			{TXid: txID, Vout: 0},
		},
		TXOutputCount: 2,
		TXOutput: []TXOutput{
			{Address: []byte("1G1yPBfSeLRY2sSzLGNbAUWRFBgSv58F4r"), Value: 600000},
			{Address: []byte("1CyssrDhEvZv2jXci6F5oueZwMXszm6kLs"), Value: 390000},
		},
	}

	var tx = Transaction{
		Version:       1,
		ID:            txID,
		TXOutputCount: 1,
		TXOutput: []TXOutput{
			{Index: 0, Address: []byte("1CyssrDhEvZv2jXci6F5oueZwMXszm6kLs"), Value: 2500000000},
		},
	}

	transaction.SignTransaction(w, []Transaction{tx})
}

func Test_ParseScriptSig(t *testing.T) {
	fmt.Println(parseScriptSig("473044022051377fe821e8c353ea848a4b573ddaa51cb63c32b588e0a5b5601e76a7becfb9022037341d11a1b8187dec4efe729df635e56558ea5d1ad175bd615448dfd59badc2010140f380c9afcf9e392bd4f7a0e68aff98ff253d0950720b4fc3185db1fe43f6ec9dc118d8d0067ff731b0a466d6f35edc6400b16f2df567afe56c27fa555c1cccd3"))
}

func Test_ParseP2pkh(t *testing.T) {
	fmt.Println(parseSignature("304402200513c67f2d46c0f2494e1c20b7d3100e006826a91578f5ab2028feb93266b9d80220f4edd33df15925ea10f851c6007ddee332852164d89c246080483901aa7d513601"))
}
