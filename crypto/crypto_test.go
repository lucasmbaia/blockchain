package crypto

import (
  "testing"
  "fmt"
)

func Test_ToECDSA(t *testing.T) {
  fmt.Println(ToECDS("9b178dcbb5a7dad2b3dc189baf9337138227145ef5972d0736d79f453a92a3d2"))
}
