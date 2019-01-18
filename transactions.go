package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/lucasmbaia/blockchain/utils"
	"log"
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"time"
	"sync"
)

type SigHashType uint32

const (
	REWARD     uint64      = 2500000000
	SigHashAll SigHashType = 0x1
)

var (
  mutex			  = &sync.RWMutex{}
  UnprocessedTransactions = make(map[utils.Hash]*Transaction)
)

type Transaction struct {
	Version       int32
	ID            utils.Hash
	Created       time.Time
	TXInputCount  int32
	TXInput       []TXInput
	TXOutputCount int32
	TXOutput      []TXOutput
	//LockTime      time.Time
}

type TXInput struct {
	TXid      utils.Hash //identificador da transacao
	Vout      int        //identifica uma saida especifica a ser gasta
	ScriptSig ScriptSig  //fornece parametros de dados que satisfazem as condicoes do script pubkey
	Coinbase  string
	Sequence  string
}

type ScriptSig struct {
	Asm string
	Hex string
}

type TXOutput struct {
	Index        int
	Value        uint64
	ScriptPubKey ScriptPubKey
	Address      []byte
}

type ScriptPubKey struct {
	Asm string
	Hex string
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
	return len(tx.TXInput) == 1 && tx.TXInput[0].Vout == -1
}

func (txOut *TXOutput) Lock(address []byte) {
	hashPubKey := AddressToPubHash(address)

	txOut.ScriptPubKey = ScriptPubKey{
		Asm: fmt.Sprintf("OP_DUP OP_HASH160 %s OP_EQUALVERIFY OP_CHECKSIG", hex.EncodeToString(hashPubKey)),
		Hex: utils.AddressHashSPK(address),
	}
}

func (txOut *TXOutput) Unlock(pubHash string) bool {
	return strings.Compare(txOut.ScriptPubKey.Hex, pubHash) == 0
}

func (tx *Transaction) SetID() {
	var hash = utils.CalcDoubleHash(tx.Serialize())
	tx.ID = hash
}

