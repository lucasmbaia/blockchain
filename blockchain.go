package blockchain

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/lucasmbaia/blockchain/utils"
	"log"
	"strconv"
	"encoding/hex"
	"math/big"
	"errors"
	//"bytes"
)

const (
	DBFILE        = "/root/workspace/go/src/github.com/lucasmbaia/blockchain/blockchain.db"
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

	ctbx = NewCoinbase(bc.wa, "Coinbase Transaction")

	index++
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

func (bc *Blockchain) UnspentTransaction(pubHash []byte) ([]Transaction, error) {
  var (
    unspentTX []Transaction
    spentTX   = make(map[string][]int)
    bci	      = bc.Iterator()
    txID      string
    block     *Block
    err	      error
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
	//if bytes.Equal(out.ScriptPubKey, address) {
	  unspentTX = append(unspentTX, *tx)
	}
      }

      if !tx.IsCoinbase() {
	  fmt.Println("NAO ERA PRA ENTRAR AQUI")
	for _, in := range tx.TXInput {
	  if in.ScriptSig == string(pubHash) {
	    var inTxID = hex.EncodeToString(in.TXid)
	    spentTX[inTxID] = append(spentTX[inTxID], in.Vout)
	  }
	}
      }
    }

    if HashToBig(&block.Header.PrevBlock).Cmp(stop) == 0 {
      break
    }
  }

  return unspentTX, nil
}

func (bc *Blockchain) FindUTXO(address []byte) ([]TXOutput, error) {
  var (
    txOUT   []TXOutput
    tr	    []Transaction
    err	    error
    pubHash []byte
    decoded []byte
  )

  decoded = utils.Base58Decode(address)
  pubHash = decoded[1:len(decoded)-4]

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

func (bc *Blockchain) NewTransaction(from, to string, amount uint64) (*Transaction, error) {
    var (
	transaction *Transaction
	unspentOut  = make(map[utils.Hash][]int)
	err	    error
	//change	    = make(map[utils.Hash]int)
	decoded	    []byte
	pubHash	    []byte
	total	    uint64
	inputs	    []TXInput
	outputs	    []TXOutput
	output	    TXOutput
	hash	    utils.Hash
    )

    decoded = utils.Base58Decode([]byte(from))
    pubHash = decoded[1:len(decoded)-4]

    if unspentOut, _, total, err = bc.FindSpendable(pubHash, amount); err != nil {
	return transaction, err
    }

    if total < amount {
	return transaction, errors.New("Not enough funds")
    }

    for hash, txOut := range unspentOut {
	for _, out := range txOut {
	    inputs = append(inputs, TXInput{TXid: hash[:], Vout: out, ScriptSig: from})
	}
    }

    output = TXOutput{Value: amount}
    output.Lock([]byte(to))
    outputs = append(outputs, output)

    if total > amount {
      output = TXOutput{Value: total - (amount + MINING_RATE)}
      output.Lock([]byte(from))
      outputs = append(outputs, output)
    }

    transaction = &Transaction{1, hash, int32(len(inputs)), inputs, int32(len(outputs)), outputs}
    transaction.SetID()

    return transaction, nil
}

func (bc *Blockchain) FindSpendable(pubHash []byte, amount uint64) (map[utils.Hash][]int, map[utils.Hash]int, uint64, error) {
    var (
      unspentTX	  []Transaction
      acumulated  uint64 = 0
      err	  error
      unspentOut  = make(map[utils.Hash][]int)
      change	  = make(map[utils.Hash]int)
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
			var ctbx = NewCoinbase(wa, "Coinbase Transaction")
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
