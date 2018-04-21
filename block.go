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

func Deserialize(b []byte) *Block {
	var block Block
	var decoder *gob.Decoder = gob.NewDecoder(bytes.NewReader(b))

	if err := decoder.Decode(&block); err != nil {
		log.Printf("Error to deserialize: %s\n", err)
	}

	return &block
}

func NewBlock(index int32, transactions []*Transaction, data []byte, prevBlockHash utils.Hash) *Block {
	bh := BlockHeader{
		Version:   1,
		PrevBlock: prevBlockHash,
		Bits:      BITS,
		Timestamp: time.Unix(time.Now().Unix(), 0),
	}

	var hash utils.Hash
	var merkle *MerkleRoot

	merkle = NewMerkleTree(transactions)
	bh.MerkleRoot = merkle.MerkleNode.Hash

	_, hash = CpuMiner(&bh)

	return &Block{
		Index:  index,
		Header: bh,
		Hash:   hash,
		Data:   data,
	}
}

func NewGenesisBlock(coinbase *Transaction) *Block {
	var prevBlockHash utils.Hash
	return NewBlock(0, []*Transaction{coinbase}, []byte("Genesis Block"), prevBlockHash)
}
