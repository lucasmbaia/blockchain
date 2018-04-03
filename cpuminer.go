package blockchain

import (
	"github.com/lucasmbaia/blockchain/utils"
	"log"
	"math/big"
)

const (
	MAX_NONCE = ^uint32(0)
	//MAX_NONCE = uint32(1)
)

func CpuMiner(bh *BlockHeader) (bool, utils.Hash) {
	var (
		difficult *big.Int
		hash      utils.Hash
	)

	//difficult = CalcDifficult(bh.Bits)
	difficult = CalcDifficultEasy(int(bh.Bits))

	for i := uint32(0); i < MAX_NONCE; i++ {
		bh.Nonce = i
		hash = bh.BlockHash()

		if HashToBig(&hash).Cmp(difficult) <= 0 {
			log.Println("#################### Block Mined ####################")
			return true, hash
		}
	}

	return false, hash
}

func HashToBig(hash *utils.Hash) *big.Int {
	buf := *hash
	blen := len(buf)
	for i := 0; i < blen/2; i++ {
		buf[i], buf[blen-1-i] = buf[blen-1-i], buf[i]
	}

	return new(big.Int).SetBytes(buf[:])
}

func CalcDifficultEasy(bits int) *big.Int {
	dif := big.NewInt(1)
	dif.Lsh(dif, uint(256-bits))

	return dif
}

func CalcDifficult(n uint32) *big.Int {
	mant := n & 0x007fffff
	exp := uint(n >> 24)
	negative := n&0x00800000 != 0

	var bn *big.Int

	if exp <= 3 {
		mant >>= 8 * (3 - exp)
		bn = big.NewInt(int64(mant))
	} else {
		bn = big.NewInt(int64(mant))
		bn.Lsh(bn, 8*(exp-3))
	}

	if negative {
		bn = bn.Neg(bn)
	}

	return bn
}
