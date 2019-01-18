package server

import (
	"encoding/json"
	//"encoding/gob"
	//"bytes"
)

type Gossip struct {
	Option	string
	Error	error
	Body	[]byte
}

func EncodeGossip(g Gossip) ([]byte, error) {
	var (
		body  []byte
		err   error
	)

	if body, err = json.Marshal(g); err != nil {
		return body, err
	}

	body = append(body, byte('\n'))
	return body, nil
}
