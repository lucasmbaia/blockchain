package utils

import (
	"crypto/sha256"
)

const (
	HASH_SIZE = 32
)

type Hash [HASH_SIZE]byte

func CalcHash(b []byte) Hash {
	return Hash(sha256.Sum256(b))
}

func CalcDoubleHash(b []byte) Hash {
	first := sha256.Sum256(b)
	return sha256.Sum256(first[:])
}
