package utils

import (
	"crypto/sha256"
)

type Hash [32]byte

func CalcHash(b []byte) Hash {
	return Hash(sha256.Sum256(b))
}

func CalcDoubleHash(b []byte) Hash {
	first := sha256.Sum256(b)
	return sha256.Sum256(first[:])
}
