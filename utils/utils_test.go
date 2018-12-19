package utils

import (
    "testing"
    "fmt"
)

func Test_ReverseBytes(t *testing.T) {
    /*var hash Hash
    var buf = bytes.NewBuffer(make([]byte, 0, 80))

    buf.Write([]byte("01ecaecc96b148589be10ef3f8fffc70dcc14970ed77144a787c087dfcd0b5e2"))
    copy(hash[:], buf.Bytes())*/
    hash := []byte("01ecaecc96b148589be10ef3f8fffc70dcc14970ed77144a787c087dfcd0b5e2")
    fmt.Println(ReverseHash(hash))
}

func Test_AddressHashSPK(t *testing.T) {
    address := "1JKbpvHpTzs854ztJfJv2XVyV2tQ5LZdjV"
    fmt.Println(AddressHashSPK([]byte(address)))
}
