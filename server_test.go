package blockchain

import (
    "testing"
    "context"
)

func Test_StartFullNode(t *testing.T) {
    StartFullNode(context.Background(), []byte("1CyssrDhEvZv2jXci6F5oueZwMXszm6kLs"), []string{})
}
