package blockchain

import (
	"bytes"
	"encoding/gob"
	"github.com/lucasmbaia/blockchain/utils"
	"log"
	"time"
)

const (
	BITS = 24
)

type Block struct {
	Index        int32
	Transactions []*Transaction
	Header       BlockHeader
	Hash         utils.Hash
	Data         []byte
}

func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	var encoder *gob.Encoder = gob.NewEncoder(&result)

	if err := encoder.Encode(b); err != nil {
		log.Printf("Error to serialize block: %s\n", err)
	}

	return result.Bytes()
}

func HashTransactions(transactions []*Transaction) utils.Hash {
	var (
		txHash  [][]byte
	)

	for _, tx := range transactions {
		txHash = append(txHash, tx.ID[:])
	}

	return utils.CalcHash(bytes.Join(txHash, []byte{}))
}

func (b *Block) CheckProcessedTransactions(tx []*Transaction) {
	for idx, transaction := range tx {
		for _, btx := range b.Transactions {
			if bytes.Compare(transaction.ID[:], btx.ID[:]) == 0 {
				if len(tx) -1 == idx {
					tx = tx[:idx]
				} else {
					tx = append(tx[:idx], tx[idx+1:]...)
				}
			}
		}
	}
}

func Deserialize(b []byte) *Block {
	var block Block
	var decoder *gob.Decoder = gob.NewDecoder(bytes.NewReader(b))

	if err := decoder.Decode(&block); err != nil {
		log.Printf("Error to deserialize: %s\n", err)
	}

	return &block
}

func NewBlock(operations Operations, index int32, transactions []*Transaction, data []byte, prevBlockHash utils.Hash) *Block {
	bh := BlockHeader{
		Version:   1,
		PrevBlock: prevBlockHash,
		Bits:      BITS,
		Timestamp: time.Unix(time.Now().Unix(), 0),
	}

	var hash utils.Hash
	var merkle *MerkleRoot
	var valid bool

	merkle = NewMerkleTree(transactions)
	bh.MerkleRoot = merkle.MerkleNode.Hash
	bh.HashTransactions = HashTransactions(transactions)

	//valid, hash = CpuMiner(&bh)
	valid, hash = CpuMinerControl(operations, &bh)

	if !valid {
		return NewBlock(operations, index, transactions, data, prevBlockHash)
	}

	return &Block{
		Index:	      index,
		Transactions: transactions,
		Header:	      bh,
		Hash:	      hash,
		Data:	      data,
	}
}

func NewGenesisBlock(coinbase *Transaction) *Block {
	var prevBlockHash utils.Hash
	var operations Operations

	return NewBlock(operations, 0, []*Transaction{coinbase}, []byte("Genesis Block"), prevBlockHash)
}
