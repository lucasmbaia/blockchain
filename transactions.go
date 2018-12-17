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
    TXid      []byte	//identificador da transacao
    Vout      int	//identifica uma saida especifica a ser gasta
    ScriptSig string	//fornece parametros de dados que satisfazem as condicoes do script pubkey
}

type TXOutput struct {
    Value        float64
    ScriptPubKey []byte
}

func (tx *Transaction) IsCoinbase() bool {
    return len(tx.TXInput) == 1 && len(tx.TXInput[0].TXid) == 0 && tx.TXInput[0].Vout == -1
}

func (txOut *TXOutput) Lock(address []byte) {
    var decoded = utils.Base58Decode(address)
    fmt.Println(decoded)
    txOut.ScriptPubKey = decoded[1:len(decoded)-4]
}

func (txOut *TXOutput) Unlock(pubHash []byte) bool {
    return bytes.Compare(txOut.ScriptPubKey, pubHash) == 0
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

func NewCoinbase(to []byte, data string) *Transaction {
    var hash utils.Hash

    if data == "" {
	data = fmt.Sprintf("Reward to '%s'", to)
    }

    var txin = TXInput{[]byte{}, -1, data}
    var txout = TXOutput{Value: float64(REWARD)}
    txout.Lock(to)

    var tx = Transaction{1, hash, 1, []TXInput{txin}, 1, []TXOutput{txout}}
    tx.SetID()

    return &tx
}
