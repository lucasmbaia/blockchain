package server

import (
	"testing"
	"fmt"
)

func Test_SerializeAndDeserialize(t *testing.T) {
	var (
		body	[]byte
		err	error
		gossip	Gossip
	)

	if body, err = Serialize(&Gossip{
		Option:	"gossip",
	}); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(body)
	}

	if err = Deserialize(&gossip, body); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(gossip)
	}
}
