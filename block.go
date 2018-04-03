package blockchain

import (
	"github.com/lucasmbaia/blockchain/utils"
	"time"
)

const (
	BITS = 24
)

type Block struct {
	Index  int32
	Header BlockHeader
	Hash   utils.Hash
	Data   []byte
}

type Blockchain struct {
	blocks []*Block
}

func NewBlock(index int32, data []byte, prevBlockHash utils.Hash) *Block {
	bh := BlockHeader{
		Version:   1,
		PrevBlock: prevBlockHash,
		Bits:      BITS,
		Timestamp: time.Unix(time.Now().Unix(), 0),
	}

	var hash utils.Hash

	_, hash = CpuMiner(&bh)

	return &Block{
		Index:  index,
		Header: bh,
		Hash:   hash,
		Data:   data,
	}
}

func (bc *Blockchain) AddBlock(data []byte) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	index := prevBlock.Index + 1
	newBlock := NewBlock(index, data, prevBlock.Hash)
	bc.blocks = append(bc.blocks, newBlock)
}

func NewGenesisBlock() *Block {
	var prevBlockHash utils.Hash
	return NewBlock(0, []byte("Genesis Block"), prevBlockHash)
}

func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}
