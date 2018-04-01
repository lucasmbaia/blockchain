package blockchain

import (
	"encoding/hex"
	"fmt"
	"github.com/lucasmbaia/blockchain/utils"
	"testing"
	"time"
)

const (
	PREV_HASH_LT_BLOCK_125552   = "81cd02ab7e569e8bcd9317e2fe99f2de44d49ab2b8851ba4a308000000000000"
	MERKLE_HASH_LT_BLOCK_125552 = "e320b6c2fffc8d750423db8b1eb942ae710e951ed797f7affc8892b0f1fc122b"
	VERSION_BLOCK_125552        = 1
	BITS_BLOCK_125552           = 440711666
	NONCE_BLOCK_125552          = 2504433986
)

func Test_Block_Hash(t *testing.T) {
	var (
		prevH   utils.Hash
		merkleH utils.Hash
		bh      BlockHeader
		td      time.Time
		prev    []byte
		merkle  []byte
		err     error
	)

	//td = time.Unix(time.Now().Unix(), 0)
	td = time.Date(2011, time.May, 21, 17, 26, 31, 0, time.UTC)
	if prev, err = hex.DecodeString(PREV_HASH_LT_BLOCK_125552); err != nil {
		t.Fatal(err)
	}

	if merkle, err = hex.DecodeString(MERKLE_HASH_LT_BLOCK_125552); err != nil {
		t.Fatal(err)
	}

	copy(prevH[:], prev)
	copy(merkleH[:], merkle)

	bh = BlockHeader{
		Version:    int32(VERSION_BLOCK_125552),
		PrevBlock:  prevH,
		MerkleRoot: merkleH,
		Bits:       uint32(BITS_BLOCK_125552),
		Timestamp:  td,
		Nonce:      uint32(NONCE_BLOCK_125552),
	}

	hash := bh.BlockHash()

	fmt.Println(HashToBig(&hash).Cmp(CalcDifficult(440711666)))
}
