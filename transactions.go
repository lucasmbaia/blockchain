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
    "reflect"
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
    TXid      utils.Hash  //identificador da transacao
    Vout      int	  //identifica uma saida especifica a ser gasta
    ScriptSig ScriptSig	  //fornece parametros de dados que satisfazem as condicoes do script pubkey
}

type ScriptSig struct {
    Asm	string
    Hex	string
}

type TXOutput struct {
    Index	  int
    Value	  uint64
    ScriptPubKey  ScriptPubKey
    Address	  []byte
}

type ScriptPubKey struct {
    Asm	string
    Hex	string
}

func (tx *Transaction) Serialize() []byte {
    var result bytes.Buffer
    var encoder *gob.Encoder = gob.NewEncoder(&result)

    if err := encoder.Encode(tx); err != nil {
	log.Printf("Error to serialize transaction: %s\n", err)
    }

    return result.Bytes()
}

func DeserializeTransaction(t []byte) *Transaction {
    var transaction Transaction
    var decoder *gob.Decoder = gob.NewDecoder(bytes.NewReader(t))

    if err := decoder.Decode(&transaction); err != nil {
	log.Printf("Error to deserialize: %s\n", err)
    }

    return &transaction
}

func (tx *Transaction) IsCoinbase() bool {
    return len(tx.TXInput) == 1 && len(tx.TXInput[0].TXid) == 0 && tx.TXInput[0].Vout == -1
}

func (txOut *TXOutput) Lock(address []byte) {
    //txOut.ScriptPubKey = utils.AddressHashSPK(address)
    //txOut.Address = address
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
	raw		[]string
	inputs		string
	outputs		string
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

    /*raw = append(raw, utils.ConvertUnsigned4Bytes(uint32(tx.Version)))  //four-byte version
    raw = append(raw, inputs) //total inputs of transaction

    //fazer a magica nova aqui porra a partir dos inputs
    raw = append(raw, utils.ReverseHash(tx.ID[:])) //reverse transaction the redeem an output
    raw = append(raw, utils.ConvertUnsigned4Bytes(uint32(tx.TXOutputCount) - 1)) //four-byte of number transaction outputs "trocar essa parada aqui, tem que referenciar o numero do output"
    raw = append(raw, "ffffffff") //four-byte field denoting the sequence
    raw = append(raw, outputs) //one-byte number of outputs

    for _, out := range tx.TXOutput {
	raw = append(raw, utils.ConvertUnsigned8Bytes(out.Value)) //eigth-byte of amount transfer
	raw = append(raw, "19") //lenght of scriptSig
	raw = append(raw, out.ScriptPubKey) //address hash of scriptPubKey
    }

    raw = append(raw, "00000000") //lock-time field*/

    raw = append(raw, utils.ConvertUnsigned4Bytes(uint32(tx.Version)))  //four-byte version
    raw = append(raw, inputs) //total inputs of transaction

    for _, input := range tx.TXInput {
	raw = append(raw, utils.ReverseHash(input.TXid[:]))
	raw = append(raw, utils.ConvertUnsigned4Bytes(uint32(input.Vout)))

	if (ScriptSig{}) != input.ScriptSig {
	    if input.ScriptSig.Hex != "" {
		raw = append(raw, input.ScriptSig.Hex)
	    }
	}

	raw = append(raw, "ffffffff") //four-byte field denoting the sequence
    }

    raw = append(raw, outputs) //one-byte number of outputs

    for _, output := range tx.TXOutput {
	raw = append(raw, utils.ConvertUnsigned8Bytes(output.Value)) //eigth-byte of amount transfer
	raw = append(raw, "19") //lenght of scriptSig
	raw = append(raw,  utils.AddressHashSPK(output.Address)) //address hash of scriptPubKey
    }

    raw = append(raw, "00000000") //lock-time field*/

    return strings.Join(raw, "")
}

