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

func ReverseBytes(data []byte) {
  for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
    data[i], data[j] = data[j], data[i]
  }
}
