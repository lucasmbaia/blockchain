package blockchain

import (
    "bytes"
    "encoding/gob"
    "fmt"
    "github.com/lucasmbaia/blockchain/utils"
    "github.com/btcsuite/btcd/btcec"
    "log"
    "strconv"
    "strings"
    "encoding/hex"
    //"crypto/ecdsa"
    //"crypto/rand"
    //"math/big"
    //"time"
)

type SigHashType uint32
const (
    REWARD	uint64 = 2500000000
    SigHashAll	SigHashType = 0x1
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
    Value	  uint64
    ScriptPubKey  string
    Address	  []byte
}

func (tx *Transaction) IsCoinbase() bool {
    return len(tx.TXInput) == 1 && len(tx.TXInput[0].TXid) == 0 && tx.TXInput[0].Vout == -1
}

func (txOut *TXOutput) Lock(address []byte) {
    txOut.ScriptPubKey = utils.AddressHashSPK(address)
    txOut.Address = address
}

func (txOut *TXOutput) Unlock(pubHash []byte) bool {
    return true
    //return bytes.Compare(txOut.ScriptPubKey, pubHash) == 0
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

func (tx *Transaction) RawTransaction() string {
    var (
	raw     []string
	inputs  string
	outputs string
    )

    inputs = strconv.FormatInt(int64(tx.TXInputCount), 16)
    if len(inputs) == 1 {
	inputs = fmt.Sprintf("0%s", inputs)
    }

    outputs = strconv.FormatInt(int64(tx.TXOutputCount), 16)
    if len(outputs) == 1 {
	outputs = fmt.Sprintf("0%s", outputs)
    }

    /*raw = append(raw, utils.ConvertUnsigned4Bytes(uint32(tx.Version)))  //four-byte version
    raw = append(raw, inputs) //total inputs of transaction
    raw = append(raw, utils.ReverseHash(tx.ID[:])) //reverse transaction the redeem an output
    raw = append(raw, utils.ConvertUnsigned4Bytes(uint32(tx.TXOutputCount) - 1)) //four-byte of number transaction outputs
    raw = append(raw, "19") // lenght of scriptSig
    raw = append(raw, utils.AddressHashSPK(tx.TXOutput[0].Address)) //address hash of scriptPubKey
    raw = append(raw, "ffffffff") //four-byte field denoting the sequence
    raw = append(raw, outputs) //one-byte number of outputs
    raw = append(raw, utils.ConvertUnsigned8Bytes(amount)) //eigth-byte of amount transfer
    raw = append(raw, "19")
    raw = append(raw, utils.AddressHashSPK(tx.TXOutput[1].Address))
    raw = append(raw, "00000000") //lock-time field
    raw = append(raw, "01000000") //four-byte hash code type*/

    raw = append(raw, utils.ConvertUnsigned4Bytes(uint32(tx.Version)))  //four-byte version
    raw = append(raw, inputs) //total inputs of transaction
    raw = append(raw, utils.ReverseHash(tx.ID[:])) //reverse transaction the redeem an output
    raw = append(raw, utils.ConvertUnsigned4Bytes(uint32(tx.TXOutputCount) - 1)) //four-byte of number transaction outputs
    raw = append(raw, "ffffffff") //four-byte field denoting the sequence
    raw = append(raw, outputs) //one-byte number of outputs

    for _, out := range tx.TXOutput {
	raw = append(raw, utils.ConvertUnsigned8Bytes(out.Value)) //eigth-byte of amount transfer
	raw = append(raw, "19") //lenght of scriptSig
	raw = append(raw, out.ScriptPubKey) //address hash of scriptPubKey
    }

    raw = append(raw, "00000000") //lock-time field
    //raw = append(raw, "01000000") //four-byte hash code type

    return strings.Join(raw, "")
}

func (tx *Transaction) SignTransaction(w *Wallet) error {
    var (
	raw	    string
	d2sha256    utils.Hash
	data	    []byte
	err	    error
	signature   *btcec.Signature
	scriptSig   []byte
	key	    = (*btcec.PrivateKey)(&w.PrivateKey)
	sign	    []string
	inputs	    string
	outputs	    string
    )

    for idx, out := range tx.TXOutput {
      tx.TXOutput[idx].ScriptPubKey = utils.AddressHashSPK(out.Address)
    }

    raw = tx.RawTransaction()
    if data, err = hex.DecodeString(raw); err != nil {
	return err
    }
    d2sha256 = utils.CalcDoubleHash(data)

    /*if r, s, err = ecdsa.Sign(rand.Reader, &w.PrivateKey, d2sha256[:]); err != nil {
	return err
    }

    signature = btcec.Signature{R: r, S: s}
    sign = append(signature.Serialize(), byte(SigHashAll))
    sign = append(sign, w.PublicKey...)*/

    if signature, err = key.Sign(d2sha256[:]); err != nil {
	return err
    }

    scriptSig = append(signature.Serialize(), byte(SigHashAll))
    fmt.Println(hex.EncodeToString(signature.Serialize()))
    scriptSig = append(scriptSig, byte(0x4))
    scriptSig = append(scriptSig, w.PublicKey...)
    fmt.Println(hex.EncodeToString(scriptSig))

    inputs = strconv.FormatInt(int64(tx.TXInputCount), 16)
    if len(inputs) == 1 {
	inputs = fmt.Sprintf("0%s", inputs)
    }

    outputs = strconv.FormatInt(int64(tx.TXOutputCount), 16)
    if len(outputs) == 1 {
	outputs = fmt.Sprintf("0%s", outputs)
    }

    sign = append(sign, utils.ConvertUnsigned4Bytes(uint32(tx.Version)))  //four-byte version
    sign = append(sign, inputs) //total inputs of transaction
    sign = append(sign, utils.ReverseHash(tx.ID[:])) //reverse transaction the redeem an output
    sign = append(sign, utils.ConvertUnsigned4Bytes(uint32(tx.TXOutputCount) - 1)) //four-byte of number transaction outputs
    sign = append(sign, strconv.FormatInt(int64(len(scriptSig)), 16))
    sign = append(sign, hex.EncodeToString(scriptSig))
    sign = append(sign, "ffffffff") //four-byte field denoting the sequence
    sign = append(sign, outputs) //one-byte number of outputs

    for _, out := range tx.TXOutput[1:] {
	sign = append(sign, utils.ConvertUnsigned8Bytes(out.Value)) //eigth-byte of amount transfer
	sign = append(sign, "19") //lenght of scriptSig
	sign = append(sign, utils.AddressHashSPK(out.Address)) //address hash of scriptPubKey
    }

    sign = append(sign, "00000000") //lock-time field

    return nil
}

/*func (tx *Transaction) Unlock(w *Wallet, address []byte) bool {
    var (
	signature string
	pubkey	  string
    )

    for _, input := range tx.TXInput {
      signature, pubkey = parseScriptSig(input.ScriptSig)
    }
}*/

func parseScriptSig(scriptSig string) (string, string) {
    return scriptSig[:142], scriptSig[146:]
}

func NewCoinbase(to []byte, data string) *Transaction {
    var hash utils.Hash

    if data == "" {
	data = fmt.Sprintf("Reward to '%s'", to)
    }

    var txout = TXOutput{Value: REWARD}
    txout.Lock(to)

    var tx = Transaction{1, hash, 0, []TXInput{}, 1, []TXOutput{txout}}
    tx.SetID()

    return &tx
}
