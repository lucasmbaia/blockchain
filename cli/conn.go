package cli

import (
	"net"
	"encoding/json"
	"fmt"
	"bufio"
	"io"
)

const (
	NODE  = "192.168.75.133:5688"
)

type gossip struct {
	Option  string
	Body    []byte
	Error	error
}

func transmit(g gossip) error {
	var (
		err   error
		conn  net.Conn
		body  []byte
		done  = make(chan struct{})
		errc  = make(chan error, 1)
	)

	go func() {
		if conn, err = net.Dial("tcp", NODE); err != nil {
			errc <- err
			return
		}
		defer conn.Close()

		if body, err = encodeGossip(g); err != nil {
			errc <- err
			return
		}

		if _, err = conn.Write(body); err != nil {
			errc <- err
			return
		}

		response(conn, done, errc)
	}()

	select {
	case err = <-errc:
		return err
	case <-done:
		return nil
	}
}

func response(conn net.Conn, done chan struct{}, errc chan error) {
	var (
		option	[]byte
		err	error
		b	*bufio.Reader
	)

	b = bufio.NewReader(conn)

	for {
		if option, err = b.ReadBytes('\n'); err != nil {
			if err != io.EOF {
				errc <- err
			}
			return
		}

		var g gossip

		if err = json.Unmarshal(option[:len(option)-1], &g); err != nil {
			errc <- err
			return
		}

		switch g.Option {
		case "local_transaction":
			if g.Error != nil {
				errc <- err
			} else {
				done <- struct{}{}
			}

			return
		}
	}
}

func encodeGossip(g gossip) ([]byte, error) {
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