func (tx *Transaction) RawTransaction() string {
	var (
		raw     []string
		inputs  string
		outputs string
	)

	if tx.TXInputCount == 0 {
		inputs = "00"
	} else {
		inputs = strconv.FormatInt(int64(tx.TXInputCount), 16)
		if len(inputs) == 1 {
			inputs = fmt.Sprintf("0%s", inputs)
		}
	}

	if tx.TXOutputCount == 0 {
		outputs = "00"
	} else {
		outputs = strconv.FormatInt(int64(tx.TXOutputCount), 16)
		if len(outputs) == 1 {
			outputs = fmt.Sprintf("0%s", outputs)
		}
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

	raw = append(raw, utils.ConvertUnsigned4Bytes(uint32(tx.Version))) //four-byte version
	raw = append(raw, inputs)                                          //total inputs of transaction

	for _, input := range tx.TXInput {
		if strings.Compare(hex.EncodeToString(input.TXid[:]), "0000000000000000000000000000000000000000000000000000000000000000") == 0 {
			raw = append(raw, "0000000000000000000000000000000000000000000000000000000000000000")
			raw = append(raw, utils.ConvertUnsigned4Bytes(uint32(input.Vout)))
			raw = append(raw, fmt.Sprintf("04%s", input.Coinbase))
		} else {
			raw = append(raw, utils.ReverseHash(input.TXid[:]))
			raw = append(raw, utils.ConvertUnsigned4Bytes(uint32(input.Vout)))
		}

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
		raw = append(raw, "19")                                      //lenght of scriptSig
		raw = append(raw, utils.AddressHashSPK(output.Address))      //address hash of scriptPubKey
	}

	raw = append(raw, "00000000") //lock-time field*/

	return strings.Join(raw, "")
}

func p2pkh(r *big.Int, s *big.Int) string {
	var (
		sString         string
		rString         string
		signature       []string
		signatureLength int
	)

	rString = hex.EncodeToString(r.Bytes())
	sString = hex.EncodeToString(s.Bytes())

	signatureLength = ((len(rString) + len(sString)) / 2) + 4 // size of r, size of s and +4 bytes

	signature = append(signature, "30")                                          //DER signature marker
	signature = append(signature, strconv.FormatInt(int64(signatureLength), 16)) //declares signature lenght in bytes
	signature = append(signature, "02")                                          //r value maker
	signature = append(signature, strconv.FormatInt(int64(len(rString)/2), 16))  //r lenght in bytes
	signature = append(signature, rString)                                       //r value
	signature = append(signature, "02")                                          //s value maker
	signature = append(signature, strconv.FormatInt(int64(len(sString)/2), 16))  //s lenght in bytes
	signature = append(signature, sString)                                       //s value
	signature = append(signature, "01")                                          //SIGHASH_ALL

	return strings.Join(signature, "")
}

func (tx *Transaction) SignTransaction(w *Wallet, txs []Transaction) error {
	var (
		raw        string
		d2sha256   utils.Hash
		data       []byte
		err        error
		hashPubKey []byte
		r          *big.Int
		s          *big.Int
	)

	for idx, out := range tx.TXOutput {
		if reflect.DeepEqual(out, (TXOutput{})) {
			if hashPubKey, err = HashPubKey(out.Address); err != nil {
				return err
			}

			tx.TXOutput[idx].ScriptPubKey = ScriptPubKey{
				Asm: fmt.Sprintf("OP_DUP OP_HASH160 %s OP_EQUALVERIFY OP_CHECKSIG", hex.EncodeToString(hashPubKey)),
				Hex: utils.AddressHashSPK(out.Address),
			}
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

			/****** NEW CODE ******/
			if r, s, err = ecdsa.Sign(rand.Reader, &w.PrivateKey, d2sha256[:]); err != nil {
				return err
			}
			/***** END NEW CODE ******/

			/*if signature, err = key.Sign(d2sha256[:]); err != nil {
				return err
			}*/

			var asmSignature []string
			var hexSignature []string

			//var sign = append(signature.Serialize(), byte(SigHashAll))
			//asmSignature = append(asmSignature, hex.EncodeToString(sign))
			asmSignature = append(asmSignature, p2pkh(r, s))

			var bytesSignature = strconv.FormatInt(int64(len(asmSignature[0])/2), 16)
			hexSignature = append(hexSignature, bytesSignature)
			hexSignature = append(hexSignature, asmSignature...)
			asmSignature = append(asmSignature, "[ALL] ")
			asmSignature = append(asmSignature, hex.EncodeToString(w.PublicKey))
			hexSignature = append(hexSignature, "01")

			var bytesPubKey = strconv.FormatInt(int64(len(hex.EncodeToString(w.PublicKey))/2), 16)

			hexSignature = append(hexSignature, bytesPubKey)
			hexSignature = append(hexSignature, hex.EncodeToString(w.PublicKey))

			tx.TXInput[idx].ScriptSig = ScriptSig{strings.Join(asmSignature, ""), strings.Join(hexSignature, "")}
		}
	}

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

func (tx *Transaction) ValidTransaction(transactions map[utils.Hash]*Transaction) error {
	var (
		err         error
		signature   string
		pubkey      string
		transaction *Transaction
		decoded     []byte
		hpk         []byte
		hspk        string
		curve       elliptic.Curve
		pubKeyLen   int
		raw         string
		data        []byte
		d2sha256    utils.Hash
	)

	curve = elliptic.P256()

	for _, input := range tx.TXInput {
		transaction = transactions[input.TXid]
		if signature, pubkey, err = parseScriptSig(input.ScriptSig.Hex); err != nil {
			return err
		}

		hspk = strings.Replace(transaction.TXOutput[input.Vout].ScriptPubKey.Asm, "OP_DUP OP_HASH160 ", "", -1)
		hspk = strings.Replace(hspk, " OP_EQUALVERIFY OP_CHECKSIG", "", -1)

		if decoded, err = hex.DecodeString(pubkey); err != nil {
			return err
		}

		if hpk, err = HashPubKey(decoded); err != nil {
			return err
		}

		if strings.Compare(hex.EncodeToString(hpk), hspk) != 0 {
			return errors.New("OP_EQUALVERIFY is invalid")
		}

		var decodePubKey []byte
		if decodePubKey, err = hex.DecodeString(pubkey); err != nil {
			return err
		}

		pubKeyLen = len(decodePubKey)

		var r = &big.Int{}
		var s = &big.Int{}

		if r, s, err = parseSignature(signature); err != nil {
			return err
		}

		var x = big.Int{}
		var y = big.Int{}
		x.SetBytes(decodePubKey[:(pubKeyLen / 2)])
		y.SetBytes(decodePubKey[(pubKeyLen / 2):])

		var rawPubKey = ecdsa.PublicKey{curve, &x, &y}

		raw = transaction.RawTransaction()
		if data, err = hex.DecodeString(raw); err != nil {
			return err
		}

		d2sha256 = utils.CalcDoubleHash(data)

		if !ecdsa.Verify(&rawPubKey, d2sha256[:], r, s) {
			return errors.New("CHECKSIG is invalid")
		}
	}

	return nil
}

func parseScriptSig(scriptSig string) (string, string, error) {
	var (
		signature     string
		sizeSignature int64
		err           error
		cutePublicKey string
	)

	if sizeSignature, err = strconv.ParseInt(fmt.Sprintf("0x%s", scriptSig[:2]), 0, 64); err != nil {
		return signature, cutePublicKey, err
	}

	signature = scriptSig[2 : sizeSignature*2+2]
	cutePublicKey = scriptSig[sizeSignature*2+2:]
	return signature, cutePublicKey[4:], nil
}

func parseSignature(signature string) (*big.Int, *big.Int, error) {
	var (
		r       = &big.Int{}
		s       = &big.Int{}
		err     error
		rSize   int64
		decoder []byte
		decodes []byte
	)

	if rSize, err = strconv.ParseInt(fmt.Sprintf("0x%s", signature[6:8]), 0, 64); err != nil {
		return r, s, err
	}

	if decoder, err = hex.DecodeString(signature[8 : rSize*2+8]); err != nil {
		return r, s, err
	}

	if decodes, err = hex.DecodeString(signature[rSize*2+8+4 : len(signature)-2]); err != nil {
		return r, s, err
	}

	r.SetBytes(decoder)
	s.SetBytes(decodes)

	return r, s, nil
}

func NewCoinbase(to []byte, data string, height int64) *Transaction {
	var (
		hash        utils.Hash
		heightBytes string
	)

	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	heightBytes = utils.ConvertUnsigned4Bytes(uint32(height))

	var txout = TXOutput{Index: 0, Value: REWARD, Address: to}
	txout.Lock(to)

	var txin = TXInput{Coinbase: fmt.Sprintf("03%s", heightBytes[:len(heightBytes)-2]), Vout: -1, Sequence: "00000000"}
	var tx = Transaction{1, hash, time.Now(), 1, []TXInput{txin}, 1, []TXOutput{txout}}
	tx.SetID()

	return &tx
}

func AppendUnprocessedTransactions(tx *Transaction) {
	mutex.Lock()
	if _, ok := UnprocessedTransactions[tx.ID]; !ok {
		UnprocessedTransactions[tx.ID] = tx
	}
	mutex.Unlock()
}

func RemoveUnprocessedTransactions(tx *Transaction) {
	mutex.Lock()
	if _, ok := UnprocessedTransactions[tx.ID]; ok {
		delete(UnprocessedTransactions, tx.ID)
	}
	mutex.Unlock()
}
