package blockchain

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/lucasmbaia/blockchain/utils"
	"log"
	//"time"
)

const (
	REWARD = 25
)

type Transaction struct {
	Version       int32
	ID            utils.Hash
	TXInputCount  int32
	TXInput       []TXInput
	TXOutputCount int32
	TXOutput      []TXOutput
	//LockTime      time.Time
}

type TXInput struct {
	TXid      []byte
	Vout      int
	ScriptSig string
}

type TXOutput struct {
	Value        float64
	ScriptPubKey string
}

func (tx *Transaction) SetID() {
	var (
		encoded bytes.Buffer
		err     error
	)

	var encoder = gob.NewEncoder(&encoded)
	if err = encoder.Encode(tx); err != nil {
		log.Fatalf("Error to SetID: %s\n", err)
	}

	tx.ID = utils.CalcDoubleHash(encoded.Bytes())
}

func NewCoinbase(to, data string) *Transaction {
	var hash utils.Hash

	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	var txin = TXInput{[]byte{}, -1, data}
	var txout = TXOutput{float64(REWARD), to}
	var tx = Transaction{1, hash, 1, []TXInput{txin}, 1, []TXOutput{txout}}
	tx.SetID()

	return &tx
}
