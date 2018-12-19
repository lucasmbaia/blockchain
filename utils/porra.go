package utils

const (
  VERSION = "01000000"
)

/*func lelek(transactio string, private string, address, pub []byte) {
    pk, err := crypto.ToECDS(private)
    if err != nil {
	panic(err)
    }

    pk2, err := ecdsa.GenerateKey(pk.PublicKey.Curve, rand.Reader)
    if err != nil {
	panic(err)
    }
}

func UnsignedTransaction(transaction string) {

}*/

/*func RawTransaction(txID Hash, txInputCount int, lastSPK []byte) {
  var (
      raw []byte
  )

  raw = append(raw, VERSION)
  raw = append(raw, []byte("01"))
  //raw = append(raw, ReverseHash(txID[:]))
  raw = append(raw, "eccf7e3034189b851985d871f91384b8ee357cd47c3024736e5676eb2debb3f2")
  raw = append(raw, "01000000")
  rae = append(raw, "19")
  raw = append(raw, "76a914010966776006953d5567439e5e39f86a0d273bee88ac")
  raw = append(raw, "ffffffff")
  raw = append(raw, "01")
  raw = append(raw, "605af40500000000")
  raw = append(raw, "19")
}*/
