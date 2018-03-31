package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/lucasmbaia/blockchain/utils"
	//"time"
)

type BlockHeader struct {
	Version   int32
	PrevBlock utils.Hash
	//PrevBlock  [32]byte
	MerkleRoot utils.Hash
	Bits       uint32
	//Timestamp  time.Time
	Timestamp uint32
	Nonce     uint32
}

func (b *BlockHeader) BlockHash() utils.Hash {
	var (
		buf      *bytes.Buffer
		elements []interface{}
	)

	elements = []interface{}{b.Version, b.PrevBlock, b.MerkleRoot, b.Timestamp, b.Bits, b.Nonce}
	buf = bytes.NewBuffer(make([]byte, 0, 80))

	for _, value := range elements {
		b := make([]byte, 4)

		switch value.(type) {
		case int32:
			binary.LittleEndian.PutUint32(b, uint32(value.(int32)))
			buf.Write(b)
		case uint32:
			binary.LittleEndian.PutUint32(b, value.(uint32))
			buf.Write(b)
		case utils.Hash:
			hash := value.(utils.Hash)
			buf.Write(hash[:])
		}
	}

	fmt.Printf("%x\n", buf)
	return utils.CalcDoubleHash(buf.Bytes())
}

func main() {
	var prevH utils.Hash
	prev, _ := hex.DecodeString("81cd02ab7e569e8bcd9317e2fe99f2de44d49ab2b8851ba4a308000000000000")
	copy(prevH[:], prev)

	var merkleH utils.Hash
	merkle, _ := hex.DecodeString("e320b6c2fffc8d750423db8b1eb942ae710e951ed797f7affc8892b0f1fc122b")
	copy(merkleH[:], merkle)

	b := BlockHeader{
		Version:    1,
		PrevBlock:  prevH,
		MerkleRoot: merkleH,
		//PrevBlock:  []byte{"00000000000008a3a41b85b8b29ad444def299fee21793cd8b9e567eab02cd81"},
		//MerkleRoot: []byte{"2b12fcf1b09288fcaff797d71e950e71ae42b91e8bdb2304758dfcffc2b620e3"},
		Bits:      440711666,
		Timestamp: 1305998791,
		Nonce:     2504433986,
	}

	fmt.Printf("%x\n", b.BlockHash())
	//fmt.Println()
}
