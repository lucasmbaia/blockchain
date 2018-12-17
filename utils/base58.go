package utils

import (
  "math/big"
  "bytes"
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
