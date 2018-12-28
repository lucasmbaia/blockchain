package blockchain

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/lucasmbaia/blockchain/utils"
	"log"
	"math/big"
	"strconv"
	"time"
)

const (
	DBFILE        = "/root/go/src/github.com/lucasmbaia/blockchain/blockchain.db"
	BLOCKS_BOCKET = "blocks"
	MINING_RATE   = 10000
)

type Blockchain struct {
	tip []byte
	db  *bolt.DB
	wa  []byte
}

func (bc *Blockchain) AddBlock(data []byte) {
	var (
		hash  utils.Hash
		index int64
		err   error
		ctbx  *Transaction
	)

	if err = bc.db.View(func(tx *bolt.Tx) error {
		var bucket *bolt.Bucket = tx.Bucket([]byte(BLOCKS_BOCKET))
		var h = bucket.Get([]byte("l"))
		var i = bucket.Get([]byte("i"))

		copy(hash[:], h)

		if index, err = strconv.ParseInt(string(i), 10, 64); err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Fatalf("Error to add new block: %s\n", err)
	}

	index++
	ctbx = NewCoinbase(bc.wa, "Coinbase Transaction", index)

	var newBlock = NewBlock(int32(index), []*Transaction{ctbx}, data, hash)

	if err = bc.db.Update(func(tx *bolt.Tx) error {
		var bucket *bolt.Bucket = tx.Bucket([]byte(BLOCKS_BOCKET))

		if err = bucket.Put(newBlock.Hash[:], newBlock.Serialize()); err != nil {
			return err
		}

		if err = bucket.Put([]byte("l"), newBlock.Hash[:]); err != nil {
			return err
		}

		if err = bucket.Put([]byte("i"), []byte(fmt.Sprintf("%d", newBlock.Index))); err != nil {
			return err
		}

		bc.tip = newBlock.Hash[:]

		return nil
	}); err != nil {
		log.Fatalf("Error to add new block: %s\n", err)
	}
}

func (bc *Blockchain) UnspentTransaction(pubHash string) ([]Transaction, error) {
	var (
		unspentTX []Transaction
		spentTX   = make(map[string][]int)
		bci       = bc.Iterator()
		txID      string
		block     *Block
		err       error
		stop      = big.NewInt(0)
	)

	for {
		if block, err = bci.Next(); err != nil {
			return unspentTX, err
		}

		for _, tx := range block.Transactions {
			txID = hex.EncodeToString(tx.ID[:])
		Outputs:
			for idx, out := range tx.TXOutput {
				if spentTX[txID] != nil {
					for _, sout := range spentTX[txID] {
						if sout == idx {
							continue Outputs
						}
					}
				}

				if out.Unlock(pubHash) {
					unspentTX = append(unspentTX, *tx)
				}
			}

			/*if !tx.IsCoinbase() {
				fmt.Println("NAO ERA PRA ENTRAR AQUI")
				for _, in := range tx.TXInput {
					if in.ScriptSig == string(pubHash) {
						var inTxID = hex.EncodeToString(in.TXid)
						spentTX[inTxID] = append(spentTX[inTxID], in.Vout)
					}
				}
			}*/
		}

		if HashToBig(&block.Header.PrevBlock).Cmp(stop) == 0 {
			break
		}
	}

	return unspentTX, nil
}

func (bc *Blockchain) FindTransaction(id utils.Hash) (bool, Transaction, error) {
	var (
		bci   = bc.Iterator()
		err   error
		block *Block
		stop  = big.NewInt(0)
	)

	for {
		if block, err = bci.Next(); err != nil {
			return false, Transaction{}, err
		}

		for _, transaction := range block.Transactions {
			if bytes.Compare(transaction.ID[:], id[:]) == 0 {
				return true, *transaction, nil
			}
		}

		if HashToBig(&block.Header.PrevBlock).Cmp(stop) == 0 {
			break
		}
	}

	return false, Transaction{}, nil
}

func (bc *Blockchain) FindUTXO(address []byte) ([]TXOutput, error) {
	var (
		txOUT   []TXOutput
		tr      []Transaction
		err     error
		pubHash string
	)

	pubHash = utils.AddressHashSPK(address)

	if tr, err = bc.UnspentTransaction(pubHash); err != nil {
		return txOUT, err
	}

	for _, tx := range tr {
		for _, out := range tx.TXOutput {
			txOUT = append(txOUT, out)
		}
	}

	return txOUT, nil
}

func (bc *Blockchain) NewTransaction(from, to []byte, amount uint64) (*Transaction, error) {
	var (
		transaction *Transaction
		unspentOut  = make(map[utils.Hash][]int)
		err         error
		pubHash     string
		total       uint64
		inputs      []TXInput
		outputs     []TXOutput
		output      TXOutput
		hash        utils.Hash
	)

	pubHash = utils.AddressHashSPK(from)

	if unspentOut, _, total, err = bc.FindSpendable(pubHash, amount); err != nil {
		return transaction, err
	}

	if total < amount {
		return transaction, errors.New("Not enough funds")
	}

	for hash, txOut := range unspentOut {
		for _, out := range txOut {
			inputs = append(inputs, TXInput{TXid: hash, Vout: out})
		}
	}

	output = TXOutput{Index: 0, Value: amount, Address: to}
	output.Lock(to)
	outputs = append(outputs, output)

	if total > (amount + MINING_RATE) {
		output = TXOutput{Index: 1, Value: total - (amount + MINING_RATE), Address: from}
		output.Lock(from)
		outputs = append(outputs, output)
	}

	transaction = &Transaction{1, hash, time.Now(), int32(len(inputs)), inputs, int32(len(outputs)), outputs}
	transaction.SetID()

	return transaction, nil
}

func (bc *Blockchain) FindSpendable(pubHash string, amount uint64) (map[utils.Hash][]int, map[utils.Hash]int, uint64, error) {
	var (
		unspentTX  []Transaction
		acumulated uint64 = 0
		err        error
		unspentOut = make(map[utils.Hash][]int)
		change     = make(map[utils.Hash]int)
	)

	amount += MINING_RATE
	if unspentTX, err = bc.UnspentTransaction(pubHash); err != nil {
		return nil, nil, 0, err
	}

Unspent:
	for _, tx := range unspentTX {
		for idxOut, out := range tx.TXOutput {
			acumulated += out.Value
			unspentOut[tx.ID] = append(unspentOut[tx.ID], idxOut)

			if acumulated > amount {
				change[tx.ID] = idxOut
				break Unspent
			}
		}
	}

	return unspentOut, change, acumulated, nil
}

type BlockchainIterator struct {
	hash []byte
	db   *bolt.DB
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{bc.tip, bc.db}
}

func (bci *BlockchainIterator) Next() (*Block, error) {
	var (
		block *Block
		err   error
	)

	if err = bci.db.View(func(tx *bolt.Tx) error {
		var bucket *bolt.Bucket = tx.Bucket([]byte(BLOCKS_BOCKET))
		var encode = bucket.Get(bci.hash)
		block = Deserialize(encode)

		return nil
	}); err != nil {
		return block, err
	}

	bci.hash = block.Header.PrevBlock[:]

	return block, nil
}

func NewBlockchain(wa []byte) *Blockchain {
	var (
		tip   []byte
		db    *bolt.DB
		err   error
		valid bool
	)

	if valid = CheckValidAddress(wa); !valid {
		log.Fatalf("Invalid Wallet Address")
	}

	if db, err = bolt.Open(DBFILE, 0600, nil); err != nil {
		log.Fatalf("Error to create blockchain: %s\n", err)
	}

	if err = db.Update(func(tx *bolt.Tx) error {
		var bucket *bolt.Bucket = tx.Bucket([]byte(BLOCKS_BOCKET))

		if bucket == nil {
			var ctbx = NewCoinbase(wa, "Coinbase Transaction", 0)
			var genesis = NewGenesisBlock(ctbx)

			if bucket, err = tx.CreateBucket([]byte(BLOCKS_BOCKET)); err != nil {
				return err
			}

			if err = bucket.Put(genesis.Hash[:], genesis.Serialize()); err != nil {
				return err
			}

			if err = bucket.Put([]byte("l"), genesis.Hash[:]); err != nil {
				return err
			}

			if err = bucket.Put([]byte("i"), []byte(fmt.Sprintf("%d", genesis.Index))); err != nil {
				return err
			}

			tip = genesis.Hash[:]
		} else {
			tip = bucket.Get([]byte("l"))
		}

		return nil
	}); err != nil {
		log.Fatalf("Error to create blockchain: %s\n", err)
	}

	return &Blockchain{tip, db, wa}
}
