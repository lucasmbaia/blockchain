package utils

import (
  "crypto/sha256"
  "bytes"
  "strings"
  "encoding/hex"
  "encoding/binary"
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

func ReverseHash(data []byte) string {
    var b = bytes.NewBuffer(make([]byte, 0, 32))

    for i := len(data); i > 0; i = i-1 {
	b.Write(data[i-1:i])
    }

    return hex.EncodeToString(b.Bytes())
}

func AddressHashSPK(data []byte) string {
  var decoded = Base58Decode(data)

  return strings.Join([]string{"76a914", hex.EncodeToString(decoded[1:len(decoded)-4]), "88ac"}, "")
}

func EncodeAmountToString(amount int) string {
  var b = make([]byte, 8)
  binary.LittleEndian.PutUint64(b, uint64(amount))
  return hex.EncodeToString(b)
}

func ConvertUnsigned4Bytes(n uint32) string {
  var b = make([]byte, 4)
  binary.LittleEndian.PutUint32(b, n)
  return hex.EncodeToString(b)
}

func ConvertUnsigned8Bytes(n uint64) string {
  var b = make([]byte, 8)
  binary.LittleEndian.PutUint64(b, n)
  return hex.EncodeToString(b)
}
