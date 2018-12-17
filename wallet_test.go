package blockchain

import (
  "testing"
  "fmt"
)

func Test_NewWallet(t *testing.T) {
  if w, err := NewWallet(); err != nil {
    t.Fatal(err)
  } else {
    fmt.Println(w.PrivateToHex())
    fmt.Println(string(w.Address))
  }
}

func Test_CheckValidAddress(t *testing.T) {
  valid := CheckValidAddress([]byte("18xB6w3WrNDriHruodzxsRFsiqYx6VufYY"))
  fmt.Println(valid)
}

func Test_UnlockWallet(t *testing.T) {
  if valid, _, err := UnlockWallet("668211e92d7030820d4c529a7cbf6da2b6a504cdd7527a9c892db5a08df3813b", "1N2JixT4qq6X48ww3qcpnDPHhUCDesDcH2"); err != nil {
    t.Fatal(err)
  } else {
    fmt.Println(valid)
  }
}

func Test_AddressToPubHash(t *testing.T) {
  AddressToPubHash([]byte("17v71YtfJoUCtKAfyYg8A4ppTQryqGDsNz"))
}
