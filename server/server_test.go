package server

import (
	"testing"
	"context"
)

func Test_StartNode(t *testing.T) {
	StartNode(context.Background(), []byte("1CyssrDhEvZv2jXci6F5oueZwMXszm6kLs"), []string{"192.168.75.133:5688"})
}
