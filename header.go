package blockchain

import (
	"bytes"
	"encoding/binary"
	"github.com/lucasmbaia/blockchain/utils"
	"time"
)

type BlockHeader struct {
	Version    int32
	PrevBlock  utils.Hash
	MerkleRoot utils.Hash
	Bits       uint32
	Timestamp  time.Time
	Nonce      uint32
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
		case time.Time:
			binary.LittleEndian.PutUint32(b, uint32(value.(time.Time).Unix()))
			buf.Write(b)
		}
	}

	return utils.CalcDoubleHash(buf.Bytes())
}
