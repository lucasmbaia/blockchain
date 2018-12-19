package utils

import (
  "math/big"
  "bytes"
  "crypto/sha256"
)

var (
  BIG_RADIX = big.NewInt(58)
  BIG_ZERO  = big.NewInt(0)
)

func Base58Encode(input []byte) []byte {
  var response []byte
  var x = big.NewInt(0).SetBytes(input)
  var mod = &big.Int{}

  for x.Cmp(BIG_ZERO) != 0 {
    x.DivMod(x, BIG_RADIX, mod)
    response = append(response, alphabet[mod.Int64()])
  }

  ReverseBytes(response)
  for str := range input {
    if str == 0x00 {
      response = append([]byte{alphabet[0]}, response...)
    } else {
      break
    }
  }

  return response
}

func Base58Decode(input []byte) []byte {
  var response = big.NewInt(0)
  var count = 0

  for str := range input {
    if str == 0x00 {
      count++
    }
  }

  var pub = input[count:]
  for _, str := range pub {
    response.Mul(response, BIG_RADIX)
    response.Add(response, big.NewInt(int64(bytes.IndexByte([]byte(alphabet), str))))
  }

  var decoded = response.Bytes()
  decoded = append(bytes.Repeat([]byte{byte(0x00)}, count), decoded...)

  return decoded
}

func CheckDecode(input string) (result []byte, version byte, err error) {
    decoded := Base58Decode([]byte(input))
    if len(decoded) < 5 {
	return nil, 0, nil
    }
    version = decoded[0]
    var cksum [4]byte
    copy(cksum[:], decoded[len(decoded)-4:])
    if checksum(decoded[:len(decoded)-4]) != cksum {
	return nil, 0, nil
    }
    payload := decoded[1 : len(decoded)-4]
    result = append(result, payload...)
    return
}

func checksum(input []byte) (cksum [4]byte) {
    h := sha256.Sum256(input)
    h2 := sha256.Sum256(h[:])
    copy(cksum[:], h2[:4])
    return
}
