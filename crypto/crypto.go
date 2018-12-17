package crypto

import (
  "crypto/elliptic"
  "crypto/ecdsa"
  "math/big"
  "encoding/hex"
)

const (
  PRIVATE_KEY_LEN = 32
)

func ToECDS(input string) (*ecdsa.PrivateKey, error) {
  var (
    dec	  []byte
    err	  error
    curve elliptic.Curve
    x	  *big.Int
    y	  *big.Int
  )

  if dec, err = hex.DecodeString(input); err != nil {
    return nil, err
  }

  curve = elliptic.P256()
  x, y = curve.ScalarBaseMult(dec)

  return &ecdsa.PrivateKey{
    PublicKey:	ecdsa.PublicKey{
      Curve:  curve,
      X:      x,
      Y:      y,
    },
    D:	new(big.Int).SetBytes(dec),
  }, nil
}

func ECDSToHEX(input []byte) string {
  var private = make([]byte, 0, PRIVATE_KEY_LEN)
  private = zeroAppend(PRIVATE_KEY_LEN, private, input)
  return hex.EncodeToString(private)
}

func zeroAppend(size int, dst, src []byte) []byte {
  for i := 0; i < (size - len(src)); i++ {
    dst = append(dst, 0)
  }

  return append(dst, src...)
}
