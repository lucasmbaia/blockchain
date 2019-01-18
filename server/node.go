package server

import (
	"bufio"
	"net"
	"github.com/lucasmbaia/blockchain"
	"encoding/json"
	"log"
	"io"
)

func (c *Client) listen() {
	var (
		l   net.Listener
		err error
	)

	if l, err = net.Listen("tcp", DEFAULT_PORT_SERVER); err != nil {
		panic(err)
	}

	for {
		var (
			connection  *connection
		)

		if connection.conn, err = l.Accept(); err != nil {
			panic(err)
		}

		connection.write = bufio.NewWriter(connection.conn)
		connection.read = bufio.NewReader(connection.conn)

		go c.handleNode(connection)
	}
}

func (c *Client) handleNode(conn *connection) {
	defer func() {
		if addr, ok := conn.conn.RemoteAddr().(*net.TCPAddr); ok {
			log.Printf("Connection Close if node IP %s\n", addr.IP.String())
		}
		conn.conn.Close()
	}()

	var (
		err	error
		option	[]byte
	)

	for {
		if option, err = conn.read.ReadBytes('\n'); err != nil {
			if err != io.EOF {
				log.Printf("Error to read bytes: %s\n", err.Error())
			}
			return
		}

		var gossip *gossip

		if err = json.Unmarshal(option[:len(option)-1], &gossip); err != nil {
			log.Printf("Error to deserialize gossip: %s\n", err.Error())
			continue
		}

		switch gossip.Option {
		case "transaction":
			var tx *blockchain.Transaction

			if err = json.Unmarshal(gossip.Body, &tx); err != nil {
				log.Printf("Error to deserialize transaction: %s\n", err.Error())
			}

			if err = tx.ValidTransaction(getTransactionsInBlocks(tx)); err == nil {
				blockchain.AppendUnprocessedTransactions(tx)
				c.transaction <- tx
			}
		}
	}
}