func (tx *Transaction) SignTransaction(w *Wallet, txs []Transaction) error {
    var (
	raw	    string
	d2sha256    utils.Hash
	data	    []byte
	err	    error
	signature   *btcec.Signature
	//scriptSig   []byte
	key	    = (*btcec.PrivateKey)(&w.PrivateKey)
	//sign	    []string
	//inputs	    string
	//outputs	    string
	//porra	    []string
	hashPubKey	[]byte
    )

    for idx, out := range tx.TXOutput {
	if hashPubKey, err = HashPubKey(out.Address); err != nil {
	    return err
	}

	tx.TXOutput[idx].ScriptPubKey = ScriptPubKey {
	    Asm:  fmt.Sprintf("OP_DUP OP_HASH160 %s OP_EQUALVERIFY OP_CHECKSIG", hex.EncodeToString(hashPubKey)),
	    Hex:  utils.AddressHashSPK(out.Address),
	}
    }

    for idx, input := range tx.TXInput {
	var transaction Transaction

	for _, txin := range txs {
	    if bytes.Compare(input.TXid[:], txin.ID[:]) == 0 {
		transaction = txin
	    }
	}

	if !reflect.DeepEqual(transaction, (Transaction{})) {
	    raw = transaction.RawTransaction()
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

	    var asmSignature []string
	    var hexSignature []string

	    var sign = append(signature.Serialize(), byte(SigHashAll))
	    asmSignature = append(asmSignature, hex.EncodeToString(sign))

	    var bytesSignature = strconv.FormatInt(int64(len(asmSignature[0]) / 2), 16);
	    hexSignature = append(hexSignature, bytesSignature)
	    hexSignature = append(hexSignature, asmSignature...)
	    asmSignature = append(asmSignature, "[ALL] ")
	    asmSignature = append(asmSignature, hex.EncodeToString(w.PublicKey))
	    hexSignature = append(hexSignature, "01")

	    var bytesPubKey = strconv.FormatInt(int64(len(hex.EncodeToString(w.PublicKey)) / 2), 16)

	    hexSignature = append(hexSignature, bytesPubKey)
	    hexSignature = append(hexSignature, hex.EncodeToString(w.PublicKey))

	    tx.TXInput[idx].ScriptSig = ScriptSig{strings.Join(asmSignature, ""), strings.Join(hexSignature, "")}
	    /*scriptSig = append(scriptSig, byte(0x4))
	    scriptSig = append(scriptSig, w.PublicKey...)

	    porra = append(porra, "01")
	    bytesHash := strconv.FormatInt(int64(len(hex.EncodeToString(w.PublicKey)) / 2), 16)
	    porra = append(porra, bytesHash)
	    porra = append(porra, hex.EncodeToString(w.PublicKey))*/
	}
    }

    /*inputs = strconv.FormatInt(int64(tx.TXInputCount), 16)
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

    fmt.Println("RAW: ", raw)
    fmt.Println("Raw Hash: ", hex.EncodeToString(d2sha256[:]))
    fmt.Println("Script SIG: ", hex.EncodeToString(scriptSig))
    fmt.Println("ALL: ", string(strings.Join(sign, "")))*/
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

func unlock(scriptSig string, w *Wallet) {
    var (
	signature string
	pubkey	  string
	stack	  []string
    )

    signature, pubkey, _ = parseScriptSig(scriptSig)

    decoded, _ := hex.DecodeString(pubkey)
    h1, _ := HashPubKey(decoded)
    h2, _ := HashPubKey(w.PublicKey)

    stack = append(stack, signature) //sig stack
    stack = append(stack, pubkey) //pubKey stack
    //stack = append(stack, HashPubKey([]byte(pubKey)))
    //stack = append(stack, HashPubKey(w.PublicKey))

    fmt.Println(stack)
    fmt.Println(hex.EncodeToString(h1))
    fmt.Println(hex.EncodeToString(h2))

    fmt.Println(hex.EncodeToString(w.PublicKey))
    fmt.Println(pubkey)
    fmt.Println(scriptSig)
}

func parseScriptSig(scriptSig string) (string, string, error) {
    var (
	signature     string
	sizeSignature int64
	err	      error
	cutePublicKey string
    )

    if sizeSignature, err = strconv.ParseInt(fmt.Sprintf("0x%s", scriptSig[:2]), 0, 64); err != nil {
	return signature, cutePublicKey, err
    }

    signature = scriptSig[2:sizeSignature*2+2]
    cutePublicKey = scriptSig[sizeSignature*2+2:]
    return signature, cutePublicKey[4:], nil
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
