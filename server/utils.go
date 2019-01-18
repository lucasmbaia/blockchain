package server

import (
	"encoding/json"
	"encoding/gob"
	"bytes"
)

type Gossip struct {
	Option	string
	Error	error
	Body	[]byte
}

func Serialize(i interface{}) ([]byte, error) {
	var (
		result	bytes.Buffer
		err	error
	)

	var encoder *gob.Encoder = gob.NewEncoder(&result)

	err = encoder.Encode(i)
	return result.Bytes(), err
}

func Deserialize(i interface{}, input []byte) error {
	var decoder *gob.Decoder = gob.NewDecoder(bytes.NewReader(input))
	return decoder.Decode(i)
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
