package utils

import (
    "testing"
    "fmt"
    "encoding/hex"
)

func Test_Base58Encode(t *testing.T) {
    address := "1EoPUX89wRFGHRNGJFTzqfNcRgdzgs5JhK"
    p, _, _ := CheckDecode(address)
    fmt.Println(hex.EncodeToString(p))
}
